package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/narukaz/EnvironmentConfigManager/pkg/connection"
)

type Database struct {
	Host     string `json:"db_host"`
	Port     string `json:"db_port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Server struct {
	Host string `json:"srv_host"`
	Port string `json:"srv_port"`
}
type JsonObject struct {
	Env        string   `json:"env"`
	Database   Database `json:"database"`
	Server     Server   `json:"server"`
	Debug_mode bool     `json:"debug_mode"`
	Log_level  string   `json:"log_level"`
}

func main() {
	var file []byte
	var err error
	var data JsonObject
	value := os.Getenv("Environment")
	fmt.Println("development mode is :", value)
	if value == "" {
		fmt.Println("Please specify Environment")
		return
	} else if value == "development" {
		file, err = os.ReadFile("config/env.development.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(file, &data)
		if err != nil {
			fmt.Println("error unmarshallign the data")
			fmt.Println(err)
		}

	} else if value == "staging" {
		file, err = os.ReadFile("config/env.staging.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(file, &data)
		if err != nil {
			fmt.Println("error unmarshallign the data")
			fmt.Println(err)
		}

	}
	// fmt.Printf("%#v", data)
	client, err := connection.ConnectToMongo(data.Database.Host + ":" + data.Database.Port)

	connection.ServerConnect(data.Server.Host, func(value string) (data int) {
		data, _ = strconv.Atoi(value)
		return data
	}(data.Server.Port), value, client)

	if err != nil {
		fmt.Println("port conversion to string failed")
	}

}
