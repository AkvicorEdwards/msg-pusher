package app

import (
	"net"
	"net/http"
	"strings"
)

func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown"
	}

	if net.ParseIP(ip) != nil {
		return ip
	}

	return "unknown"
}

func LastPage(w http.ResponseWriter, r *http.Request) {
	RespRedirect(w, r, r.URL.String()[:strings.LastIndexByte(r.URL.String(), '/')])
}

func Reload(w http.ResponseWriter, r *http.Request) {
	RespRedirect(w, r, r.URL.String())
}
