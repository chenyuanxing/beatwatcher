package beatmanage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func Heart(uuid,user string)  {
	for range time.Tick(5 * time.Second) {
		params := url.Values{}
		params.Set("uuidKey",uuid)

		params.Set("userId",user)

		// 发送心跳前将collection状态更新,只检查状态为on的
		for i :=range CollectionStatusSlice{
			if CollectionStatusSlice[i].Status=="on" {

				//updateCollectionStatus(&CollectionStatusSlice[i])
			}
		}

		buf,err := json.Marshal(CollectionStatusSlice)
		if err != nil {
			fmt.Println("Heart - Marshal Error: ", err.Error());
			return
		}
		params.Set("collectionStatuses",string(buf))
		//fmt.Println("heart message:")
		//fmt.Println(params)
		http.PostForm("http://10.108.210.194:8080/agents/aliveAgent",params)

	}
}