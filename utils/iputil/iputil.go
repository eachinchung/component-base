package iputil

//goland:noinspection SpellCheckingInspection
import (
	stdnet "net"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// http headers.
const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
	XClientIP     = "x-client-ip"
)

// GetLocalIP 返回主机的本地IP。
func GetLocalIP() string {
	address, err := stdnet.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range address {
		if ipNet, ok := addr.(*stdnet.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// GetRemoteIP 返回请求的远程ip
func GetRemoteIP(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XClientIP); ip != "" {
		remoteAddr = ip
	} else if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = stdnet.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

// GetLocalExternalIPV4 获取本机外网IP
func GetLocalExternalIPV4() (string, error) {
	client := resty.New()
	resp, err := client.R().Get("https://ipv4.icanhazip.com")
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}
