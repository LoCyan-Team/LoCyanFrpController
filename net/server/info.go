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

func sendRequest(url string) (*http.Response, []byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	req.SetBasicAuth(getBasicAuthInfo())
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	var response bytes.Buffer
	if _, err := io.Copy(&response, resp.Body); err != nil {
		return nil, nil, err
	}

	return resp, response.Bytes(), nil
}

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

	resp, body, err := sendRequest(url)
	if err != nil {
		return frps.ServerInfoResponse{}, err
	}
	defer resp.Body.Close()

	var serverInfo frps.ServerInfoResponse
	err = json.Unmarshal(body, &serverInfo)
	if err != nil {
		log.Fatalf("error unmarshalling JSON: %v", err)
	}
	return serverInfo, nil
}

func GetProxyList(proxyType string) (frps.Proxy, error) {
	url := fmt.Sprintf(
		"%s/%s",
		getUrl("/proxy"),
		proxyType,
	)

	resp, body, err := sendRequest(url)
	if err != nil {
		return frps.Proxy{}, err
	}
	defer resp.Body.Close()

	var proxyInfo frps.Proxy
	err = json.Unmarshal(body, &proxyInfo)
	if err != nil {
		return frps.Proxy{}, err
	}
	return proxyInfo, nil
}
