package server

import (
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

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

func (s Server) GetBlockUsedSizeSize(ctx context.Context, in *pb.GetBlockUsedSizeRequest) (*pb.GetBlockUsedSizeReply, error) {
	infos := make([]*pb.GetBlockUsedSizeReplyInfo, 0)
	for _, p := range in.Paths {
		size, err := getBlockUsedSize(p)
		if err != nil {
			return nil, err
		} else {
			info := &pb.GetBlockUsedSizeReplyInfo{
				Path: p,
				Size: size,
			}
			infos = append(infos, info)
		}
	}
	return &pb.GetBlockUsedSizeReply{
		Infos: infos,
	}, nil
}

func getBlockUsedSize(path string) (int64, error) {
	var size int64
	command := "df " + path + "|tail -n1|awk '{print $3}'"
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return size, err
	}
	if err := cmd.Start(); err != nil {
		return size, err
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return size, err
	}
	if err := cmd.Wait(); err != nil {
		return size, err
	}
	s := strings.Replace(string(bytes[:]), "\n", "", -1)
	size, err = strconv.ParseInt(s, 10, 64)
	return size, err
}
