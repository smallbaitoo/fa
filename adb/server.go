package adb

import (
	"net"
	"sync"

	"github.com/qiniu/log"
)

// ADBDaemon implement service for command: adb connect
type ADBDaemon struct {
	device  *Device
	remotes map[string]net.Conn
	mu      sync.Mutex
}

func NewADBDaemon(device *Device) *ADBDaemon {
	return &ADBDaemon{
		device:  device,
		remotes: make(map[string]net.Conn),
	}
}

func (s *ADBDaemon) ListenAndServe(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(ln)
}

func (s *ADBDaemon) Serve(ln net.Listener) error {
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		remoteAddress := conn.RemoteAddr().String()
		log.Infof("Incomming request from: %v", remoteAddress)

		s.mu.Lock()
		s.remotes[remoteAddress] = conn
		s.mu.Unlock()

		go func() {
			s.device.ServeTCP(conn)

			s.mu.Lock()
			delete(s.remotes, remoteAddress)
			s.mu.Unlock()
		}()
	}
}
