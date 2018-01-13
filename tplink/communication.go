package tplink

import (
	"net"
	"encoding/json"
	"time"
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
	result := make([]byte, len(request))
	for i,  c := range request {
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
	return nil, decrypt(buff[4:n], 171)
}

type Discovered struct {
	Addr	*net.UDPAddr
	Data	[]byte
}

func discovered(r chan Discovered, addr *net.UDPAddr, rlen int, buff []byte) {
	r <- Discovered{
		Addr: addr,
		Data: decrypt(buff[:rlen], 171),
	}
}

func Discover(request interface{}, timeout int) (error, []Discovered) {

	found := []Discovered{}
	
	BroadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9999")
	if err != nil {
		return err, found
	}

	FromAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8755")
	if err != nil {
		return err, found
	}

	sock, err := net.ListenUDP("udp", FromAddr)
	defer sock.Close()
	if err != nil {
		return err, found
	}
	sock.SetReadBuffer(2048)

	r := make(chan Discovered)

	go func(s *net.UDPConn, request interface{}) {
		for {
			buff := make([]byte, 2048)
			rlen, addr, err := s.ReadFromUDP(buff)
			if err != nil {
				break
			}
			go discovered(r, addr, rlen, buff)
		}
	}(sock, request)

	e_json, err := json.Marshal(&request)
	if err != nil {
		return err, found
	}

	var encrypted = encrypt(e_json, 171)[4:]
	_, err = sock.WriteToUDP(encrypted, BroadcastAddr)
	if err != nil {
		return err, found
	}
	started := time.Now()
	Q:
		for {
			select {
				case x := <-r:
					found = append(found, x)
				default:
					if now := time.Now(); now.Sub(started) >= time.Duration(timeout) * time.Second {
						break Q
					}
			}
		}
	return nil, found
}