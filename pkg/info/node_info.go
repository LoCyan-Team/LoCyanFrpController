package info

import (
	_frps "lcf-controller/pkg/type/frps"
)

func GetNodeInfo(info _frps.ServerInfoResponse) (data map[string]any) {
	data = make(map[string]any)
	inbound := info.TotalTrafficIn
	outbound := info.TotalTrafficOut
	clientCount := info.ClientCounts
	data["inbound"] = inbound
	data["outbound"] = outbound
	data["client_count"] = clientCount
	return data
}
