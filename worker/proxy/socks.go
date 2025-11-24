package proxy

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/pubudu2003060/go-proxy-prototype/worker/auth"
	"github.com/pubudu2003060/go-proxy-prototype/worker/config"
	"github.com/pubudu2003060/go-proxy-prototype/worker/usage"
)

var PROXY_VERSION = byte(0x05)
var Uandp = byte(0x02)

const CMD_CONNECT = byte(0x01)

type SocksProxy struct {
	ConfigManager *config.ConfigManager
	authClient    *auth.AuthClient
	usageRepoter  *usage.UsageRepoter
	sessionMap    map[string]string
	mu            sync.RWMutex
}

func NewSocksProxy(configManager *config.ConfigManager, authClient *auth.AuthClient, usageRepoter *usage.UsageRepoter) *SocksProxy {
	return &SocksProxy{
		ConfigManager: configManager,
		authClient:    authClient,
		usageRepoter:  usageRepoter,
		sessionMap:    make(map[string]string),
	}
}

func (s *SocksProxy) HandleConnection(client net.Conn) {
	defer client.Close()

	//authentication
	if err := s.authHandShake(client); err != nil {
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

func (s *SocksProxy) authHandShake(client io.ReadWriter) error {
	header := make([]byte, 2)
	if _, err := io.ReadFull(client, header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	if header[0] != PROXY_VERSION {
		return fmt.Errorf("unsupported SOCKS5 version: %d", header[0])
	}

	nMethods := int(header[1])
	methods := make([]byte, nMethods)
	if _, err := io.ReadFull(client, methods); err != nil {
		return fmt.Errorf("failed to read methods: %w", err)
	}

	hasAuth := false
	for _, method := range methods {
		if method == Uandp {
			hasAuth = true
			break
		}
	}

	if !hasAuth {
		return fmt.Errorf("client is not suported username and pasword for socks5")
	}

	if _, err := client.Write([]byte{PROXY_VERSION, Uandp}); err != nil {
		return err
	}

	upHeader := make([]byte, 2)
	if _, err := io.ReadFull(client, upHeader); err != nil {
		return fmt.Errorf("failed to read auth header: %w", err)
	}

	ulen := upHeader[1]
	uname := make([]byte, ulen)
	if _, err := io.ReadFull(client, uname); err != nil {
		return fmt.Errorf("failed to read username in socks5: %w", err)
	}

	plen := make([]byte, 1)
	if _, err := io.ReadFull(client, plen); err != nil {
		return fmt.Errorf("failed to read username in socks5: %w", err)
	}

	password := make([]byte, plen[0])
	if _, err := io.ReadFull(client, password); err != nil {
		return fmt.Errorf("failed to read username in socks5: %w", err)
	}

	_, err := s.authClient.Authenticate(string(uname), string(password))
	if err != nil {
		client.Write([]byte{0x05, 0x01})
		return fmt.Errorf("auth credentials worng in socks5")
	}

	client.Write([]byte{0x05, 0x00})

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
