package main

import (
	"google.golang.org/grpc"
	"srvs/goods_srv/proto"
)

var goodsClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	// 初始化数据库
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	goodsClient = proto.NewGoodsClient(conn)
}

func main() {
	Init()

	//TestGetBrandsList()

	//TestGetCategoryList()

	//TestGetSubCategories()

	TestGetCategoryBrand()

	conn.Close()
}
