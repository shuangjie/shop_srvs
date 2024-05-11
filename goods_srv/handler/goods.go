package handler

import (
	"srvs/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

//func (s *GoodsServer) GoodsList(context.Context, *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
//	return nil, nil
//}

/*
// 商品接口
GoodsList(context.Context, *GoodsFilterRequest) (*GoodsListResponse, error)
// 现在用户提交订单有多个商品，你得批量查询商品的信息吧
BatchGetGoods(context.Context, *BatchGoodsIdInfo) (*GoodsListResponse, error)
CreateGoods(context.Context, *CreateGoodsInfo) (*GoodsInfoResponse, error)
DeleteGoods(context.Context, *DeleteGoodsInfo) (*emptypb.Empty, error)
UpdateGoods(context.Context, *CreateGoodsInfo) (*emptypb.Empty, error)
GetGoodsDetail(context.Context, *GoodInfoRequest) (*GoodsInfoResponse, error)
*/
