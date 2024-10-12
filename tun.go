package main

import (
	"errors"
	"net"
	"time"

	"golang.zx2c4.com/wireguard/tun"
)

// tunConn wraps a tun.Device to provide some net.Conn-like functionality.
type tunConn struct {
	device tun.Device
}

// Read reads data from the TUN device.
func (t *tunConn) Read(b []byte) (int, error) {
	return t.device.Read(b, 0)
}

// Write writes data to the TUN device.
func (t *tunConn) Write(b []byte) (int, error) {
	return t.device.Write(b, 0)
}

// Close closes the TUN device.
func (t *tunConn) Close() error {
	return t.device.Close()
}

// LocalAddr and RemoteAddr are not applicable to a TUN device, but you must
// define them to satisfy the net.Conn interface. You could return a dummy value.
func (t *tunConn) LocalAddr() net.Addr {
	// Not applicable for TUN devices, return nil or a dummy value.
	return nil
}

func (t *tunConn) RemoteAddr() net.Addr {
	// Not applicable for TUN devices, return nil or a dummy value.
	return nil
}

// SetDeadline, SetReadDeadline, and SetWriteDeadline are not supported by TUN devices.
// You could define them but have them return an error indicating they are not supported.
func (t *tunConn) SetDeadline(_ time.Time) error {
	return errors.New("SetDeadline not supported on TUN devices")
}

func (t *tunConn) SetReadDeadline(_ time.Time) error {
	return errors.New("SetReadDeadline not supported on TUN devices")
}

func (t *tunConn) SetWriteDeadline(_ time.Time) error {
	return errors.New("SetWriteDeadline not supported on TUN devices")
}
