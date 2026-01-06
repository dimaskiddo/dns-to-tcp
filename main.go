package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
)

// --- Constants for DNS Limits ---
const (
	MaxUDPSize = 4096
	MaxDNSSize = 8192
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, MaxUDPSize)
	},
}

// --- Types ---
type Bridge struct {
	remoteAddr string
	poolSize   int
	connPool   chan *connWrapper
}

type connWrapper struct {
	conn      net.Conn
	createdAt time.Time
}

func main() {
	lh := flag.String("lh", "127.0.0.1", "Local IP")
	lp := flag.Int("lp", 5353, "Local UDP Port")
	rh := flag.String("rh", "127.0.0.1", "Remote IP")
	rp := flag.Int("rp", 5353, "Remote TCP Port")
	ps := flag.Int("p", 1024, "Warm Pool TCP Connections")

	flag.Parse()

	if *rp == 0 {
		fmt.Println("Error: Remote port (-rp) is required")
		os.Exit(1)
	}

	bridge := &Bridge{
		remoteAddr: fmt.Sprintf("%s:%d", *rh, *rp),
		poolSize:   *ps,
		connPool:   make(chan *connWrapper, *ps),
	}

	fmt.Printf("Connection Pool Size: %d\n", *ps)
	bridge.listenUDP(fmt.Sprintf("%s:%d", *lh, *lp))
}

func (b *Bridge) listenUDP(addr string) {
	// Using ListenConfig to Set SO_REUSEPORT
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
				syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, 4*1024*1024)
			})
		},
	}

	pc, err := lc.ListenPacket(context.Background(), "udp", addr)
	if err != nil {
		log.Fatalf("Critical Error: %v", err)
	}
	defer pc.Close()

	for {
		buf := bufferPool.Get().([]byte)

		n, clientAddr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}

		go b.handleRequest(pc, clientAddr, buf[:n])
	}
}

func (b *Bridge) getConnection() (*connWrapper, error) {
	for {
		select {
		case cw := <-b.connPool:
			// Check Connection Max Age (30s)
			if time.Since(cw.createdAt) < 30*time.Second {
				return cw, nil
			}
			cw.conn.Close()

		default:
			// Create on-demand connection with Low Latency socket options
			d := net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 60 * time.Second,
			}

			conn, err := d.Dial("tcp", b.remoteAddr)
			if err != nil {
				return nil, err
			}

			// Apply TCP_NODELAY
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				tcpConn.SetNoDelay(true)
				tcpConn.SetKeepAlive(true)
			}

			return &connWrapper{conn: conn, createdAt: time.Now()}, nil
		}
	}
}

func (b *Bridge) handleRequest(pc net.PacketConn, addr net.Addr, data []byte) {
	defer bufferPool.Put(data[:cap(data)])

	// Attempt the request up to 2 times (Initial + 1 Retry)
	for attempt := range 2 {
		cw, err := b.getConnection()
		if err != nil {
			if attempt == 1 {
				return
			}

			continue
		}

		success := b.executeProxy(pc, addr, cw.conn, data)
		if success {
			// Release connection back to pool
			select {
			case b.connPool <- cw:

			default:
				cw.conn.Close()
			}

			return
		}

		// If failed, discard connection and retry
		cw.conn.Close()
	}
}

func (b *Bridge) executeProxy(pc net.PacketConn, clientAddr net.Addr, tcpConn net.Conn, data []byte) bool {
	tcpConn.SetDeadline(time.Now().Add(1500 * time.Millisecond))

	// 1. Send Length (Pre-packed struct equivalent) + Data
	lenBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBuf, uint16(len(data)))

	if _, err := tcpConn.Write(lenBuf); err != nil {
		return false
	}

	if _, err := tcpConn.Write(data); err != nil {
		return false
	}

	// 2. Read Length
	if _, err := io.ReadFull(tcpConn, lenBuf); err != nil {
		return false
	}
	respLen := binary.BigEndian.Uint16(lenBuf)

	if respLen > MaxDNSSize {
		return false
	}

	// 3. Read Response Data
	respData := make([]byte, respLen)
	if _, err := io.ReadFull(tcpConn, respData); err != nil {
		return false
	}

	// 4. Check for NXDOMAIN (3) or SERVFAIL (2)
	if len(respData) >= 4 {
		rcode := respData[3] & 0x0F
		if rcode == 2 || rcode == 3 {
			pc.WriteTo(respData, clientAddr)
			return true
		}
	}

	// 5. Success
	pc.WriteTo(respData, clientAddr)
	return true
}
