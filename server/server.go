package server

import (
	"path/filepath"

	"golang.org/x/net/context"
	"os"
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
	infos, err := getBlockUsedSize(in.Paths, in.Type)
	return &pb.GetBlockUsedSizeReply{
		Infos: infos,
	}, err
}

func getBlockUsedSize(paths []string, t string) (map[string]int64, error) {
	infos := make(map[string]int64)
	out, _ := exec.Command("df", paths...).Output()
	outputs := strings.Split(string(out), "\n")
	n := 4
	switch t {
	case "t":
		n = 5
	case "u":
		n = 4
	case "f":
		n = 3
	}
	for i := 1; i < len(outputs); i++ {
		if !strings.Contains(outputs[i], "%") {
			continue
		}
		line := strings.Fields(outputs[i])
		num := len(line)
		size, err := strconv.ParseInt(line[num-n], 10, 64)
		if err != nil {
			return infos, err
		}
		infos[line[num-1]] = size
	}
	return infos, nil
}
