package beatmanage

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Name of this beat


func init() {
	fmt.Println("get in here --------------------------")
	//curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-6.3.2-x86_64.rpm
	//sudo rpm -vi filebeat-6.3.2-x86_64.rpm


	//curl -L -O https://artifacts.elastic.co/downloads/beats/metricbeat/metricbeat-6.3.2-x86_64.rpm
	//sudo rpm -vi metricbeat-6.3.2-x86_64.rpm
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


func preDownloadandInstall(beat,beatV string){
	yescmd := exec.Command("echo", "y","y","y","y","y","y","y","y","y","y")
	stdout1, err := yescmd.StdoutPipe()       //yescmd上建立一个输出管道，为*io.Reader类型
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdout pipe for command: %s", err)
		return
	}
	if err := yescmd.Start(); err != nil {
		fmt.Printf("Error: The command can not running: %s\n", err)
		return
	}
	outputBuf1 := bufio.NewReader(stdout1)  //避免数据过多带来的困扰，使用带缓冲的读取器来获取输出管道中的数据
	beatPKG := beatV+".rpm"
	url := "https://artifacts.elastic.co/downloads/beats/"+beat+"/"+beatPKG

	downloadcmd := exec.Command("curl", "-L","-O",url)
	rpmcmd := exec.Command("sudo","rpm", "-vi",beatPKG)

	fmt.Println("down:    "+"curl", "-L","-O",url)
	fmt.Println("rpmcmd:    "+"sudo","rpm", "-vi",beatPKG)

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

