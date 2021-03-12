package dmc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

const (
	lrpc2Version             uint32 = 4          // LRPC version for handshake
	handShakeAck             uint32 = 1          // handshare acknowledged
	msgIDAdminClientGetToken uint16 = 1500       // Command
	msgDataTypeString        byte   = 4          // String data type
	msgEnd                   byte   = 0          // End of message
	magicNumber              uint32 = 0xABCD8012 // magic number
	headerLength             uint16 = 34         // Header length
	msgDataTypeInt32         byte   = 2          // Int32 data type
	msgDataTypeSet           byte   = 7          // Set data type
)

func constructPayload(scope string) []byte {
	var payload []byte
	payload = append(payload, uint16ToByteArray(msgIDAdminClientGetToken)...)
	payload = append(payload, []byte{msgDataTypeString}...)
	payload = append(payload, uint32ToByteArray(uint32(len(scope)))...)
	payload = append(payload, []byte(scope)...)
	payload = append(payload, []byte{msgEnd}...)

	return payload
}

func sendRequest(payload []byte, pid int, maxLength uint32, readFunc func(int) ([]byte, error), writeFunc func([]byte) (int, error)) ([]interface{}, error) {
	// Construct request header
	var data []byte
	data = append(data, uint32ToByteArray(magicNumber)...)  // magic number
	data = append(data, uint16ToByteArray(headerLength)...) // header length
	data = append(data, uint32ToByteArray(lrpc2Version)...) // LRPC2 version
	data = append(data, uint64ToByteArray(uint64(pid))...)  // process id
	seq := rand.Int31()

	data = append(data, uint32ToByteArray(uint32(seq))...) // sequence number
	now := time.Now()
	data = append(data, uint64ToByteArray(uint64(now.Unix()))...) // seconds passed since epoch
	length := uint32(len(payload))
	if length > maxLength {
		return nil, fmt.Errorf("LRPC payload length %d exceeds the max limit %d", length, maxLength)
	}
	data = append(data, uint32ToByteArray(length)...)
	data = append(data, payload...) // payload

	// Send request header
	_, err := writeFunc(data) // send request
	if err != nil {
		return nil, fmt.Errorf("Write error: %v", err)
	}

	// Read response
	scan := true
	var items []interface{}
	for scan {
		// Read the first 34 bytes
		header, _ := readFunc(int(headerLength))
		bmagicNumber := header[0:4]
		bheaderLength := header[4:6]
		_ = bheaderLength
		bversion := header[6:10]
		_ = bversion
		bpid := header[10:18]
		_ = bpid
		bseq := header[18:22]
		bts := header[22:30]
		_ = bts
		blength := header[30:34]
		payloadLength := byteArrayToUInt16(blength)
		payload, _ = readFunc(int(payloadLength))
		responseSeq := byteArrayToUInt32(bseq)

		// Compare magic number
		if bytes.Compare(bmagicNumber, uint32ToByteArray(magicNumber)) != 0 {
			return nil, fmt.Errorf("Unrecognized LRPC2 server")
		}
		// Compare sequence number
		if int32(responseSeq) != seq {
			// Skip unmatched seq number.  It could be the response from prevoius requests
			continue
		}
		command := payload[0:2]
		_ = command
		payload = payload[2:]

		// Process payload
		for len(payload) > 0 {
			itemType := payload[0:1]
			if bytes.Compare(itemType, []byte{msgDataTypeInt32}) == 0 {
				item := byteArrayToUInt32(payload[1:5])
				items = append(items, item)
				payload = payload[5:]
			} else if bytes.Compare(itemType, []byte{msgDataTypeString}) == 0 {
				strlen := int32(byteArrayToUInt32(payload[1:5]))
				if strlen < 0 {
					items = append(items, "")
					payload = payload[5:]
				} else {
					item := string(payload[5 : 5+strlen])
					items = append(items, item)
					payload = payload[5+strlen:]
				}
			} else if bytes.Compare(itemType, []byte{msgDataTypeSet}) == 0 {
				count := byteArrayToUInt32(payload[1:5])
				payload = payload[5:]
				var strset []string
				for i := 0; i < int(count); i++ {
					strlen := byteArrayToUInt32(payload[1:5])
					if strlen < 0 {
						strset = append(strset, "")
						payload = payload[5:]
					} else {
						item := string(payload[5 : 5+strlen])
						strset = append(strset, item)
						payload = payload[5+strlen:]
					}
				}
				items = append(items, strset)
			} else if bytes.Compare(itemType, []byte{msgEnd}) == 0 {
				scan = false
				break
			} else {
				return nil, fmt.Errorf("Unrecognized data type %v", itemType)
			}
		}

	}

	return items, nil
}

func constructRequestHeader(payload []byte, seq int32, pid int, maxLenght uint32) ([]byte, error) {
	// Construct request header
	var data []byte
	data = append(data, uint32ToByteArray(magicNumber)...)  // magic number
	data = append(data, uint16ToByteArray(headerLength)...) // header length
	data = append(data, uint32ToByteArray(lrpc2Version)...) // LRPC2 version
	data = append(data, uint64ToByteArray(uint64(pid))...)  // process id
	data = append(data, uint32ToByteArray(uint32(seq))...)  // sequence number
	now := time.Now()
	data = append(data, uint64ToByteArray(uint64(now.Unix()))...) // seconds passed since epoch
	length := uint32(len(payload))
	if length > maxLenght {
		return nil, fmt.Errorf("LRPC payload length %d exceeds the max limit %d", length, maxLenght)
	}
	data = append(data, uint32ToByteArray(length)...)
	data = append(data, payload...) // payload

	return data, nil
}

func uint16ToByteArray(num uint16) []byte {
	arr := make([]byte, 2)
	binary.LittleEndian.PutUint16(arr, num)
	return arr
}

func uint32ToByteArray(num uint32) []byte {
	arr := make([]byte, 4)
	binary.LittleEndian.PutUint32(arr, num)
	return arr
}

func uint64ToByteArray(num uint64) []byte {
	arr := make([]byte, 8)
	binary.LittleEndian.PutUint64(arr, num)
	return arr
}

func byteArrayToUInt32(arr []byte) uint32 {
	val := uint32(0)
	val = binary.LittleEndian.Uint32(arr)
	return val
}

func byteArrayToUInt16(arr []byte) uint16 {
	val := uint16(0)
	val = binary.LittleEndian.Uint16(arr)
	return val
}
