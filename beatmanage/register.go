package beatmanage

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
)

func Regist(key,user,version,tag,port string)  {

	params := url.Values{}
	params.Set("uuidKey",key)
	params.Set("agentVersion",version)
	params.Set("tag",tag)

	params.Set("system",runtime.GOOS)
	params.Set("kernelVersion",runtime.GOARCH)
	params.Set("port",port)

	resp, _:= http.PostForm("http://10.108.210.194:8080/agents/registAgent",params)
	buf := make([]byte, 512)
	resp.Body.Read(buf)
	fmt.Println("regist data : \n" , params)
	fmt.Println("register result : \n", string(buf[:]))

}
