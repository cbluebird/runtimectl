package util

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestJson(t *testing.T) {
	// Replace single quotes with double quotes to make it a valid JSON string
	configStr := "{\"appPorts\":[{\"name\":\"devbox-app-port\",\"port\":8080,\"protocol\":\"TCP\"}],\"ports\":[{\"containerPort\":22,\"name\":\"devbox-ssh-port\",\"protocol\":\"TCP\"}],\"releaseArgs\":[\"/home/devbox/project/entrypoint.sh\"],\"releaseCommand\":[\"/bin/bash\",\"-c\"],\"user\":\"devbox\",\"workingDir\":\"/home/devbox/project\"}"

	var config map[string]interface{}
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	log.Println(config)
	fmt.Println("Parsed JSON")
}
