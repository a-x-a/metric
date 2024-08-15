package sender

import (
	"errors"
	"net"
)

var ErrOutboundIP = errors.New("не удалось получить основной ip адрес")

// getOutboundIP получет основной исходящий адрес компьютера.
func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.IP{}, err
	}
	defer conn.Close()

	if localAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return localAddr.IP, nil
	}

	return net.IP{}, ErrOutboundIP
}
