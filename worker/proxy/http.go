package proxy

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pubudu2003060/go-proxy-prototype/worker/auth"
	"github.com/pubudu2003060/go-proxy-prototype/worker/config"
	"github.com/pubudu2003060/go-proxy-prototype/worker/models"
	"github.com/pubudu2003060/go-proxy-prototype/worker/usage"
)

type HTTPProxy struct {
	ConfigManager *config.ConfigManager
	authClient *auth.AuthClient
	usageRepoter *usage.UsageRepoter
	sessionMap map[string]string
	mu sync.RWMutex
}

func NewHTTPProxy(configManager *config.ConfigManager,authClient *auth.AuthClient,usageRepoter *usage.UsageRepoter) *HTTPProxy{
	return &HTTPProxy{
		ConfigManager: configManager,
		authClient: authClient,
		usageRepoter: usageRepoter,
		sessionMap: make(map[string]string),
	}
}

func (p *HTTPProxy) Start(addr string) error {
	ln,err := net.Listen("tcp",addr)
	if err != nil {
		return err
	}

	for {
		conn,err := ln.Accept()
		if err != nil {
			log.Printf("accept error %v,",err)
		}
		go p.handleConnection(conn)
	}
}

func (p *HTTPProxy) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte,4096)
	n,err := conn.Read(buf)
	if err != nil {
		return
	}

	reader := bufio.NewReader(bytes.NewReader(buf[:n]))
	req,err := http.ReadRequest(reader)
	if err != nil {
		return
	}

	username,password,hasAuth := p.extractCredentials(req)
	if !hasAuth{
		p.sendAuthRequired(conn)
		return
	}

	authresp,err := p.authClient.Authenticate(username,password)
	if err != nil || authresp.Success {
		p.sendAuthFailed(conn)
		return
	}

	if authresp.DataLimit >= authresp.DataLimit {
		p.sendDataLimitExceeded(conn)
		return
	}

	sessionId := p.extractSessionID(req)
	if sessionId == "" {
		sessionId = generateSessionID()
	}

	p.routeRequest(conn,buf[:n],authresp,sessionId)

}

func (p *HTTPProxy) routeRequest(conn net.Conn,initialData []byte,authResp *models.AuthResponse,sessionId string){
	pools := p.ConfigManager.GetPools()

	var selectedPool *models.Pool
	for _,poolName := range authResp.AllowedPools {
		if pool,exit := pools[poolName];exit {
			selectedPool = pool
			break
		}
	}

	if selectedPool == nil {
		p.sendNoAllowedPools(conn)
		return
	}

	// Select upstream with session stickiness
	upstream := p.selectUpstream(selectedPool, sessionId)
	if upstream == nil {
		p.sendNoUpstreamAvailable(conn)
		return
	}
	
	// Connect to upstream
	targetConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", upstream.Domain, upstream.UpstreamPort))
	if err != nil {
		log.Printf("Failed to connect to upstream: %v", err)
		return
	}
	defer targetConn.Close()
	
	// For CONNECT method (HTTPS)
	if string(initialData[:7]) == "CONNECT" {
		conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	} else {
		// For HTTP, write the initial data
		targetConn.Write(initialData)
	}
	
	// Track data transfer
	var bytesSent, bytesReceived int64
	
	// Copy data between client and upstream
	go func() {
		defer conn.Close()
		defer targetConn.Close()
		
		n, _ := io.Copy(conn, targetConn)
		bytesReceived += n
	}()
	
	n, _ := io.Copy(targetConn, conn)
	bytesSent += n
	
	// Report usage
	totalBytes := bytesSent + bytesReceived
	if totalBytes > 0 {
		p.usageRepoter.ReportUsage(authResp.UserID, totalBytes)
	}
}

func (p *HTTPProxy) selectUpstream(pool *models.Pool, sessionID string) *models.Out {
	p.mu.RLock()
	upstreamKey, exists := p.sessionMap[sessionID]
	p.mu.RUnlock()
	
	if exists {
		for _, out := range pool.Outs {
			// Simple key matching - in real implementation, use proper keys
			if out.Domain == upstreamKey {
				return &out
			}
		}
	}
	
	// Weighted random selection for new session
	if len(pool.Outs) == 0 {
		return nil
	}
	
	selected := &pool.Outs[0] // Simple selection for prototype
	
	p.mu.Lock()
	p.sessionMap[sessionID] = selected.Domain
	p.mu.Unlock()
	
	return selected
}

func (p *HTTPProxy) extractCredentials(req *http.Request) (string,string,bool){
	auth := req.Header.Get("proxy-Authorization")
	if auth == ""{
		return "","",false
	}

	if !strings.HasPrefix(auth,"Basic ") {
		return "","",false
	}

	payload,err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return "","",false
	}

	pair := strings.SplitN(string(payload),":",2)
	if len(pair) != 2 {
		return  "","",false
	}

	return pair[0],pair[1],true

}

func (p *HTTPProxy) extractSessionID(req *http.Request) string {
	auth := req.Header.Get("proxy-authentication")
	if auth != ""{
		return auth
	}
	return ""
}

func (p *HTTPProxy) sendAuthRequired(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 407 Proxy Authentication Required\r\nProxy-Authenticate: Basic realm=\"Upstream Y\"\r\n\r\n"))
}

func (p *HTTPProxy) sendAuthFailed(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 407 Proxy Authentication Required\r\n\r\nAuthentication failed"))
}

func (p *HTTPProxy) sendDataLimitExceeded(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 507 Insufficient Storage\r\n\r\nData limit exceeded"))
}

func generateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (p *HTTPProxy) sendNoAllowedPools(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 403 Forbidden\r\n\r\nNo allowed pools available"))
}

func (p *HTTPProxy) sendNoUpstreamAvailable(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\n\r\nNo upstream available"))
}
