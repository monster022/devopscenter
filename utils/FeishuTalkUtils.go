package utils

import (
	"bytes"
	"devopscenter/configuration"
	"fmt"
	"net/http"
)

func Alter(people, message string) int {
	var jsonStr = []byte(fmt.Sprintf(`
		{
			"msg_type": "text",
			"content": {
				"text": "%s \n @ %s"
			}
		}`, message, people))
	req, err := http.NewRequest("POST", configuration.Configs.FeishuTalk, bytes.NewBuffer(jsonStr))
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	return resp.StatusCode
}
