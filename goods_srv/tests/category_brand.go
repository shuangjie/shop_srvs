package main

import (
	"context"
	"fmt"
	"srvs/goods_srv/proto"
)

func TestGetCategoryBrand() {
	rsp, err := goodsClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.Data)
}
