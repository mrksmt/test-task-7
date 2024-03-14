package helpers

import (
	"bytes"
	"io"
	"net"
	"syscall"
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

func ConnCheck(conn net.Conn) error {
	var sysErr error
	rawConn, err := conn.(syscall.Conn).SyscallConn()
	if err != nil {
		return err
	}
	err = rawConn.Read(func(fd uintptr) bool {
		buf := make([]byte, 1)
		n, _, err := syscall.Recvfrom(int(fd), buf, syscall.MSG_PEEK|syscall.MSG_DONTWAIT)
		switch {
		case n == 0 && err == nil:
			sysErr = io.EOF
		case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
			// no-op
		default:
			sysErr = err
		}
		return true
	})
	if err != nil {
		return err
	}

	return sysErr
}
