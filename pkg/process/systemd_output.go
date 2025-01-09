package process

import (
	"LoCyanFrpController/pkg/database"
	log2 "LoCyanFrpController/pkg/log"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

func HookServiceLogs(serviceName string) (err error) {
	// 持续 hook 日志
	for {
		cmd := exec.Command("journalctl", "-u", serviceName, "--follow", "--output", "cat")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}

		if err := cmd.Start(); err != nil {
			fmt.Println("Error starting command:", err)
			time.Sleep(5 * time.Second) // 等待一段时间后重试
			continue
		}

		processLogs(stdout)

		if err := cmd.Wait(); err != nil {
			return err
		}

		fmt.Println("Restarting journalctl command...")
		time.Sleep(5 * time.Second)
	}
}

func processLogs(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		parts := log2.SplitLog(line)
		log.Printf(line)
		logType := parts[2]
		jsonText := parts[3]

		// 将每一条记录存入时序库, 包含 消息 ID, src, dst, 规则名, action
		// 后端可以根据 dst 获取到远程端口, 即可定位隧道, 若触犯封禁规则, 同时在数据库内存入
		// 数据库新增表存储规则列表, 并存储触犯该规则是否进行封禁

		db, err := database.InitDb()
		if err != nil {
			log.Fatalf("Error while initing database connection, err: %v", err)
		}
		switch logType {
		case "TCPstreamaction":
			var data TCPStreamAction
			if err := json.Unmarshal([]byte(jsonText), &data); err != nil {
				id := data.ID
				dst := data.Dst
				src := data.Src
				action := data.Action
				query := fmt.Sprintf("SELECT * FROM `opengfw_record` WHERE `dst`=%q", dst)
				rs, err := db.QueryRecord(query)
				if err != nil {
					log.Fatalf("Error while insert data, err: %v", err)
				}
				if len(rs) != 0 {
					// 这个远程端口已经封过了
					continue
				}
				query = fmt.Sprintf("SELECT * FROM `opengfw_record` WHERE `id`=%q", id)
				rs, err = db.QueryRecord(query)
				if err != nil {
					log.Fatalf("Error while insert data, err: %v", err)
				}
				if len(rs) != 0 {
					// 已经有记录了，更新
					query = fmt.Sprintf("UPDATE `opengfw_record` SET `dst` = %s, `src` = %s, `action` = %s WHERE `id` = %s", dst, src, action, id)
					err = db.UpdateRecord(query)
					if err != nil {
						log.Fatalf("Error while updating data, err: %v", err)
					}
				} else {
					// 没有记录，插入
					query = fmt.Sprintf("INSERT INTO `opengfw_record` (id, rule_name, dst, src, time, action) VALUES (%s, '', %s, %s, NOW(), %s)", id, dst, src, action)
					err = db.InsertRecord(query)
					if err != nil {
						log.Fatalf("Error while inserting data, err: %v", err)
					}
				}
			} else {
				log.Fatalf("Error while read OpenGFW's output, err: %v", err)
			}
		case "rulesetlog":
			var data RulesetLog
			if err := json.Unmarshal([]byte(jsonText), &data); err != nil {
				id := data.ID
				dst := data.Dst
				src := data.Src
				ruleName := data.Name
				query := fmt.Sprintf("SELECT * FROM `opengfw_record` WHERE `dst`=%q", dst)
				rs, err := db.QueryRecord(query)
				if err != nil {
					log.Fatalf("Error while insert data, err: %v", err)
				}
				if len(rs) != 0 {
					// 这个远程端口已经封过了
					continue
				}
				query = fmt.Sprintf("SELECT * FROM `opengfw_record` WHERE `id`=%q", id)
				rs, err = db.QueryRecord(query)
				if err != nil {
					log.Fatalf("Error querying data, err: %v", err)
				}
				if len(rs) != 0 {
					// 已经有记录了，更新
					query = fmt.Sprintf("UPDATE `opengfw_record` SET `rule_name` = %q WHERE `id` = %q", ruleName, id)
					err = db.UpdateRecord(query)
					if err != nil {
						log.Fatalf("Error while updating rule_name, err: %v", err)
					}
				} else {
					// 没有记录，插入
					query = fmt.Sprintf("INSERT INTO `opengfw_record` (id, rule_name, dst, src, time, action) VALUES (%q, %q, %q, %q, NOW(), '')", id, ruleName, dst, src)
					err = db.InsertRecord(query)
					if err != nil {
						log.Fatalf("Error while inserting data, err: %v", err)
					}
				}
			} else {
				log.Fatalf("Error while read OpenGFW's output, err: %v", err)
			}
		default:
			log.Fatalf("Unknown type: %s, text: %s", logType, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading log line:", err)
	}
}

type TCPStreamAction struct {
	ID     string `json:"id"`
	Src    string `json:"src"`
	Dst    string `json:"dst"`
	Action string `json:"action"`
}

// HTTPReqHeaders 定义了嵌套的 JSON 对象中 "http.req.headers" 字段的结构
type HTTPReqHeaders struct {
	Accept    string `json:"accept"`
	Host      string `json:"host"`
	UserAgent string `json:"user-agent"`
}

// HTTPReq 定义了嵌套的 JSON 对象中 "http.req" 字段的结构
type HTTPReq struct {
	Headers HTTPReqHeaders `json:"headers"`
	Method  string         `json:"method"`
	Path    string         `json:"path"`
	Version string         `json:"version"`
}

// HTTPProps 定义了嵌套的 JSON 对象中 "http" 字段的结构
type HTTPProps struct {
	Req HTTPReq `json:"req"`
}

type Props struct {
	Fet  map[string]interface{} `json:"fet"`
	HTTP HTTPProps              `json:"http"`
}

// RulesetLog 定义了主映射结构
type RulesetLog struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Src   string `json:"src"`
	Dst   string `json:"dst"`
	Props Props  `json:"props"`
}
