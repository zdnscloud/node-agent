package clusteragent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/zdnscloud/cement/log"
	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/gok8s/client/config"
)

var (
	flag bool
	cli  client.Client
)

const (
	HeartbeatInterval   = 10 //10 second
	ZCloudNamespace     = "zcloud"
	ClusterAgentSvcName = "cluster-agent"
	ClusterAgentSvcPort = "80"
)

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

func StartHeartbeat(nodeName, addr string) {
	go func() {
		for {
			clusterAgentAddr, err := getClusterAgentSvcAddr()
			if err != nil {
				log.Errorf("get cluster agent service addr failed:%s, start to retry", err.Error())
			}
			err = registerWithClusterAgent(clusterAgentAddr, nodeName, addr)
			if err != nil {
				log.Errorf("register to cluster agent failed:%s, start to retry", err.Error())
			}
			<-time.After(time.Duration(HeartbeatInterval+rand.Intn(5)) * time.Second)
		}
	}()
}

func getClusterAgentSvcAddr() (string, error) {
	if !flag {
		cfg, err := config.GetConfig()
		if err != nil {
			return "", err
		}
		cli, err = client.New(cfg, client.Options{})
		if err != nil {
			return "", err
		}
		flag = true
	}
	service := corev1.Service{}
	err := cli.Get(context.TODO(), k8stypes.NamespacedName{ZCloudNamespace, ClusterAgentSvcName}, &service)
	if err != nil {
		return "", err
	}
	return service.Spec.ClusterIP + ":" + ClusterAgentSvcPort, nil
}
