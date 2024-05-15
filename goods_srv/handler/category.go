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
	//preloads := "SubCategory"
	//if category.Level == 1 {
	//	preloads = "SubCategory.SubCategory"
	//}
	//global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload(preloads).Find(&subCategories)
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

/*
// 商品分类
GetAllCategorysList(context.Context, *emptypb.Empty) (*CategoryListResponse, error)
// 获取子分类
GetSubCategory(context.Context, *CategoryListRequest) (*SubCategoryListResponse, error)
CreateCategory(context.Context, *CategoryInfoRequest) (*CategoryInfoResponse, error)
DeleteCategory(context.Context, *DeleteCategoryRequest) (*emptypb.Empty, error)
UpdateCategory(context.Context, *CategoryInfoRequest) (*emptypb.Empty, error)
*/
