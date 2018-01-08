package tplink

import (
	"net"
	"encoding/json"
)

func encrypt(request []byte, key uint8) []byte {
	result := make([]byte, 4 + len(request))
	result[0] = 0x0
	result[1] = 0x0
	result[2] = 0x0
	result[3] = 0x0
	for i,  c := range request {
		var a = key ^ uint8(c)
		key = uint8(a)
		result[i+4] = a
	}
	return result
}

func decrypt(request []byte, key uint8) []byte {
	result := make([]byte, len(request) - 4)
	for i,  c := range request[4:] {
		var a = key ^ uint8(c)
		key = uint8(c)
		result[i] = a
	}
	return result
}

func Query(address string, request interface{}) (error, []byte) {
	conn, err := net.Dial("tcp", address)

	defer conn.Close()

	if err != nil {
		return err, make([]byte, 0)
	}

	e_json, err := json.Marshal(&request)
	if err != nil {
		return err, make([]byte, 0)
	}

	var encrypted = encrypt(e_json, 171)
	conn.Write(encrypted)

	buff := make([]byte, 2048)
	n, err := conn.Read(buff)
	if err != nil {
		return err, make([]byte, 0)
	}
	return nil, decrypt(buff[:n], 171)
}