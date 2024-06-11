package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"srvs/inventory_srv/global"
	"srvs/inventory_srv/model"
	"srvs/inventory_srv/proto"
	"sync"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

// SetInv 设置库存
func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

// InvDetail 获取库存详情
func (*InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有找到库存信息")
	}

	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

// 锁
var m sync.Mutex

// Sell 扣减库存
func (c *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 开启事务
	tx := global.DB.Begin()
	//m.Lock()
	for _, good := range req.GoodsInfo {
		var inv model.Inventory
		/*if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: good.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有找到库存信息")
		}*/

		for {
			if result := global.DB.Where(&model.Inventory{Goods: good.GoodsId}).First(&inv); result.RowsAffected == 0 {
				tx.Rollback()
				return nil, status.Errorf(codes.InvalidArgument, "没有找到库存信息")
			}

			if inv.Stocks < good.Num {
				tx.Rollback()
				return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
			}
			// 扣减库存
			inv.Stocks -= good.Num

			//update inventory set stocks = stocks - 1 and version = version + 1 where goods = goods and version = version
			if result := tx.Model(&model.Inventory{}).Select("stocks", "version").Where("goods = ? and version = ?", good.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
				zap.S().Info("扣减库存失败")
			} else {
				break
			}
		}

		//tx.Save(&inv)
	}
	tx.Commit()
	//m.Unlock()
	return &emptypb.Empty{}, nil
}

// ReBack 回退库存
/*
 * 库存回退场景：
 * 1. 订单超时未支付
 * 2. 订单创建失败
 * 3. 订单取消 （手动）
 */
func (c *InventoryServer) ReBack(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 开启事务
	tx := global.DB.Begin()
	m.Lock()
	for _, good := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: good.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有找到库存信息")
		}
		// 回退库存
		inv.Stocks += good.Num
		tx.Save(&inv)
	}
	tx.Commit()
	m.Unlock()
	return &emptypb.Empty{}, nil
}
