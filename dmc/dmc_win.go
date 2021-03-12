// +build windows

package dmc

// This file is for Windows platform
import (
	"fmt"
	"reflect"

	// npipe package only works on Windows
	"gopkg.in/natefinch/npipe.v2"
)

// LRPC2 represents local RPC data structure
type LRPC2 struct {
	winNamePipe      string
	maxPayloadLength uint32
	pid              int
	client           *npipe.PipeConn
}

// NewLRPC2 initiates a new local RPC client
func NewLRPC2() *LRPC2 {
	lrpc := LRPC2{}
	lrpc.winNamePipe = `\\.\pipe\cagent_admins`

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
	c, err := npipe.Dial(lrpc.winNamePipe)
	if err != nil {
		return err
	}

	lrpc.client = c
	// Do handshake to check version number
	lrpc.client.Write(uint32ToByteArray(lrpc2Version))
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

	return "", fmt.Errorf("Invaid reply from server: %+v", "")
}
