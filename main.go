package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/ini.v1"
)

type CommandPayload struct {
	Command string `json:"command"`
	Timeout int    `json:"timeout"`
	Param   string `json:"param"`
	Authkey string `json:"authkey"`
}

func main() {
	//[0] =
	//1   = timeout
	//[2] = base_binary
	//[n] = parameter
	cfg, err := ini.Load("/app/Content/configurations/remote-agent.ini")
	agent_address := cfg.Section("").Key("agent_address").String()
	static_key := cfg.Section("").Key("agent_address").String()
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	httpposturl := agent_address + "/execute"
	timeout, err := strconv.Atoi(os.Args[1])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	binary_param := ""
	base_binary := os.Args[2]
	for i := 3; i < len(os.Args); i++ {
		if i == 3 {
			binary_param = os.Args[i]
		} else {
			binary_param = binary_param + " " + os.Args[i]
		}
		fmt.Println("add " + os.Args[i] + " jadi->" + binary_param)
	}

	fmt.Println("param " + binary_param)
	fmt.Println("base_binary " + base_binary)
	command := CommandPayload{
		Command: base_binary,
		Timeout: timeout,
		Param:   binary_param,
		Authkey: static_key,
	}

	commandJson, err := json.Marshal(command)
	if err != nil {
		panic(err)
	}
	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(commandJson))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	os.Exit(0)
}
