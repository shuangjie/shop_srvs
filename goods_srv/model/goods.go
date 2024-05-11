package model

type Category struct {
	BaseModel
	Name             string `gorm:"type:varchar(50);not null"`
	ParentCategoryID int32
	ParentCategory   *Category
	Level            int32 `gorm:"type:int;not null;default:1"`
	IsTab            bool  `gorm:"type:tinyint;not null;default:0"`
}

type Brand struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);not null;default:''"`
}

// GoodsCategoryBrand 可以用 gorm 的 many2many 多对多关联来自动生成
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   *Category

	BrandID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brand   *Brand
}

func (GoodsCategoryBrand) TableName() string {
	return "goods_category_brand"
}

// Banner Image 的长度为 200，实际上可能会有问题，可以根据实际情况调整
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;not null;default:1"`
}
