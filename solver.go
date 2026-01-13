package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"
)

var CHARSET = func() []byte {
	out := []byte{}
	for i := 32; i <= 126; i++ {
		if i == '\n' || i == '\r' || i == '\t' || i == ' ' {
			continue
		}
		out = append(out, byte(i))
	}
	return out
}()

func sha1hex(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

func suffixAt(n int, length int) string {
	base := len(CHARSET)
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		out[i] = CHARSET[n%base]
		n /= base
	}
	return string(out)
}

func tlsConnect() (*tls.Conn, error) {
	cert, err := tls.LoadX509KeyPair(CERT_FILE, CERT_FILE)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", HOST, PORT))
	if err != nil {
		return nil, err
	}

	return tls.Client(conn, cfg), nil
}

func solvePOW(auth string, difficulty int) string {
	workers := runtime.NumCPU()
	done := make(chan struct{})
	result := make(chan string, 1)

	prefix := strings.Repeat("0", difficulty)

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func(worker int) {
			defer wg.Done()

			i := worker
			for {
				select {
				case <-done:
					return
				default:
				}

				suffix := suffixAt(i, 12)
				i += workers

				if strings.HasPrefix(sha1hex(auth+suffix), prefix) {
					select {
					case result <- suffix:
						close(done)
					default:
					}
					return
				}
			}
		}(w)
	}

	found := <-result
	wg.Wait()
	return found
}
