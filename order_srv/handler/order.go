package handler

import (
	"context"

	"srvs/goods_srv/global"
	"srvs/order_srv/model"
	"srvs/order_srv/proto"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

// CartItemList 获取用户购物车列表
func (*OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart
	var rsp proto.CartItemListResponse

	if result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}

	for _, cart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      cart.ID,
			UserId:  cart.User,
			GoodsId: cart.Goods,
			Nums:    cart.Nums,
			Checked: cart.Checked,
		})
	}

	return &rsp, nil
}

// CreateCartItem 添加购物车
func (*OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	// 判断商品是否存在： 1、不存在则新建 2、存在则合并
	var cart model.ShoppingCart

	if result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Goods: req.GoodsId}).First(&cart); result.RowsAffected == 1 {
		cart.Nums += req.Nums
	} else {
		cart = model.ShoppingCart{
			User:    req.UserId,
			Goods:   req.GoodsId,
			Nums:    req.Nums,
			Checked: false,
		}
	}
	global.DB.Save(&cart)
	return &proto.ShopCartInfoResponse{Id: cart.ID}, nil
}
