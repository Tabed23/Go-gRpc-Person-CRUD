package cmd

import (
	"fmt"
	"log"
	"net"
	"persons-gRpc/pb"
	"persons-gRpc/service"

	"google.golang.org/grpc"
)
const (
	port = ":50001"
)
func StartGrpcServer(){
	opts := []grpc.ServerOption{}
	rpcServer := grpc.NewServer(opts...)
	pb.RegisterPersonServiceServer(rpcServer,&service.PersonServer{})
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("error with server %v", err)
	}
	fmt.Printf("gRpc server started at port [%v]", port)
	if err = rpcServer.Serve(lis); err != nil {
		log.Fatalf("cannot server the request %v", err)
	}
}