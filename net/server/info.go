package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lcf-controller/pkg/config"
	"lcf-controller/pkg/type/frps"
	"log"
	"net/http"
	"strconv"
)

var cfg = config.ReadCfg().FrpServerConfig

func getBasicAuthInfo() (string, string) {
	return cfg.Username, cfg.Password
}

func getUrl(path string) string {
	return fmt.Sprintf(
		"http://%s:%s/api%s",
		cfg.AdminApiHost,
		strconv.FormatInt(int64(cfg.AdminApiPort), 10),
		path,
	)
}

func GetServerInfo() (frps.ServerInfoResponse, error) {
	url := getUrl("/serverinfo")

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return frps.ServerInfoResponse{}, err
	}

	req.SetBasicAuth(getBasicAuthInfo())
	resp, err := client.Do(req)
	if err != nil {
		return frps.ServerInfoResponse{}, err
	}
	defer resp.Body.Close()

	var response bytes.Buffer
	if _, err := io.Copy(&response, resp.Body); err != nil {
		return frps.ServerInfoResponse{}, err
	}

	var serverInfo frps.ServerInfoResponse
	body := response.Bytes()
	err = json.Unmarshal(body, &serverInfo)
	if err != nil {
		log.Fatalf("error unmarshalling JSON: %v", err)
	}
	return serverInfo, nil
}

func GetTunnelList(tunnelType string) (frps.Tunnel, error) {
	url := fmt.Sprintf(
		"%s/%s",
		getUrl("/proxy"),
		tunnelType,
	)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return frps.Tunnel{}, err
	}

	req.SetBasicAuth(getBasicAuthInfo())
	resp, err := client.Do(req)
	if err != nil {
		return frps.Tunnel{}, err
	}
	defer resp.Body.Close()

	var response bytes.Buffer
	if _, err := io.Copy(&response, resp.Body); err != nil {
		return frps.Tunnel{}, err
	}

	var tunnelInfo frps.Tunnel
	body := response.Bytes()
	err = json.Unmarshal(body, &tunnelInfo)
	if err != nil {
		return frps.Tunnel{}, err
	}
	return tunnelInfo, nil
}
