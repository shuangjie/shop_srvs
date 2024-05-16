package main

import (
	"context"
	"fmt"

	"srvs/goods_srv/proto"
)

func TestGetGoodsList() {
	rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		TopCategory: 130361,
		PriceMin:    90,
		PriceMax:    100,
		//KeyWords:    "深海速冻",
	})
	if err != nil {
		fmt.Printf("GoodsList failed, err:%v\n", err)
	}
	fmt.Println(rsp.Total)
	for _, good := range rsp.Data {
		fmt.Println(good.Name, good.ShopPrice)
	}

	//fmt.Println(rsp.Data)
}

func TestBatchGetGoods() {
	rsp, err := goodsClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: []int32{421, 422, 423},
	})
	if err != nil {
		fmt.Printf("BatchGetGoods failed, err:%v\n", err)
	}
	fmt.Println(rsp.Total)
	for _, good := range rsp.Data {
		fmt.Println(good.Name, good.ShopPrice)
	}
}

func TestGetGoodsDetail() {
	rsp, err := goodsClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: 421,
	})
	if err != nil {
		fmt.Printf("GetGoodsDetail failed, err:%v\n", err)
	}
	fmt.Println(rsp.Name, rsp.ShopPrice)
}
