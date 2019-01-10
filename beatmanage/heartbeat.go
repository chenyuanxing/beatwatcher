package beatmanage

import (
	"net/http"
	"net/url"
	"time"
)

func Heart(uuid,user string)  {
	for range time.Tick(5 * time.Second) {
		params := url.Values{}
		params.Set("uuidKey",uuid)

		params.Set("userId",user)
		http.PostForm("http://10.108.210.194:8080/agents/aliveAgent",params)
		//buf := make([]byte, 512)
		//resp.Body.Read(buf)
		//fmt.Println("resp.Body : ", string(buf[:]))




	}
}