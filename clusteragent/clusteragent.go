package clusteragent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/zdnscloud/cement/log"
)

const HeartbeatInterval = 10 //10 second

func registerWithClusterAgent(clusterAgentAddr, nodeName, nodeAddr string) error {
	client := &http.Client{}
	url := fmt.Sprintf("http://%s/apis/agent.zcloud.cn/v1/nodeagents", clusterAgentAddr)
	requestBody, _ := json.Marshal(map[string]interface{}{
		"name":    nodeName,
		"address": nodeAddr,
	})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		return errors.New(string(body))
	}

	return nil
}

func StartHeartbeat(clusterAgentAddr, nodeName, addr string) {
	go func() {
		for {
			err := registerWithClusterAgent(clusterAgentAddr, nodeName, addr)
			if err != nil {
				log.Errorf("register to cluster agent failed:%s, start to retry", err.Error())
			}
			<-time.After(time.Duration(HeartbeatInterval+rand.Intn(5)) * time.Second)
		}
	}()
}
