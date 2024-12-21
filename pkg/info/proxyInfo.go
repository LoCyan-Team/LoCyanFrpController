package info

import (
	"LoCyanFrpController/net/server"
	"strings"
)

func GetProxies(proxyType string) (listMap []map[string]any) {
	proxies := server.GetProxyList(proxyType)
	proxyList := make([]map[string]any, 0)
	for _, p := range proxies.Proxies {
		tmp := make(map[string]any)
		tmp["proxy_name"] = strings.Split(p.Name, "0")[1]
		tmp["inbound"] = p.TodayTrafficIn
		tmp["outbound"] = p.TodayTrafficOut
		proxyList = append(proxyList, tmp)
	}
	return proxyList
}
