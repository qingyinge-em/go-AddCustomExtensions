package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	hc := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				// Rand:               &RandGen{0x12},
				CustomExtensionGen: func(rand []byte) (uint16, []byte, error) {
					// logger.Infof("rand: %x", rand)
					h := hmac.New(sha1.New, rand[len(rand)/2:])
					h.Write(rand[:len(rand)/2])
					h.Write([]byte("xxx"))
					return 0xffff, h.Sum(nil), nil
				},
			},
			Dial: func(network, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(network, addr, time.Second*10)
				if err != nil {
					// logger.Error("DialTimeout err:", err)
					return nil, err
				}
				return c, nil
			},
		},
	}

	rsp, err := hc.Get("https://localhost")
	if err != nil {
		// logger.Error(err)
		return
	}
	defer rsp.Body.Close()

	_, err = io.ReadAll(rsp.Body)
	if err != nil {
		// logger.Error(err)
		return
	}

	// logger.Info(string(rspdata[:128]))
}
