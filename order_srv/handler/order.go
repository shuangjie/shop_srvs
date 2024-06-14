package handler

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/rand"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"srvs/order_srv/global"
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
	// 1. 判断商品是否存在：
	if _, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: req.GoodsId}); err != nil {
		return nil, status.Error(codes.NotFound, "商品不存在")
	}

	// 2. 判断商品是否存在于购物车： A、不存在则新建 B、存在则合并
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
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
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

// OrderDetail 获取订单详情
func (*OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp proto.OrderInfoDetailResponse

	// 如果是后台查询，UserID 为 0，gorm 忽略此条件
	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	// 获取订单详情
	rsp.OrderInfo = &proto.OrderInfoResponse{
		Id:      order.ID,
		UserId:  order.User,
		OrderSn: order.OrderSn,
		PayType: order.PayType,
		Status:  order.Status,
		Total:   order.OrderMount,
		Address: order.Address,
		Name:    order.SignerName,
		Mobile:  order.SignerMobile,
	}

	// 获取订单商品
	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}
	for _, goods := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId:    goods.Goods,
			GoodsName:  goods.GoodsName,
			GoodsPrice: goods.GoodsPrice,
			Nums:       goods.Nums,
		})
	}

	return &rsp, nil
}

// CreateOrder 创建订单
func (*OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		创建订单
		   A. 计算订单总价 - 访问商品服务获取商品价格
		   B. 库存扣减 - 访问商品服务扣减库存
		   C. 订单的基本信息 + 订单的商品信息（从购物车获取选中的商品）
		   D. 从购物车删除已购买的记录
	*/
	// 1. 从购物车获取选中的商品
	var shopCarts []model.ShoppingCart
	var goodsIds []int32
	var goodsNumsMap = make(map[int32]int32)
	if result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "没有选中结算的商品")
	}

	for _, cart := range shopCarts {
		goodsIds = append(goodsIds, cart.Goods)
		goodsNumsMap[cart.Goods] = cart.Nums
	}

	// 2. 访问商品服务获取商品价格。 跨服务调用
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		return nil, status.Error(codes.Internal, "获取商品信息失败")
	}

	// 3. 计算订单总价
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsPrice: good.ShopPrice,
			GoodsImage: good.GoodsFrontImage,
			Nums:       goodsNumsMap[good.Id],
		})

		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}

	// 4. 访问库存服务扣减库存
	if _, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{GoodsInfo: goodsInvInfo}); err != nil {
		return nil, status.Error(codes.ResourceExhausted, "扣减库存失败")
	}

	// 5. 创建订单 - 开启事务 todo: 应该开启分布式事务
	tx := global.DB.Begin()
	order := model.OrderInfo{
		User:         req.UserId,
		OrderSn:      GenerateOrderSn(req.UserId),
		OrderMount:   orderAmount,
		Address:      req.Address,
		SignerName:   req.Name,
		SignerMobile: req.Mobile,
		Post:         req.Post,
	}

	if result := tx.Save(&order); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "创建订单失败")
	}

	for _, orderGood := range orderGoods {
		orderGood.Order = order.ID
	}
	// 6. 批量插入订单商品
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "创建订单失败")
	}

	// 7. 删除购物车中已购买的商品
	if result := tx.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "创建订单失败")
	}

	return &proto.OrderInfoResponse{Id: order.ID, OrderSn: order.OrderSn, Total: orderAmount}, tx.Commit().Error
}

// UpdateOrderStatus 更新订单状态
func (*OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	return &emptypb.Empty{}, nil
}

// GenerateOrderSn 生成唯一订单号,但这里的用户id会导致订单号长度不一致
/*func GenerateOrderSn(userId int32) string {
	// 规则: 年月日时分秒(nano) + 用户ID + 2位随机数
	now := time.Now()
	rand.Seed(uint64(now.UnixNano()))
	return fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId,
		rand.Intn(90)+10,
	)
}*/

// GenerateOrderSn 生成唯一订单号
func GenerateOrderSn(userId int32) string {
	// 规则: 1位订单类型+13位毫秒级别时间戳+2位自增值+2位随机数
	now := time.Now()
	rand.Seed(uint64(now.UnixNano()))
	millisecond := now.UnixNano()
	orderType := 1
	orderIndex := GetOrderIndex(int(millisecond), userId)
	return fmt.Sprintf("%d%d%s%d", orderType, millisecond, orderIndex, rand.Intn(90)+10)

}

// GetOrderIndex 获取订单索引,通过redis的原子性自增操作，获得一毫秒内新增订单的个数预留2位
func GetOrderIndex(millisecond int, userId int32) string {
	key := GetOrderIndexKey(millisecond, userId)
	index := global.RedisCli.Incr(context.Background(), key).Val()
	if index == 1 {
		global.RedisCli.Expire(context.Background(), key, time.Second*1)
	}
	return fmt.Sprintf("%02d", index)
}

func GetOrderIndexKey(millisecond int, userId int32) string {
	return fmt.Sprintf("order:index:user:%d:%d", userId, millisecond)
}
