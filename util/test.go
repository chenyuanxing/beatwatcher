package main

import (
	"encoding/json"
	"fmt"
	"../conf"
	"github.com/bitly/go-simplejson"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"../beatmanage"
)



type Operate struct {
	Operate   string  `json:"operate"`
	Param     int     `json:"param"`
	Id        string  `json:"id"`
	Timestamp int64 `json:"timestamp"`
	File      simplejson.Json     `json:"file"`
	//other  []string
}
type BeatJson struct {
	Name            string  `json:"name"`
	JsonFile         string     `json:"jsonFile"`
	ModulesJsonFile  string  `json:"modulesJsonFile"`
}

func test1()  {

	// 512 是数组的长度并且也是切片的初始长度，可增加。
	buf := []byte("{\n" +
		"  \"file\": {\n" +
		"    \"jsonFile\": {\n" +
		"      \"setup.template.settings\": {\n" +
		"        \"index.number_of_shards\": 1,\n" +
		"        \"index.codec\": \"best_compression\"\n" +
		"      },\n" +
		"      \"output.elasticsearch\": {\n" +
		"        \"hosts\": [\n" +
		"          \"localhost:9200\",\n" +
		"          \"10.108.210.194:8080\"\n" +
		"        ]\n" +
		"      },\n" +
		"      \"metricbeat.config.modules\": {\n" +
		"        \"reload.enabled\": false,\n" +
		"        \"path\": \"${path.config}/modules.d/smj1547363524203.yml\"\n" +
		"      }\n" +
		"    },\n" +
		"    \"modulesJsonFile\": {\n" +
		"      \"metricbeat.modules\": [\n" +
		"        {\n" +
		"          \"module\": \"system\",\n" +
		"          \"metricsets\": [\n" +
		"            \"cpu\",\n" +
		"            \"filesystem\",\n" +
		"            \"fsstat\",\n" +
		"            \"load\",\n" +
		"            \"memory\",\n" +
		"            \"network\",\n" +
		"            \"process\",\n" +
		"            \"process_summary\",\n" +
		"            \"uptime\"\n" +
		"          ],\n" +
		"          \"enabled\": true,\n" +
		"          \"period\": \"10s\",\n" +
		"          \"processes\": [\n" +
		"            \".*\"\n" +
		"          ],\n" +
		"          \"cpu.metrics\": [\n" +
		"            \"percentages\"\n" +
		"          ],\n" +
		"          \"core.metrics\": [\n" +
		"            \"percentages\"\n" +
		"          ]\n" +
		"        },\n" +
		"        {\n" +
		"          \"module\": \"mysql\",\n" +
		"          \"period\": \"10s\",\n" +
		"          \"hosts\": [\n" +
		"            \"root:secret@tcp(127.0.0.1:3306)/\"\n" +
		"          ]\n" +
		"        }\n" +
		"      ]\n" +
		"    },\n" +
		"    \"name\": \"new_Collection\"\n" +
		"  },\n" +
		"  \"id\": \"\",\n" +
		"  \"operate\": \"metricbeat\",\n" +
		"  \"param\": 1,\n" +
		"  \"timestamp\": 1547368715987\n" +
		"}")

	var operate Operate
	err := json.Unmarshal(buf[:], &operate)
	if err != nil {
		fmt.Println("Unmarshal Error:", err.Error());
		return
	}
	fmt.Println("Operate after Unmarshal: \n", operate)
	fmt.Println(operate.Operate)
	if(operate.Operate=="metricbeat"){
		configName,err:=operate.File.Get("name").String();
		if err != nil {
			fmt.Println("MarshalJSON Error:", err.Error());
			return
		}
		////启动
		if(operate.Param==1){
			operate.File.Get("jsonFile").Get("metricbeat.config.modules").Set("path","${path.config}/"+operate.Operate+"_"+configName+"Modules.yml")
			jsonFilebuf,err :=operate.File.Get("jsonFile").MarshalJSON()
			if err != nil {
				fmt.Println("MarshalJSON Error:\n", err.Error());
				return
			}
			modulesJsonFilebuf,err :=operate.File.Get("modulesJsonFile").MarshalJSON()
			if err != nil {
				fmt.Println("MarshalJSON Error:\n", err.Error());
				return
			}

			ymlfile, err :=yaml.JSONToYAML(jsonFilebuf)
			if err != nil {
				fmt.Println("JSONToYAML Error:", err.Error());
				return
			}

			ModulesYmlFile, err :=yaml.JSONToYAML(modulesJsonFilebuf)
			if err != nil {
				fmt.Println("JSONToYAML Error:", err.Error());
				return
			}
			// WriteFile 向文件 filename 中写入数据 data
			// 如果文件不存在，则以 perm 权限创建该文件
			// 如果文件存在，则先清空文件，然后再写入

			ioutil.WriteFile(operate.Operate+"_"+configName+".yml",ymlfile,os.ModeAppend)
			ioutil.WriteFile(operate.Operate+"_"+configName+"Modules.yml",ModulesYmlFile,os.ModeAppend)
			if err != nil {
				fmt.Println("WriteFile Error:", err.Error());
				return
			}

		}
	}

}

func testConfig()  {
	fmt.Println("testing")
	fmt.Println(conf.Config.MetricbeatFolder)
}
func testMar(){

	beatmanage.CollectionStatusSlice = append(beatmanage.CollectionStatusSlice, )
	buf,err := json.Marshal(beatmanage.CollectionStatusSlice)
	if err != nil {
		fmt.Println("Heart - Marshal Error: ", err.Error());
		return
	}
	fmt.Println(string(buf))
}
func main() {
	testMar()
}
