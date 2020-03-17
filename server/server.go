package server

import (
	"fmt"

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

func (s Server) IscsiDiscovery(ctx context.Context, in *pb.IscsiDiscoveryRequest) (*pb.IscsiDiscoveryReply, error) {
	addr := in.Host + ":" + in.Port
	if err := command.DiscoverTarget(addr, in.Iqn); err != nil {
		return nil, err
	}
	ok, err := command.IsTargetDiscovered(addr, in.Iqn)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("iscsi check discovery faild")
	}
	return &pb.IscsiDiscoveryReply{}, nil
}

func (s Server) IscsiChap(ctx context.Context, in *pb.IscsiChapRequest) (*pb.IscsiChapReply, error) {
	if err := command.IscsiChap(in.Host+":"+in.Port, in.Iqn, in.Username, in.Password); err != nil {
		return nil, err
	}
	return &pb.IscsiChapReply{}, nil
}
func (s Server) IscsiLogin(ctx context.Context, in *pb.IscsiLoginRequest) (*pb.IscsiLoginReply, error) {
	if err := command.LoginTarget(in.Host+":"+in.Port, in.Iqn); err != nil {
		return nil, err
	}
	return &pb.IscsiLoginReply{}, nil
}
func (s Server) IscsiLogout(ctx context.Context, in *pb.IscsiLogoutRequest) (*pb.IscsiLogoutReply, error) {
	if err := command.LogoutTarget(in.Host+":"+in.Port, in.Iqn); err != nil {
		return nil, err
	}
	return &pb.IscsiLogoutReply{}, nil
}
func (s Server) IscsiGetBlocks(ctx context.Context, in *pb.IscsiGetBlocksRequest) (*pb.IscsiGetBlocksReply, error) {
	output, err := command.GetIscsiDevices(in.Host, in.Iqn)
	if err != nil {
		return nil, err
	}
	info := make(map[string]*pb.IscsiDevice)
	for lun, devs := range output {
		info[lun] = &pb.IscsiDevice{
			Blocks: devs,
		}
	}
	return &pb.IscsiGetBlocksReply{
		IscsiBlock: info,
	}, nil
}

func (s Server) IscsiGetMultipaths(ctx context.Context, in *pb.IscsiGetMultipathsRequest) (*pb.IscsiGetMultipathsReply, error) {
	path, err := command.GetIscsiMultipath(in.Devs)
	if err != nil {
		return nil, err
	}
	return &pb.IscsiGetMultipathsReply{
		Dev: path,
	}, nil
}

func (s Server) ReplaceInitiatorname(ctx context.Context, in *pb.ReplaceInitiatornameRequest) (*pb.ReplaceInitiatornameReply, error) {
	if err := command.ReplaceInitiatorname(in.SrcFile, in.DstFile); err != nil {
		return nil, err
	}
	return &pb.ReplaceInitiatornameReply{}, nil
}
