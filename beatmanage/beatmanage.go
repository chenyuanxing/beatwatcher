package beatmanage

import (
	"../conf"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Name of this beat
type collectionStatus struct {
	Agentuuid string    `json:"agentuuid"`
	Configname string      `json:"configname"`
	Pid int           `json:"pid"`
	Status string    `json:"status"`
	Other string     `json:"other"`
}

// collection 状态  cmds为启动的收集器的exec.Cmd 的集合 ，在添加时，必须同时操作
// collection 状态一旦改为 off ，将不能还原为 on ，同时，会将cmds中对应的cmd 删除
var CollectionStatusSlice [] collectionStatus
var cmds [] *exec.Cmd


func init() {
	fmt.Println("get in here --------------------------")
	//CollectionStatusSlice = make([]collectionStatus,0)

	//curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-6.3.2-x86_64.rpm
	//sudo rpm -vi filebeat-6.3.2-x86_64.rpm


	//curl -L -O https://artifacts.elastic.co/downloads/beats/metricbeat/metricbeat-6.3.2-x86_64.rpm
	//sudo rpm -vi metricbeat-6.3.2-x86_64.rpm
/**
	var metricbeat = "metricbeat"
	var metricbeatT = "metricbeat-6.3.2"
	var metricbeatV = "metricbeat-6.3.2-x86_64"

	var filebeat = "filebeat"
	var filebeatT = "filebeat-6.3.2"
	var filebeatV = "filebeat-6.3.2-x86_64"
	//判断是否需要下载安装，若是
	InstalledList := rpmInstalledList()

	if strings.Contains(strings.ToLower(InstalledList),strings.ToLower(metricbeatT))==false {
		preDownloadandInstall(metricbeat,metricbeatV);
	}
	if strings.Contains(strings.ToLower(InstalledList),strings.ToLower(filebeatT))==false {
		preDownloadandInstall(filebeat,filebeatV);
	}
*/

}

func rpmInstalledList() string{
	cmd := exec.Command("rpm", "-aq")
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error: The command can not be startup: %s\n", err)
		return ""
	}
	if err := cmd.Wait(); err != nil {     //为了获取cmd的所有输出内容，调用Wait()方法一直阻塞到其所属所有命令执行完
		fmt.Printf("Error: Can not wait for the command: %s\n", err)
		return ""
	}
	return out.String();
}

func getYcmdReader() (*bufio.Reader,error)   {
	yescmd := exec.Command("echo", "y","y","y","y","y","y","y","y","y","y")
	stdout1, err := yescmd.StdoutPipe()       //yescmd上建立一个输出管道，为*io.Reader类型
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdout pipe for command: %s", err)
		return nil,err
	}
	if err := yescmd.Start(); err != nil {
		fmt.Printf("Error: The command can not running: %s\n", err)
		return nil,err
	}
	outputBuf1 := bufio.NewReader(stdout1)  //避免数据过多带来的困扰，使用带缓冲的读取器来获取输出管道中的数据
	return outputBuf1,nil
}

func preDownloadandInstall(beat,beatV string){

	beatPKG := beatV+".rpm"
	url := "https://artifacts.elastic.co/downloads/beats/"+beat+"/"+beatPKG

	downloadcmd := exec.Command("curl", "-L","-O",url)
	rpmcmd := exec.Command("sudo","rpm", "-vi",beatPKG)

	fmt.Println("down:    "+"curl", "-L","-O",url)
	fmt.Println("rpmcmd:    "+"sudo","rpm", "-vi",beatPKG)

	outputBuf1,err := getYcmdReader()
	if err != nil {
		fmt.Printf("Error: getYcmdReader Error: %s", err)
		return
	}
	executeCmd(downloadcmd,outputBuf1)
	executeCmd(rpmcmd,outputBuf1)

}


func executeCmd(cmd  *exec.Cmd,outputBuf1 *bufio.Reader) {

	stdin2, err := cmd.StdinPipe()         //cmd上建立一个输入管道
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdin pipe for command: %s\n", err)
		return ;
	}
	outputBuf1.WriteTo(stdin2)              //将缓冲读取器里的输出管道数据写入输入管道里

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error: The command can not be startup: %s\n", err)
		return
	}

	err = stdin2.Close()                    //关闭cmd的输入管道


	if err := cmd.Wait(); err != nil {     //为了获取cmd的所有输出内容，调用Wait()方法一直阻塞到其所属所有命令执行完
		fmt.Printf("Error: Can not wait for the command: %s\n", err)
		return
	}
	fmt.Printf("%s\n", out.Bytes())  //输出执行结果

}
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

