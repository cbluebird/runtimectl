package util

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestJson(t *testing.T) {
	// Replace single quotes with double quotes to make it a valid JSON string
	configStr := "{\n  \"appPorts\": [\n    {\n      \"name\": \"devbox-app-port\",\n      \"port\": 4200,\n      \"protocol\": \"TCP\"\n    }\n  ],\n  \"ports\": [\n    {\n      \"containerPort\": 22,\n      \"name\": \"devbox-ssh-port\",\n      \"protocol\": \"TCP\"\n    }\n  ],\n  \"releaseArgs\": [\n    \"/home/devbox/project/entrypoint.sh\"\n  ],\n  \"releaseCommand\": [\n    \"/bin/bash\",\n    \"-c\"\n  ],\n  \"user\": \"devbox\",\n  \"workingDir\": \"/home/devbox/project\"\n}"

	var config map[string]interface{}
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	log.Println(config)
	fmt.Println("Parsed JSON")
}
