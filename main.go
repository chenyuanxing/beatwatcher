package main

import (
	"./beatmanage"
	"flag"
	"fmt"
	"./conf"
	"net"
	"strconv"
)

type addrobject struct {

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

	conf.Uuid = *key

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

		go beatmanage.DoServerStuff(conn)
	}
}
