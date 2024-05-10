package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"srvs/user_srv/global"
	"srvs/user_srv/handler"
	"srvs/user_srv/initialize"
	"srvs/user_srv/proto"
	"srvs/user_srv/utils"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip address")
	Port := flag.Int("port", 0, "port")

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	zap.S().Info(global.ServerConfig)

	flag.Parse()
	zap.S().Infof("ip:%s,port:%d", *IP, *Port)

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("port:", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	//注册健康检查服务
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("192.168.2.49:%d", *Port),
		Interval:                       "5s",
		Timeout:                        "3s",
		DeregisterCriticalServiceAfter: "10s",
	}

	registration := new(api.AgentServiceRegistration)
	serverID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serverID
	registration.Name = global.ServerConfig.Name
	registration.Port = *Port
	registration.Tags = []string{"zeng", "user", "grpc"}
	registration.Address = "192.168.2.49"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to serve: " + err.Error())
		}
	}()

	//优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serverID); err != nil {
		zap.S().Info("注销失败")
	} else {
		zap.S().Info("注销成功")
	}

}
