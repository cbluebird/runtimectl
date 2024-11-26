package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtimectl/model"
)

func ParseJson() *model.Config {
	// Open the JSON file
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return nil
	}
	defer jsonFile.Close()

	// Read the JSON file
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return nil
	}

	// Unmarshal the JSON data into the struct
	var config model.Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println("Error unmarshalling JSON data:", err)
		return nil
	}
	return &config
}
