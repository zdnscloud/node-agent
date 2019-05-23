package server

import (
	"golang.org/x/net/context"
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
	infos, err := getDirectorySize(in.Path)
	return &pb.GetDirectorySizeReply{
		Infos: infos,
	}, err
}

func getDirectorySize(path string) (map[string]int64, error) {
	infos := make(map[string]int64)
	out, _ := exec.Command("du", "-d1", path).Output()
	outputs := strings.Split(string(out), "\n")
	for _, l := range outputs {
		if !strings.Contains(l, "/") {
			continue
		}
		line := strings.Fields(l)
		size, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			return infos, err
		}
		infos[line[1]] = size
	}
	delete(infos, path)
	return infos, nil
}

func (s Server) GetMountpointsSize(ctx context.Context, in *pb.GetMountpointsSizeRequest) (*pb.GetMountpointsSizeReply, error) {
	infos, err := getMountpointsSize(in.Paths)
	return &pb.GetMountpointsSizeReply{
		Infos: infos,
	}, err
}

func getMountpointsSize(paths []string) (map[string]*pb.Sizes, error) {
	infos := make(map[string]*pb.Sizes)
	out, _ := exec.Command("df", paths...).Output()
	outputs := strings.Split(string(out), "\n")
	for i := 1; i < len(outputs); i++ {
		if !strings.Contains(outputs[i], "%") {
			continue
		}
		line := strings.Fields(outputs[i])
		num := len(line)
		tsize, err := strconv.ParseInt(line[num-5], 10, 64)
		if err != nil {
			return infos, err
		}
		usize, err := strconv.ParseInt(line[num-4], 10, 64)
		if err != nil {
			return infos, err
		}
		fsize, err := strconv.ParseInt(line[num-3], 10, 64)
		if err != nil {
			return infos, err
		}
		infos[line[num-1]] = &pb.Sizes{
			Tsize: tsize,
			Usize: usize,
			Fsize: fsize,
		}
	}
	return infos, nil
}
