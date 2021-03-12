// +build !windows

package dmc

// This file is for non Windows platform
import (
	"fmt"
	"net"
	"os"
	"reflect"
)

// LRPC2 represents local RPC data structure
type LRPC2 struct {
	unixSocketFile   string
	maxPayloadLength uint32
	pid              int
	client           net.Conn
}

// NewLRPC2 initiates a new local RPC client
func NewLRPC2() *LRPC2 {
	lrpc := LRPC2{}
	lrpc.unixSocketFile = "/var/centrify/cloud/daemon2"

	return &lrpc
}

func (lrpc *LRPC2) reader(size int) ([]byte, error) {
	buf := make([]byte, size)
	_, err := lrpc.client.Read(buf[:])
	if err != nil {
		return nil, fmt.Errorf("Error reading from server: %v", err)
	}

	return buf, nil
}

func (lrpc *LRPC2) connect() error {
	c, err := net.Dial("unix", lrpc.unixSocketFile)

	if err != nil {
		fmt.Printf("Error connecting to local rpc server: %v", err)
		return fmt.Errorf("Error connecting to local rpc server: %v", err)
	}
	lrpc.client = c
	//defer c.Close()

	lrpc.pid = os.Getegid()
	// Do handshake to check version number
	_, err = lrpc.client.Write(uint32ToByteArray(lrpc2Version))
	if err != nil {
		return fmt.Errorf("Error in handshake: %v", err)
	}

	// Read response for version
	data, err := lrpc.reader(4)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	answer := byteArrayToUInt32(data)
	if answer != handShakeAck {
		return fmt.Errorf("Server doesn't support LRPC2 version 4")
	}
	data, err = lrpc.reader(4)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	lrpc.maxPayloadLength = byteArrayToUInt32(data)

	return nil
}

func (lrpc *LRPC2) request(payload []byte) ([]interface{}, error) {
	return sendRequest(payload, lrpc.pid, lrpc.maxPayloadLength, lrpc.reader, lrpc.client.Write)
}

// GetToken gets dmc token from Centrify Client service
func (lrpc *LRPC2) GetToken(scope string) (string, error) {
	// Connect to lrpc server
	err := lrpc.connect()
	if err != nil {
		return "", fmt.Errorf("Failed to connect to server: %v", err)
	}
	defer lrpc.client.Close()

	// Send payload to lrpc server
	reply, err := lrpc.request(constructPayload(scope))
	if err != nil {
		return "", fmt.Errorf("Request error: %v", err)
	}

	if len(reply) > 2 {
		statusCode := fmt.Sprintf("%v", reply[0])
		if statusCode != "0" {
			return "", fmt.Errorf("%+v", reply[1])
		}
		if reflect.TypeOf(reply[2]).Kind() != reflect.String {
			return "", fmt.Errorf("Invaid reply from server: %+v", reply)
		}
		return reply[2].(string), nil
	}

	return "", fmt.Errorf("Invaid reply from server: %+v", reply)
}
