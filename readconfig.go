package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func ReadConfigfile(ConfigName string) map[string]interface{} {
	var configJson map[string]interface{}
	configfile, err := os.Open(ConfigName)

	if err != nil {
		log.Fatal(err)
	}
	jsonResult, err := io.ReadAll(configfile)
	defer configfile.Close()
	err = json.Unmarshal([]byte(jsonResult), &configJson)
	if err != nil {
		log.Fatal(err)
	}
	return configJson
}
