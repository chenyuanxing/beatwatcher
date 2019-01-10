package main

import (
	"./beatmanage"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type addrobject struct {

}

type Operate struct {
	Operate string  `json:"operate"`
	Param   int     `json:"param"`
	Id      string  `json:"id"`
	Timestamp int64 `json:"timestamp"`
	File string     `json:"file"`
	//other  []string
}

func doServerStuff(conn net.Conn) {
	remote := conn.RemoteAddr().String()
	fmt.Println(remote, " connected!")
	for {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read Error:", err.Error());
			return
		}
		//fmt.Println("data from client:",string(buf),"size:",size)
		var operate Operate
		s := string(buf[:size])
		fmt.Println("get string:", s)
		err = json.Unmarshal(buf[:size], &operate)
		if err != nil {
			fmt.Println("Unmarshal Error:", err.Error());
			return
		}
		fmt.Println("Operate after Unmarshal:", operate)
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
			var operateR Operate;
			operateR.Operate = "success"
			operateR.Timestamp =time.Now().Unix()
			buf, err = json.Marshal(operateR)
			if err != nil {
				fmt.Println("Marshal Error:", err.Error());
				return
			}
			conn.Write(buf)
			conn.Close()
			os.Exit(1);
		}

		operate.Timestamp += 2
		buf, err = json.Marshal(operate)
		if err != nil {
			fmt.Println("Marshal Error:", err.Error());
			return
		}
		conn.Write(buf)
		conn.Close()
		break
	}
}


func main() {
	ok := flag.Bool("ok", false, "is ok")

	id := flag.Int("id", 0, "id")
	port := flag.String("port", "50001", "http listen port")
	key := flag.String("k","test_key123","key of the beatwather")
	user := flag.String("u","chen","user of the beatwather")
	tag := flag.String("tag","test_tag","tag")

	flag.Parse()

	fmt.Println("ok:", *ok)
	fmt.Println("id:", *id)
	fmt.Println("port:", *port)
	fmt.Println("key:",*key)
	fmt.Println("user:",*user)
	fmt.Println("tag:",*tag)


	fmt.Println("Starting the server...")
	listener, err := net.Listen("tcp", "0.0.0.0:"+string(*port))

	version:="0.0.1"

	realport := *port;
	if err != nil {
		fmt.Println("Listen Error:", err.Error())

		for i:=2; i<10;i++{
			realport = strconv.Itoa(50000+i)
			fmt.Println("try another port: ",realport)
			listener, err = net.Listen("tcp", "0.0.0.0:"+realport)
			if err != nil {
				fmt.Println("Listen Error:", err.Error())
			}else {
				break
			}
		}
		if err != nil {
			fmt.Println("Listen Error:", err.Error())
			return
		}
	}

	beatmanage.Regist(*key,*user,version,*tag,realport)
	go beatmanage.Heart(*key,*user)

	fmt.Println(" is Starting...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept Error:", err.Error())
			return
		}

		go doServerStuff(conn)
	}
}
