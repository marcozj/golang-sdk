package restapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	logger "github.com/marcozj/golang-sdk/logging"
)

const (
	SourceHeader string = "golang-sdk"
)

type HttpClientFactory func() *http.Client

// BaseAPIResponse represents the most basic standard Centrify API response,
//	where Result itself is left as raw json
type BaseAPIResponse struct {
	Success   bool `json:"success"`
	Result    json.RawMessage
	Message   string
	Exception string
}

type StringResponse struct {
	BaseAPIResponse
	Result string
}

type BoolResponse struct {
	BaseAPIResponse
	Result bool
}

// GenericMapResponse represents Centrify API responses where results are map[string]interface{},
//	this type allows direct access to these without further decoding.
type GenericMapResponse struct {
	BaseAPIResponse
	Result map[string]interface{}
}

type SliceResponse struct {
	BaseAPIResponse
	Result []interface{}
}

type HttpError struct {
	error          // error type
	StatusCode int // HTTP status
}

type RestClientMode uint32

// RestClient represents a stateful API client (cookies maintained between calls, single service etc)
type RestClient struct {
	Service         string
	Client          *http.Client
	Headers         map[string]string
	SourceHeader    string
	ResponseHeaders http.Header
}

// GetNewRestClient creates a new RestClient for the specified endpoint.  If a factory for creating
//	http.Client's is not provided, you'll get a new: &http.Client{}.
func GetNewRestClient(service string, httpFactory HttpClientFactory) (*RestClient, error) {
	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, err
	}

	// Munge on the service a little bit, force it to have no trailing / and always start with https://
	url, err := url.Parse(service)
	if err != nil {
		return nil, err
	}
	url.Scheme = "https"
	url.Path = ""

	client := &RestClient{}
	client.Service = url.String()
	if httpFactory != nil {
		client.Client = httpFactory()
	} else {
		client.Client = &http.Client{}
	}
	client.Client.Jar = jar
	client.Headers = make(map[string]string)
	client.SourceHeader = SourceHeader
	return client, nil
}

func (r *RestClient) CallRawAPI(method string, args map[string]interface{}) ([]byte, error) {
	return r.postAndGetBody(method, args)
}

func (r *RestClient) CallBaseAPI(method string, args map[string]interface{}) (*BaseAPIResponse, error) {
	body, err := r.postAndGetBody(method, args)
	if err != nil {
		return nil, err
	}
	return bodyToBaseAPIResponse(body)
}

func (r *RestClient) CallGenericMapAPI(method string, args map[string]interface{}) (*GenericMapResponse, error) {
	body, err := r.postAndGetBody(method, args)
	if err != nil {
		return nil, err
	}
	return bodyToGenericMapResponse(body)
}

func (r *RestClient) CallStringAPI(method string, args map[string]interface{}) (*StringResponse, error) {
	body, err := r.postAndGetBody(method, args)
	if err != nil {
		return nil, err
	}
	return bodyToStringResponse(body)
}

func (r *RestClient) CallBoolAPI(method string, args map[string]interface{}) (*BoolResponse, error) {
	body, err := r.postAndGetBody(method, args)
	if err != nil {
		return nil, err
	}
	return bodyToBoolResponse(body)
}

func (r *RestClient) CallSliceAPI(method string, args map[string]interface{}) (*SliceResponse, error) {
	body, err := r.postAndGetBody(method, args)
	if err != nil {
		return nil, err
	}
	return bodyToSliceResponse(body)
}

func (r *RestClient) postAndGetBody(method string, args map[string]interface{}) ([]byte, error) {
	postreq, err := r.formHttpRequest(method, args)
	if err != nil {
		logger.ErrorTracef(err.Error())
		return nil, err
	}

	httpresp, err := r.Client.Do(postreq)
	if err != nil {
		r.ResponseHeaders = nil
		logger.ErrorTracef(err.Error())
		return nil, err
	}
	defer httpresp.Body.Close()

	// save response heasder
	r.ResponseHeaders = httpresp.Header

	if httpresp.StatusCode == 200 {
		body, err := ioutil.ReadAll(httpresp.Body)
		return body, err
	}

	body, _ := ioutil.ReadAll(httpresp.Body)
	return nil, &HttpError{error: fmt.Errorf("POST to %s failed with code %d, body: %s", method, httpresp.StatusCode, body), StatusCode: httpresp.StatusCode}
}

// GetLastResponseHeaders returns the response headers from last REST call
func (r *RestClient) GetLastResponseHeaders() http.Header {
	return r.ResponseHeaders
}

// This function converts a map[string]interface{} into json string
func payloadFromMap(input map[string]interface{}) string {
	if input != nil {
		p, _ := json.Marshal(input)
		return string(p)
	}

	return ""
}

func bodyToBaseAPIResponse(body []byte) (*BaseAPIResponse, error) {
	reply := &BaseAPIResponse{}
	err := json.Unmarshal(body, &reply)
	if err != nil {
		logger.ErrorTracef(err.Error())
		return nil, fmt.Errorf("Failed to unmarshal BaseApiResponse from HTTP response: %v", err)
	}
	return reply, nil
}

