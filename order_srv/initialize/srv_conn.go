package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"srvs/order_srv/config"
	"srvs/order_srv/global"
	"srvs/order_srv/proto"
)

/*func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	// 初始化商品服务客户端
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 连接【商品服务】失败: %v", err)
	}

	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	// 初始化库存服务客户端
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 连接【库存服务】失败: %v", err)
	}

	global.InventorySrvClient = proto.NewInventoryClient(inventoryConn)
}*/

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo

	// 初始化商品服务客户端
	goodsConn, err := createGRPCConn(consulInfo, global.ServerConfig.GoodsSrvInfo.Name)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 连接【商品服务】失败: %v", err)
	}
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	// 初始化库存服务客户端
	inventoryConn, err := createGRPCConn(consulInfo, global.ServerConfig.InventorySrvInfo.Name)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 连接【库存服务】失败: %v", err)
	}
	global.InventorySrvClient = proto.NewInventoryClient(inventoryConn)
}

func createGRPCConn(consulInfo config.ConsulConfig, srvName string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, srvName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	return conn, err
}
