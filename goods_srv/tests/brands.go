package main

import (
	"context"
	"fmt"
	"srvs/goods_srv/proto"
)

func TestGetBrandsList() {
	rsp, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{
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
