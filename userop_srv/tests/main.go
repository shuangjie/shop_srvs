package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"srvs/userop_srv/proto"
)

var (
	addressClient proto.AddressClient
	messageClient proto.MessageClient
	userFavClient proto.UserFavClient
	conn          *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	addressClient = proto.NewAddressClient(conn)
	messageClient = proto.NewMessageClient(conn)
	userFavClient = proto.NewUserFavClient(conn)
}

func TestAddressList() {
	rsp, err := addressClient.GetAddressList(context.Background(), &proto.AddressRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
	for _, addr := range rsp.Data {
		fmt.Println(addr.Address)
	}
}

func TestMessageList() {
	rsp, err := messageClient.MessageList(context.Background(), &proto.MessageRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
	for _, msg := range rsp.Data {
		fmt.Println(msg.Message)
	}
}

func TestUserFavList() {
	rsp, err := userFavClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
	for _, fav := range rsp.Data {
		fmt.Printf("UserID: %d, GoodsID: %d\n", fav.UserId, fav.GoodsId)
	}
}

func main() {
	Init()
	TestAddressList()
	TestMessageList()
	TestUserFavList()
	conn.Close()
}
