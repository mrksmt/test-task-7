package helpers

import (
	"bytes"
	"net"
)

const chunkSize = 1024

func Receive(conn net.Conn) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	for {
		chunk := make([]byte, chunkSize)
		read, err := conn.Read(chunk)
		if err != nil {
			return buffer.Bytes(), err
		}
		buffer.Write(chunk[:read])
		if read < chunkSize {
			break
		}
	}
	return buffer.Bytes(), nil
}
