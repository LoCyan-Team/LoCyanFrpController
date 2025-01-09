package info

import (
	_type "LoCyanFrpController/pkg/type"
)

func GetNodeInfo(info _type.FrpsServerInfoResponse) (data map[string]any) {
	data = make(map[string]any)
	inbound := info.TotalTrafficIn
	outbound := info.TotalTrafficOut
	clientCount := info.ClientCounts
	data["inbound"] = inbound
	data["outbound"] = outbound
	data["client_count"] = clientCount
	return data
}
