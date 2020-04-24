package client

import (
	"time"

	"google.golang.org/grpc"

	pb "github.com/zdnscloud/node-agent/proto"
)

type NodeAgentClient struct {
	pb.NodeAgentClient
	conn *grpc.ClientConn
}

func NewClient(addr string, timeout time.Duration) (*NodeAgentClient, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithTimeout(timeout),
	}

	_conn, err := grpc.Dial(addr, dialOptions...)
	if err != nil {
		return nil, err
	}

	return &NodeAgentClient{
		NodeAgentClient: pb.NewNodeAgentClient(_conn),
		conn:            _conn,
	}, nil
}

func (c *NodeAgentClient) Close() error {
	return c.conn.Close()
}
