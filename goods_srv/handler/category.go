package handler

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"

	"srvs/goods_srv/global"
	"srvs/goods_srv/model"
	"srvs/goods_srv/proto"
)

// GetAllCategorysList 获取所有分类
func (g *GoodsServer) GetAllCategorysList(ctx context.Context, req *emptypb.Empty) (*proto.CategoryListResponse, error) {
	//返回 id、name、level、parent_id、is_tab、sub_categories
	var categories []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categories)
	b, _ := json.Marshal(&categories)
	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

// GetSubCategory 获取子分类
func (g *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListResponse := proto.SubCategoryListResponse{}

	var category model.Category
	if err := global.DB.First(&category, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}

	var subCategories []model.Category
	var subCategoriesResponse []*proto.CategoryInfoResponse

	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategories)

	for _, subCategory := range subCategories {
		subCategoriesResponse = append(subCategoriesResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}

	categoryListResponse.SubCategorys = subCategoriesResponse
	return &categoryListResponse, nil
}

// CreateCategory 创建分类
func (g *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	var category model.Category

	if category.Level != 1 {
		// 查询父分类是否存在, 前端如果是拉去的数据（下拉菜单），可以不用判断
		if err := global.DB.First(&category, category.ParentCategoryID); err != nil {
			return nil, status.Errorf(codes.NotFound, "父分类不存在")
		}
		category.ParentCategoryID = req.ParentCategory
	}
	category.Name = req.Name
	category.Level = req.Level
	category.IsTab = req.IsTab

	if err := global.DB.Create(&category).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "创建商品分类失败")
	}

	return &proto.CategoryInfoResponse{Id: category.ID}, nil
}

// DeleteCategory 删除分类
func (g *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	return &emptypb.Empty{}, nil
}

// UpdateCategory 更新分类
func (g *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}

	if err := global.DB.Save(&category).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "更新商品分类失败")
	}

	return &emptypb.Empty{}, nil
}
