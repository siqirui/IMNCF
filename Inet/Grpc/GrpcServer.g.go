package Grpc

import (
	"IMCF/Ilog"
	pb "IMCF/Inet/Grpc/proto"
	utils "IMCF/Utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type UserInfoService struct {
}

var user = UserInfoService{}

func (this *UserInfoService) GetUserInfo(ctx context.Context, req *pb.UerRequest) (resp *pb.UerResponse, err error) {
	name := req.Name

	if name == "zs" {
		resp = &pb.UerResponse{
			Id:    1,
			Name:  name,
			Age:   11,
			Hobby: []string{"sing", "Run"},
		}
	}
	err = nil
	return resp, err
}

func NewGrpcServer() {

}
func (this *UserInfoService) initLister() {
	lister, err := net.Listen("tcp", fmt.Sprintf("%s:%d", utils.GlobalOBJ.GrpcServerHost, utils.GlobalOBJ.GrpcServerPort))
	if err != nil {
		Ilog.Error("GrpcServer :Listern error :", err)
		return
	}
	Ilog.Info("GrpcServer:Listern succ!")

	server := grpc.NewServer()
	pb.RegisterUserInfoServiceServer(server, &user)
	server.Serve(lister)
}
