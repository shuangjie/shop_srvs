package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"srvs/order_srv/proto"
)

var orderClient proto.OrderClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	orderClient = proto.NewOrderClient(conn)
}

func TestCreateCartItem() {
	rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  1,
		GoodsId: 421,
		Nums:    1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("CreateCartItem success, ID:", rsp.Id)
}

func main() {
	Init()

	TestCreateCartItem()

	conn.Close()
}
