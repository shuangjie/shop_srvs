package main

import (
	"flag"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"srvs/user_srv/handler"
	"srvs/user_srv/initialize"
	"srvs/user_srv/proto"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip address")
	Port := flag.Int("port", 50051, "port")

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()
	zap.S().Info("ip: %s, port: %d\n", *IP, *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}
	err = server.Serve(lis)
	if err != nil {
		panic("failed to serve: " + err.Error())
	}

}