func DoServerStuff(conn net.Conn) {
	remote := conn.RemoteAddr().String()
	fmt.Println(remote, " connected!")
	//for {
		// 1024 是数组的长度并且也是切片的初始长度，可增加。
		buf := make([]byte, 5120)

		//一定要等到有error或EOF的时候才会返回结果，因此只能等到客户端退出时才会返回结果。 因此不用此方法
		//buf,err :=ioutil.ReadAll(conn)

		size, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Read Error:", err.Error());
			return
		}
		fmt.Println("data from client:",string(buf),"size:",size)
		var operate Operate
		err = json.Unmarshal(buf[:size], &operate)
		if err != nil {
			fmt.Println("Unmarshal Error:", err.Error());
			return
		}
		fmt.Println("Operate after Unmarshal:", operate)
		var operateReturn Operate;
		operateReturn.Timestamp =time.Now().Unix()

		fmt.Println(operate.Operate)
		if(operate.Operate=="start"){
			fmt.Println("get in")
			cmd := exec.Command("ls", "-l")
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Error %v executing command!", err)
				os.Exit(1)
			}
			fmt.Printf("The command is %v", cmd)
		}else if(operate.Operate=="stop"){
			fmt.Println("get the operate stop")
			operateReturn.Operate = "success"
			operateReturn.Timestamp =time.Now().Unix()
			buf, err = json.Marshal(operateReturn)
			if err != nil {
				fmt.Println("Marshal Error:", err.Error());
				return
			}
			conn.Write(buf)
			conn.Close()
			os.Exit(1);
		}else if(operate.Operate=="metricbeat"){
			configName,err:=operate.File.Get("name").String();
			if err != nil {
				fmt.Println("MarshalJSON Error:", err.Error());
				return
			}
				metricbeatYml:=operate.Operate+"_"+configName+".yml"
				metricbeatModulesYml:=operate.Operate+"_"+configName+"Modules.yml"

				operate.File.Get("jsonFile").Get("metricbeat.config.modules").Set("path","${path.config}/modules.d"+"/"+metricbeatModulesYml)
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


				ioutil.WriteFile(conf.Config.MetricbeatFolder+"/"+metricbeatYml,ymlfile,os.ModeAppend)
				ioutil.WriteFile(conf.Config.MetricbeatFolder+"/modules.d"+"/"+metricbeatModulesYml,ModulesYmlFile,os.ModeAppend)
				if err != nil {
					fmt.Println("WriteFile Error:", err.Error());
					return
				}

				operateReturn.Operate = "success"

				// 此处启动 待续..
				//  ./metricbeat-6.5.4-linux-x86_64/metricbeat -e -c ./metricbeat-6.5.4-linux-x86_64/metricbeat_new_Collection.yml
				launchcmd := exec.Command("./"+conf.Config.MetricbeatFolder+"/"+conf.Config.Metricbeat, "-c","./"+conf.Config.MetricbeatFolder+"/"+metricbeatYml)

				fmt.Println("launchcmd ")
				fmt.Println( "./"+conf.Config.MetricbeatFolder+"/"+conf.Config.Metricbeat, "-c","./"+conf.Config.MetricbeatFolder+"/"+metricbeatYml)

				err =launchcmd.Start()
				if err != nil {
					fmt.Println("launchcmd start Error:", err.Error());
					return
				}

				var cStatus collectionStatus;
				cStatus.Pid = launchcmd.Process.Pid
				cStatus.Status = "on"
				cStatus.Configname = configName
				cStatus.Agentuuid = conf.Uuid
				CollectionStatusSlice = append(CollectionStatusSlice,cStatus)

				cmds = append(cmds, launchcmd)

		}else if(operate.Operate=="filebeat"){
			//启动



		}else if(operate.Operate=="metricbeat_stop"){
			//停止
			operateReturn.Operate = killbeat(operate.Param)
		}else if(operate.Operate=="filebeat_stop"){
			//停止
			operateReturn.Operate = killbeat(operate.Param)
		}

		operateReturn.Timestamp =time.Now().Unix()
		buf, err = json.Marshal(operateReturn)
		if err != nil {
			fmt.Println("Marshal Error:", err.Error());
			return
		}
		conn.Write(buf)

		conn.Close()
		//break
	//}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) (bool) {
	var exist = true;
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false;
	}
	return exist;
}

func updateCollectionStatus(status *collectionStatus){
	cmd := exec.Command("cat","/proc/"+strconv.Itoa(status.Pid)+"/cmdline")
	var out bytes.Buffer
	cmd.Stdout = &out
	err:=cmd.Run()
	if err != nil {
		fmt.Printf("Error: execute cmd "+"cat","/proc/"+strconv.Itoa(status.Pid)+"/cmdline "+":\n %s", err)
		status.Status = "off"
		return
	}
	if strings.Contains(out.String(),status.Configname){
		status.Status = "on"
	}else {
		// 说明该pid虽然存在，但是不是之前启动的pid
		status.Status = "off"
	}
}

// 杀掉某个状态为on的 beat 进程
func killbeat(pid int) string{
	stopPid := pid
	index :=-1
	for i :=range cmds{
		if cmds[i].Process.Pid == stopPid{
			index = i
			break
		}
	}
	if index >=0 {
		//杀死该进程
		err:=cmds[index].Process.Kill()
		if err != nil {
			return "kill_failed"
		}
		//删除cmds中被删除的元素
		cmds =  append(cmds[:index], cmds[index+1:]...)      // 最后面的“...”不能省略
	}else {
		return "notfind"
	}
	return "success"
}