func bodyToGenericMapResponse(body []byte) (*GenericMapResponse, error) {
	reply := &GenericMapResponse{}
	err := json.Unmarshal(body, &reply)
	if err != nil {
		logger.ErrorTracef(err.Error())
		return nil, fmt.Errorf("Failed to unmarshal GenericMapResponse from HTTP response: %v", err)
	}
	return reply, nil
}

func bodyToStringResponse(body []byte) (*StringResponse, error) {
	reply := &StringResponse{}
	err := json.Unmarshal(body, &reply)
	if err != nil {
		logger.ErrorTracef(err.Error())
		return nil, fmt.Errorf("Failed to unmarshal StringResponse from HTTP response: %v", err)
	}
	return reply, nil
}

func bodyToBoolResponse(body []byte) (*BoolResponse, error) {
	reply := &BoolResponse{}
	err := json.Unmarshal(body, &reply)
	if err != nil {
		logger.ErrorTracef(err.Error())
		return nil, fmt.Errorf("Failed to unmarshal BoolResponse from HTTP response: %v", err)
	}
	return reply, nil
}

func bodyToSliceResponse(body []byte) (*SliceResponse, error) {
	reply := &SliceResponse{}
	err := json.Unmarshal(body, &reply)
	if err != nil {
		logger.ErrorTracef(err.Error())
		return nil, fmt.Errorf("Failed to unmarshal BoolResponse from HTTP response: %v", err)
	}
	return reply, nil
}

// Custom functions

// This function converts a list of map[string]interface{} into json string
func payloadFromList(input []map[string]interface{}) string {
	if input != nil {
		p, _ := json.Marshal(input)
		return string(p)
	}

	return ""
}

// CallGenericMapListAPI is currently used by admin right assignment and removal for Role
func (r *RestClient) CallGenericMapListAPI(method string, args []map[string]interface{}) (*GenericMapResponse, error) {
	body, err := r.postAndGetBodyList(method, args)
	if err != nil {
		return nil, err
	}
	return bodyToGenericMapResponse(body)
}

func (r *RestClient) postAndGetBodyList(method string, args []map[string]interface{}) ([]byte, error) {
	service := strings.TrimSuffix(r.Service, "/")
	method = strings.TrimPrefix(method, "/")
	postdata := strings.NewReader(payloadFromList(args))
	logger.Debugf("Post url: %s", service+"/"+method)
	logger.Debugf("Post json: %+v", postdata)
	postreq, err := http.NewRequest("POST", service+"/"+method, postdata)

	if err != nil {
		return nil, err
	}

	postreq.Header.Add("Content-Type", "application/json")
	postreq.Header.Add("X-CENTRIFY-NATIVE-CLIENT", "Yes")
	postreq.Header.Add("X-CFY-SRC", r.SourceHeader)

	for k, v := range r.Headers {
		postreq.Header.Add(k, v)
	}

	httpresp, err := r.Client.Do(postreq)
	if err != nil {
		r.ResponseHeaders = nil
		logger.ErrorTracef(err.Error())
		return nil, err
	}

	defer httpresp.Body.Close()

	// save response heasder
	r.ResponseHeaders = httpresp.Header

	if httpresp.StatusCode == 200 {
		return ioutil.ReadAll(httpresp.Body)
	}

	body, _ := ioutil.ReadAll(httpresp.Body)
	return nil, &HttpError{error: fmt.Errorf("POST to %s failed with code %d, body: %s", method, httpresp.StatusCode, body), StatusCode: httpresp.StatusCode}
}

func (r *RestClient) DownloadFile(method string, args map[string]interface{}, filepath string) error {
	postreq, err := r.formHttpRequest(method, args)
	if err != nil {
		logger.ErrorTracef(err.Error())
		return err
	}

	httpresp, err := r.Client.Do(postreq)
	if err != nil {
		r.ResponseHeaders = nil
		logger.ErrorTracef(err.Error())
		return err
	}
	defer httpresp.Body.Close()

	// save response heasder
	r.ResponseHeaders = httpresp.Header

	if httpresp.StatusCode == 200 {
		// Create the file
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, httpresp.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RestClient) formHttpRequest(method string, args map[string]interface{}) (*http.Request, error) {
	service := strings.TrimSuffix(r.Service, "/")
	method = strings.TrimPrefix(method, "/")
	postdata := strings.NewReader(payloadFromMap(args))
	logger.Debugf("Post url: %s", service+"/"+method)
	logger.Debugf("Post json: %+v", postdata)
	postreq, err := http.NewRequest("POST", service+"/"+method, postdata)

	if err != nil {
		logger.ErrorTracef(err.Error())
		return nil, err
	}

	postreq.Header.Add("Content-Type", "application/json")
	postreq.Header.Add("X-CENTRIFY-NATIVE-CLIENT", "Yes")
	postreq.Header.Add("X-CFY-SRC", r.SourceHeader)

	for k, v := range r.Headers {
		postreq.Header.Add(k, v)
	}

	return postreq, nil
}
