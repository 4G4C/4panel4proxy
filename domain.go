package main

import (
	"strings"

	"github.com/kataras/iris/v12"
)

type Domain struct {
	Type         string `json:"type"`
	Domain       string `json:"domain"`
	ReverseProxy string `json:"reverse_proxy"`
	ForceHTTPS   bool   `json:"force_https"`
	TLSCert      string `json:"tls_cert"`
	TLSKey       string `json:"tls_key"`
}

func GetDomain(ctx iris.Context) (Domain, bool) {
	host := strings.Split(ctx.Host(), ":")[0]
	var domain Domain
	var ok bool
	if strings.Contains(host, domainMain) {
		hostSub := strings.Split(host, ".")[0]
		if strings.Contains(hostSub, "vpn") {
			hostPort := strings.Split(strings.ReplaceAll(hostSub, "vpn", ""), "-")
			port := "80"
			if len(hostPort) > 1 {
				port = hostPort[1]
			}
			domain.Domain = host
			domain.Type = "reverse_proxy"
			domain.ReverseProxy = "10.8.0." + hostPort[0] + ":" + port
			return domain, true
		}
	}
	domain, ok = domainList[host]
	return domain, ok
}
