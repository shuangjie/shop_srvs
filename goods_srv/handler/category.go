package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"google.golang.org/protobuf/types/known/emptypb"

	"srvs/goods_srv/global"
	"srvs/goods_srv/model"
	"srvs/goods_srv/proto"
)

func (g *GoodsServer) GetAllCategorysList(ctx context.Context, req *emptypb.Empty) (*proto.CategoryListResponse, error) {
	//返回 id、name、level、parent_id、is_tab、sub_categories
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	//b, _ := json.Marshal(&categorys)
	//return &proto.CategoryListResponse{JsonData: string(b)}, nil

	// 初始化输出流
	var buf bytes.Buffer

	// 创建流式编码器，并将数据编码到输出流中
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(&categorys); err != nil {
		return nil, err
	}

	// 将输出流中的数据转换为字符串，并存储在 proto.CategoryListResponse 中
	return &proto.CategoryListResponse{JsonData: buf.String()}, nil
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
