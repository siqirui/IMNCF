package Grpc

import (
	"IMCF/Ilog"
	pb "IMCF/Inet/Grpc/proto"
	utils "IMCF/Utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

func NewClent() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "127.0.0.1", utils.GlobalOBJ.GrpcServerHost), grpc.WithInsecure())
	if err != nil {

	}
	defer conn.Close()

	clinet := pb.NewUserInfoServiceClient(conn)

	req := new(pb.UerRequest)

	req.Name = "zs"

	resp, err := clinet.GetUserInfo(context.Background(), req)

	if err != nil {

	}
	Ilog.Info("recv:", resp)
}
