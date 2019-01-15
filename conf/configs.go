package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

var Config  = Configuration{}

type Configuration struct {
	Users               []string
	Groups              []string
	Metricbeat          string
	MetricbeatFolder    string
	Filebeat            string
	FilebeatFolder      string
}


func init() {
	fmt.Println("init config here-----------------")
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	Config = configuration
	fmt.Println("config is : ")
	fmt.Println(Config)

}