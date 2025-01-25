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

func GetServerInfo() frps.ServerInfoResponse {
	configInfo := config.ReadCfg()
	username := configInfo.Username
	password := configInfo.Password
	adminPort := configInfo.AdminPort

	url := fmt.Sprintf("http://127.0.0.1:%s/api/serverinfo", strconv.FormatInt(int64(adminPort), 10))

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing response body: %v", err)
		}
	}(resp.Body)

	var response bytes.Buffer
	if _, err := io.Copy(&response, resp.Body); err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var serverInfo frps.ServerInfoResponse
	body := response.Bytes()
	err = json.Unmarshal(body, &serverInfo)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	return serverInfo
}

func GetProxyList(proxyType string) frps.Proxy {
	configInfo := config.ReadCfg()
	username := configInfo.Username
	password := configInfo.Password
	adminPort := configInfo.AdminPort

	url := fmt.Sprintf("http://127.0.0.1:%s/api/proxy/%s", strconv.FormatInt(int64(adminPort), 10), proxyType)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing response body: %v", err)
		}
	}(resp.Body)

	var response bytes.Buffer
	if _, err := io.Copy(&response, resp.Body); err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var proxyInfo frps.Proxy
	body := response.Bytes()
	err = json.Unmarshal(body, &proxyInfo)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	return proxyInfo
}
