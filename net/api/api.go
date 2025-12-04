package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func HttpPost(urlString string, params map[string]any, header map[string]any) ([]byte, error) {
	formData := url.Values{}
	for key, value := range params {
		formData.Set(key, fmt.Sprintf("%v", value))
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlString, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	// 设置 Header
	for key, value := range header {
		req.Header.Set(key, fmt.Sprintf("%v", value))
	}

	req.Header.Set("User-Agent", "LoCyanFrpController/3.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	var response bytes.Buffer
	if _, err := io.Copy(&response, resp.Body); err != nil {
		return nil, err
	}
	body := response.Bytes()
	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("请求失败，状态码: %d, 返回内容: %s, ", resp.StatusCode, body)
	}

	return body, nil
}

func HttpGet(urlString string) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Get(urlString)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	var response bytes.Buffer
	if _, err := io.Copy(&response, resp.Body); err != nil {
		return nil, err
	}
	body := response.Bytes()
	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return body, nil
}
