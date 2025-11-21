package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pubudu2003060/go-proxy-prototype/worker/auth"
	"github.com/pubudu2003060/go-proxy-prototype/worker/config"
	"github.com/pubudu2003060/go-proxy-prototype/worker/proxy"
	"github.com/pubudu2003060/go-proxy-prototype/worker/usage"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go startHTTPProxy(&wg)
	go startSOCKSProxy(&wg)
	wg.Wait()
}

func startHTTPProxy(wg *sync.WaitGroup) {
	addr := ":8081"

	configManager := config.NewConfigManager("http://localhost:8080")

	authClient := auth.NewAuthClient("http://localhost:8080")

	usageReporter := usage.NewUsageReporter("http://localhost:8080")

	proxy := proxy.NewHTTPProxy(configManager, authClient, usageReporter)

	go configManager.StartSync(30 * time.Second)

	srv := &http.Server{
		Addr:         addr,
		Handler:      http.HandlerFunc(proxy.HandleConnection),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("HTTP/S proxy listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		wg.Done()
		return
	}
}

var PROXY_VERSION = byte(0x05)
var NO_AUTH = byte(0x00)

const CMD_CONNECT = byte(0x01)

func startSOCKSProxy(wg *sync.WaitGroup) {
	listner, err := net.Listen("tcp", ":1080")
	if err != nil {
		wg.Done()
		log.Fatalf("failed to listen :1080: %s", err)
	}
	defer listner.Close()
	log.Println("SOCKS5 proxy listening on :1080")

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}
		handleConnection(conn)
	}
}

func handleConnection(client net.Conn) {
	defer client.Close()

	//authentication
	if err := authHandShake(client); err != nil {
		log.Printf("Authentication failed: %s", err)
		return
	}

	//socks5 request
	destAddr, err := requestHandShake(client)
	if err != nil {
		log.Printf("Request handshake failed: %s", err)
		return
	}

	// 3. Connect to the destination
	dest, err := net.Dial("tcp", destAddr)
	if err != nil {
		log.Printf("Failed to connect to destination %s: %v", destAddr, err)
		sendReply(client, 0x04, nil)
		return
	}
	defer dest.Close()

	if err := sendReply(client, 0x00, dest.LocalAddr()); err != nil {
		log.Printf("Failed to send success reply: %v", err)
		return
	}

	// 4. Tunnel the data
	log.Printf("Tunneling data for %s", destAddr)
	go io.Copy(dest, client)
	io.Copy(client, dest)
	log.Printf("Connection closed for %s", destAddr)
}

func authHandShake(client io.ReadWriter) error {
	header := make([]byte, 2)
	if _, err := io.ReadFull(client, header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	if header[0] != PROXY_VERSION {
		return fmt.Errorf("unsupported SOCKS version: %d", header[0])
	}

	nMethods := int(header[1])
	methods := make([]byte, nMethods)
	if _, err := io.ReadFull(client, methods); err != nil {
		return fmt.Errorf("failed to read methods: %w", err)
	}

	hasNoAuth := false
	for _, method := range methods {
		if method == NO_AUTH {
			hasNoAuth = true
			break
		}
	}

	if !hasNoAuth {
		client.Write([]byte{PROXY_VERSION, 0xFF})
	}

	if _, err := client.Write([]byte{PROXY_VERSION, NO_AUTH}); err != nil {
		return err
	}

	return nil
}

func requestHandShake(client io.Reader) (string, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(client, header); err != nil {
		return "", err
	}

	if header[0] != PROXY_VERSION {
		return "", fmt.Errorf("unsupported SOCKS version: %d", header[0])
	}

	if header[1] != CMD_CONNECT {
		return "", fmt.Errorf("unsupported command: %d", header[1])
	}

	var host string
	addrType := header[3]

	switch addrType {
	case 0x01: // IPv4
		ip := make([]byte, 4)
		if _, err := io.ReadFull(client, ip); err != nil {
			return "", fmt.Errorf("reading IPv4 address: %w", err)
		}
		host = net.IP(ip).String()
	case 0x03: // Domain Name
		lenBuf := make([]byte, 1)
		if _, err := io.ReadFull(client, lenBuf); err != nil {
			return "", fmt.Errorf("reading domain length: %w", err)
		}
		domainLen := int(lenBuf[0])
		domain := make([]byte, domainLen)
		if _, err := io.ReadFull(client, domain); err != nil {
			return "", fmt.Errorf("reading domain: %w", err)
		}
		host = string(domain)
	case 0x04: // IPv6
		ip := make([]byte, 16)
		if _, err := io.ReadFull(client, ip); err != nil {
			return "", fmt.Errorf("reading IPv6 address: %w", err)
		}
		host = net.IP(ip).String()
	default:
		return "", fmt.Errorf("unknown address type: %d", addrType)
	}

	// Read Port (2 bytes, big-endian)
	portBuf := make([]byte, 2)
	if _, err := io.ReadFull(client, portBuf); err != nil {
		return "", fmt.Errorf("reading port: %w", err)
	}
	port := binary.BigEndian.Uint16(portBuf)

	return net.JoinHostPort(host, strconv.Itoa(int(port))), nil
}

func sendReply(client io.Writer, rep byte, addr net.Addr) error {
	// [VER | REP | RSV | ATYP | BND.ADDR | BND.PORT]
	reply := []byte{PROXY_VERSION, rep, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	_, err := client.Write(reply)
	return err
}
