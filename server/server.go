package server

import (
	"path/filepath"

	"bytes"
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
	cmds := []*exec.Cmd{
		exec.Command("df"),
		exec.Command("grep", path),
		exec.Command("awk", "{print $3}"),
	}
	for i := 1; i < len(cmds); i++ {
		var err error
		if cmds[i].Stdin, err = cmds[i-1].StdoutPipe(); err != nil {
			return size, err
		}

	}
	var output bytes.Buffer
	cmds[len(cmds)-1].Stdout = &output
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return size, err
		}
	}
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return size, err
		}
	}
	size, err := strconv.ParseInt(strings.Replace(string(output.Bytes()), "\n", "", -1), 10, 64)
	return size, err
}
