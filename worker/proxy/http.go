package proxy

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/pubudu2003060/go-proxy-prototype/worker/auth"
	"github.com/pubudu2003060/go-proxy-prototype/worker/config"
	"github.com/pubudu2003060/go-proxy-prototype/worker/models"
	"github.com/pubudu2003060/go-proxy-prototype/worker/usage"
)

type HTTPProxy struct {
	ConfigManager *config.ConfigManager
	authClient    *auth.AuthClient
	usageRepoter  *usage.UsageRepoter
	sessionMap    map[string]string
	mu            sync.RWMutex
}

func NewHTTPProxy(configManager *config.ConfigManager, authClient *auth.AuthClient, usageRepoter *usage.UsageRepoter) *HTTPProxy {
	return &HTTPProxy{
		ConfigManager: configManager,
		authClient:    authClient,
		usageRepoter:  usageRepoter,
		sessionMap:    make(map[string]string),
	}
}

func (p *HTTPProxy) HandleConnection(w http.ResponseWriter, r *http.Request) {
	log.Printf("new request come:%v", r.Host)
	proxyAuth := r.Header.Get("Proxy-Authorization")
	authresp, filters, err := p.authenticateProxyHeader(proxyAuth)
	if err != nil {
		log.Printf("auth failed: %v", err)
		p.send407(w)
		return
	}

	if authresp.DataUsed >= authresp.DataLimit {
		p.send429(w)
		return
	}

	pools := p.ConfigManager.GetPools()

	var selectedPool *models.Pool
	for _, poolName := range authresp.AllowedPools {
		name := strings.TrimSpace(poolName)
		name = strings.ToLower(name)
		if pool, exit := pools[name]; exit {
			selectedPool = pool
			break
		}
	}

	if selectedPool == nil {
		log.Println("there is no allowd pools.directly connect to destination")
		p.send500(w)
		return
	}

	upstream := p.selectUpstream(selectedPool, filters)
	if upstream == nil {
		log.Println("there is no allowd upstreams.directly connect to destination")
		p.send500(w)
		return
	}

	r.Header.Del("Proxy-Authorization")

	i := strings.Index(upstream.Format, "-")
	credentials := upstream.Format[:i] + filters

	upStreamURL := "http://" + credentials + "@" + upstream.Domain + ":" + strconv.Itoa(upstream.UpstreamPort)

	u, err := url.Parse(upStreamURL)
	if err != nil {
		log.Fatal("Invalid upstream proxy: ", err)
		p.send500(w)
		return
	}

	if r.Method == http.MethodConnect {
		p.handleConnect(w, r, u)
		return
	}
	p.handleHTTP(w, r, u)
}

func (p *HTTPProxy) handleHTTP(w http.ResponseWriter, r *http.Request, u *url.URL) {
	log.Println("HTTP request:", r.URL.String())

	var transport *http.Transport

	if u == nil {
		transport = &http.Transport{}
	} else {
		transport = &http.Transport{
			Proxy: http.ProxyURL(u),
		}
	}

	client := &http.Client{Transport: transport}

	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header = r.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (p *HTTPProxy) handleConnect(w http.ResponseWriter, r *http.Request, u *url.URL) {
	log.Println("HTTPS request:", r.Host)

	if u == nil {
		destConn, err := net.Dial("tcp", r.Host)
		if err != nil {
			http.Error(w, "Cannot connect to target", http.StatusBadGateway)
			return
		}

		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
			return
		}
		clientConn, _, err := hj.Hijack()
		if err != nil {
			http.Error(w, "Hijack failed", http.StatusInternalServerError)
			return
		}

		clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

		go io.Copy(destConn, clientConn)
		go io.Copy(clientConn, destConn)
		return
	}

	upstreamConn, err := net.Dial("tcp", u.Host)
	if err != nil {
		http.Error(w, "Cannot reach upstream", http.StatusBadGateway)
		return
	}

	connectReq := "CONNECT " + r.Host + " HTTP/1.1\r\n"
	if u.User != nil {
		user := u.User.String()
		auth := "Proxy-Authorization: Basic " + base64.StdEncoding.EncodeToString([]byte(user)) + "\r\n"
		connectReq += auth
	}
	connectReq += "\r\n"

	_, err = upstreamConn.Write([]byte(connectReq))
	if err != nil {
		http.Error(w, "Failed to send CONNECT", http.StatusBadGateway)
		return
	}

	buf := make([]byte, 4096)
	n, _ := upstreamConn.Read(buf)
	log.Printf("Upstream replied: %s", string(buf[:n]))

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, "Hijack failed", http.StatusInternalServerError)
		return
	}

	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	go io.Copy(upstreamConn, clientConn)
	go io.Copy(clientConn, upstreamConn)
}

func (p *HTTPProxy) selectUpstream(pool *models.Pool, filters string) *models.Out {
	var sessionID string

	if strings.Contains(filters, "session") {
		i := strings.Index(filters, "session")
		sessionID = filters[i+8 : i+15]
	} else if strings.Contains(filters, "sid") {
		i := strings.Index(filters, "sid")
		sessionID = filters[i+4 : i+12]
	}

	p.mu.RLock()
	upstreamKey, exists := p.sessionMap[sessionID]
	p.mu.RUnlock()

	if len(pool.Outs) == 0 {
		return nil
	}

	if exists {
		for _, out := range pool.Outs {
			if out.Domain == upstreamKey {
				return &out
			}
		}
	}

	selected := &pool.Outs[0]

	p.mu.Lock()
	p.sessionMap[sessionID] = selected.Domain
	p.mu.Unlock()

	return selected
}

func (p *HTTPProxy) authenticateProxyHeader(header string) (*models.AuthResponse, string, error) {
	if header == "" {
		return nil, "", errors.New("missing proxy-authorization")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return nil, "", errors.New("malformed proxy-authorization")
	}
	scheme := strings.ToLower(parts[0])
	cred := strings.TrimSpace(parts[1])

	switch scheme {
	case "basic":
		decoded, err := base64.StdEncoding.DecodeString(cred)
		if err != nil {
			return nil, "", fmt.Errorf("bad basic encoding: %w", err)
		}
		up := string(decoded)
		position := strings.IndexByte(up, ':')
		filterPosition := strings.IndexByte(up, '-')
		if position < 0 {
			return nil, "", errors.New("basic credentials missing ':'")
		}
		username, password, cred := up[:position], up[position+1:filterPosition], up[filterPosition:]
		authresp, err := p.authClient.Authenticate(username, password)
		if err != nil {
			log.Println("Error authenticate", err)
			return nil, "", errors.New("invalid user/pass")
		}
		return authresp, cred, nil
	default:
		return nil, "", fmt.Errorf("unsupported auth scheme: %s", scheme)
	}
}

func (p *HTTPProxy) send407(w http.ResponseWriter) {
	w.Header().Set("Proxy-Authenticate", `Basic realm="Proxy", Bearer`)
	w.WriteHeader(http.StatusProxyAuthRequired)
	fmt.Fprintln(w, "Proxy Authentication Required")
}

func (p *HTTPProxy) send429(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprintln(w, "too many request")
}

func (p *HTTPProxy) send500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "process faild")
}
