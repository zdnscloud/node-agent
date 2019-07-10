package server

import (
	"golang.org/x/net/context"

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

func (s Server) GetMountpointsSize(ctx context.Context, in *pb.GetMountpointsSizeRequest) (*pb.GetMountpointsSizeReply, error) {
	infos, err := getMountpointsSize(in.Paths)
	return &pb.GetMountpointsSizeReply{
		Infos: infos,
	}, err
}

func (s Server) GetDisksInfo(ctx context.Context, in *pb.GetDisksInfoRequest) (*pb.GetDisksInfoReply, error) {
	infos, err := getDisksInfo(in.Disk)
	return &pb.GetDisksInfoReply{
		Infos: infos,
	}, err
}
