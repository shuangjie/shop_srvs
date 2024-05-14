package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"

	"srvs/goods_srv/proto"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	// 初始化数据库
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	brandClient = proto.NewGoodsClient(conn)
}

func TestGetBrandsList() {
	rsp, err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pages:       1,
		PagePerNums: 2,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}

func main() {
	Init()

	TestGetBrandsList()

	conn.Close()
}
