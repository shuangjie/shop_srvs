package model

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;not null;comment:商品ID;index"`
	Stocks  int32 `gorm:"type:int;not null;comment:库存"`
	Version int32 `gorm:"type:int;not null;comment:版本-分布式乐观锁"`
}

//type InventoryHistory struct {
//	user   int32
//	goods  int32
//	nums   int32
//	order  int32
//	status int32 // 1.预扣减库存 2.已支付
//}
