package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"

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

// UpdateCartItem 更新购物车, 可以更新数量和选中状态
func (*OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	var cart model.ShoppingCart

	// 判断购物车是否存在
	if result := global.DB.First(&cart, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车不存在")
	}

	// 更新购物车
	cart.Checked = req.Checked
	if req.Nums > 0 {
		cart.Nums = req.Nums
	}

	global.DB.Save(&cart)

	return &emptypb.Empty{}, nil
}

// DeleteCartItem 删除购物车
func (*OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.ShoppingCart{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车不存在")
	}

	return &emptypb.Empty{}, nil

}

// OrderList 获取订单列表
func (*OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp proto.OrderListResponse

	var total int64
	// 如果是后台查询，UserID 为 0，gorm 忽略此条件
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)

	// 分页
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&orders)
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SignerMobile,
		})
	}

	return &rsp, nil
}
