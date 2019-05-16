package server

import (
	"os"
	"path/filepath"

	"golang.org/x/net/context"

	pb "github.com/zdnscloud/node-agent/proto"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (s Server) GetDirectorySize(ctx context.Context, in *pb.GetDirectorySizeRequest) (*pb.GetDirectorySizeReply, error) {
	size, err := getDirectorySize(in.Path)
	if err != nil {
		return nil, err
	} else {
		return &pb.GetDirectorySizeReply{
			Size: size,
		}, nil
	}
}

func getDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
