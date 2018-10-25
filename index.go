package sensu_check_push

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type CheckResult struct {
	Handlers []string `json:"handlers"`
	Name     string   `json:"name"`
	Output   string   `json:"output"`
	Status   int      `json:"status"`
	Ttl      int      `json:"ttl"`
}

var DefaultTimeout = time.Duration(5)
var DefaultHost = "localhost"
var DefaultPort = 3030

func PushCheckResultsToSensu(host string, port int, result CheckResult, timeoutSeconds time.Duration) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeoutSeconds*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Use a map to be able to omit unneeded fields
	resultMap := make(map[string]interface{})
	resultMap["handlers"] = result.Handlers
	resultMap["name"] = result.Name
	resultMap["output"] = result.Output
	resultMap["status"] = result.Status
	if result.Ttl > 0 {
		resultMap["ttl"] = result.Ttl
	}

	bytes, err := json.Marshal(resultMap)
	if err != nil {
		return err
	}

	conn.Write(bytes)
	return nil
}
