/*

Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/zdnscloud/cement/log"

	pb "github.com/zdnscloud/node-agent/proto"
	"github.com/zdnscloud/node-agent/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

func main() {
	var addr string
	var clusterAgentAddr string
	var nodeName string

	flag.StringVar(&addr, "listen", ":80", "server listen address")
	flag.StringVar(&clusterAgentAddr, "cluster", "cluster-agent.zcloud.svc", "server listen address")
	flag.StringVar(&nodeName, "node", "", "node name this node agent resides")
	flag.Parse()

	log.InitLogger(log.Debug)

	if nodeName == "" {
		log.Fatalf("node name isn't specified")
	} else if clusterAgentAddr == "" {
		log.Fatalf("cluster agent addr isn't specified")
	}
	for {
		err := registerWithClusterAgent(clusterAgentAddr, nodeName, addr)
		if err != nil {
			log.Errorf("register to cluster agent failed:%s, start to retry", err.Error())
			<-time.After(10 * time.Second)
		} else {
			log.Infof("register to cluster agent %s succeed", clusterAgentAddr)
			break
		}
	}

	svr := server.NewServer()
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterNodeAgentServer(grpcServer, &svr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen failed:%s", err.Error())
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("run grpc server failed:%s", err.Error())
	}
}
