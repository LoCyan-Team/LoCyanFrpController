package main

import (
	"LoCyanFrpController/net/server"
	"LoCyanFrpController/pkg/process"
	websocket2 "LoCyanFrpController/pkg/websocket"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func main() {
	ws := websocket2.NewWebSocket()
	err := ws.ConnectWsServer()
	if err != nil {
		log.Fatalf("Can't connect to WebSocket server, err: %v", err)
	} else {
		log.Printf("Connect to WebSocket server successfully!")
		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {
				log.Fatalf("Can't close WebSocket Connection, err: %v", err)
			}
		}(ws.Conn)
		go ws.ReadMsg()
		ticker := time.NewTicker(time.Second * ws.Config.SendDuration)
		defer ticker.Stop()

		// 监听 opengfw 日志
		err := process.HookServiceLogs("opengfw")
		if err != nil {
			log.Fatalf("Can't hook service: opengfw, err: %s", err)
		}

		serverInfo := server.GetServerInfo()
		ws.SendNodeStatsToServer(serverInfo)
		ws.SendProxyStatsToServer()

		for range ticker.C {
			serverInfo := server.GetServerInfo()
			ws.SendNodeStatsToServer(serverInfo)
			ws.SendProxyStatsToServer()
		}
	}
}
