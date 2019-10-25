package server

import (
	"golang.org/x/net/context"

	"github.com/zdnscloud/node-agent/command"
	pb "github.com/zdnscloud/node-agent/proto"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (s Server) GetDirectorySize(ctx context.Context, in *pb.GetDirectorySizeRequest) (*pb.GetDirectorySizeReply, error) {
	infos, err := command.GetDirectorySize(in.Path)
	return &pb.GetDirectorySizeReply{
		Infos: infos,
	}, err
}

func (s Server) GetMountpointsSize(ctx context.Context, in *pb.GetMountpointsSizeRequest) (*pb.GetMountpointsSizeReply, error) {
	infos, err := command.GetMountpointsSize(in.Paths)
	return &pb.GetMountpointsSizeReply{
		Infos: infos,
	}, err
}

func (s Server) GetDisksInfo(ctx context.Context, in *pb.GetDisksInfoRequest) (*pb.GetDisksInfoReply, error) {
	infos, err := command.GetDisksInfo(in.Disk)
	return &pb.GetDisksInfoReply{
		Infos: infos,
	}, err
}
