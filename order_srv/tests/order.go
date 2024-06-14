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

func TestCreateCartItem(userId int32, goodsId int32, nums int32) {
	rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId,
		GoodsId: goodsId,
		Nums:    nums,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("CreateCartItem success, ID:", rsp.Id)
}

func TestCartItemList(userId int32) {
	rsp, err := orderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: userId,
	})
	if err != nil {
		panic(err)
	}

	for _, item := range rsp.Data {
		fmt.Println(item.Id, item.GoodsId, item.Nums)
	}
}

func TestUpdateCartItem(id int32, nums int32) {
	_, err := orderClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      id,
		Nums:    nums,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateOrder(userId int32) {
	rsp, err := orderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  userId,
		Address: "广东省广州市",
		Name:    fmt.Sprintf("Zeng#%d", userId),
		Mobile:  "18888888888",
		Post:    "请尽快发货",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("CreateOrder success, ID:", rsp.Id)
}

func TestOrderDetail(orderId int32) {
	rsp, err := orderClient.OrderDetail(context.Background(), &proto.OrderRequest{
		Id: orderId,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("OrderSn:", rsp.OrderInfo.OrderSn)
	for _, item := range rsp.Goods {
		fmt.Println(item.GoodsName, item.GoodsPrice, item.Nums)
	}
}

func TestOrderList() {
	rsp, err := orderClient.OrderList(context.Background(), &proto.OrderFilterRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}

	for _, order := range rsp.Data {
		fmt.Println(order.OrderSn, order.Total)
	}

}

func main() {
	Init()

	//TestCreateCartItem(1, 422, 2)
	//TestCartItemList(1)
	//TestUpdateCartItem(3, 0)
	//TestCreateOrder(1)
	//TestOrderDetail(5)
	TestOrderList()

	conn.Close()
}
