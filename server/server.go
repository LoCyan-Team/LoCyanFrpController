package server

import (
	"fmt"
	"lcf-controller/logger"
	"lcf-controller/net/api"
	"lcf-controller/pkg/config"
	"lcf-controller/pkg/info"
)

func SendTunnelTrafficToServer(cfg *config.Config) (err error) {
	proxyType := []string{"tcp", "udp", "http", "https", "xtcp", "stcp", "tcpmux", "sudp"}
	header := make(map[string]any)
	header["X-Node-API-Key"] = cfg.ControllerConfig.NodeApiKey
	for _, pxyType := range proxyType {
		if pxyInfoList, err := info.GetTunnelInfo(cfg, pxyType); err != nil {
			return err
		} else {
			for _, pxyInfo := range pxyInfoList {
				if _, err := api.HttpPost(cfg.ControllerConfig.Endpoint+"/node/exchange/traffic", pxyInfo, header); err != nil {
					return err
				} else {
					logger.Info(fmt.Sprintf("update proxy traffic: %s, proxy_type: %s, ↑%dB ↓%dB", pxyInfo["proxy_name"], pxyType, pxyInfo["outbound_traffic"], pxyInfo["inbound_traffic"]))
				}
			}
		}
	}
	return nil
}
