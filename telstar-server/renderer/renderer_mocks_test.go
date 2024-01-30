package renderer

import (
	"net"
	"time"
)

//type Conn interface{
//	Close()
//	Write()
//}

type MockConn struct {
	buffer []byte
}

func (c *MockConn) Read(b []byte) (n int, err error) {
	b = append(b, 0xff)
	return len(b), nil
}

func (c *MockConn) Write(b []byte) (n int, err error) {
	for i := 0; i < len(b); i++ {
		c.buffer = append(c.buffer, b[i])
	}
	return len(b), nil
}

func (c *MockConn) Close() error {
	panic("implement me")
}

func (c *MockConn) LocalAddr() net.Addr {
	panic("implement me")
}

func (c *MockConn) RemoteAddr() net.Addr {
	panic("implement me")
}

func (c *MockConn) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (c *MockConn) SetReadDeadline(t time.Time) error {
	panic("implement me")
}

func (c *MockConn) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}
