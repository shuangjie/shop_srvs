package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"srvs/goods_srv/global"
	"srvs/goods_srv/handler"
	"srvs/goods_srv/initialize"
	"srvs/goods_srv/proto"
	"srvs/goods_srv/utils"
	"srvs/goods_srv/utils/register/consul"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip address")
	Port := flag.Int("port", 50051, "port")

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
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	//注册健康检查服务
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := uuid.NewV4().String()
	err = registerClient.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("goods_srv注册服务失败", err.Error())
	}
	zap.S().Debugf("goods_srv服务启动中..., 端口: %d", *Port)

	//启动服务
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to serve: " + err.Error())
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Info("goods_srv注销服务失败", err.Error())
	} else {
		zap.S().Info("goods_srv服务注销成功")
	}

}
