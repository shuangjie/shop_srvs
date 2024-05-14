package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"srvs/goods_srv/global"
	"srvs/goods_srv/model"
	"srvs/goods_srv/proto"
)

// BrandList 获取品牌列表
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	brandListResponse := proto.BrandListResponse{}

	var brands []model.Brands
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)

	brandListResponse.Total = int32(total)

	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}

	brandListResponse.Data = brandResponses

	return &brandListResponse, nil

}

// CreateBrand 创建品牌
func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	brand := model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}

	if result := global.DB.First(&brand, "name = ?", req.Name); result.RowsAffected != 0 {
		return nil, status.Error(codes.InvalidArgument, "品牌已存在")
	}

	result := global.DB.Save(&brand)
	if result.Error != nil {
		return nil, result.Error
	}

	return &proto.BrandInfoResponse{
		Id:   brand.ID,
		Name: brand.Name,
		Logo: brand.Logo,
	}, nil
}

// DeleteBrand 删除品牌
func (s *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {

	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}
	return &emptypb.Empty{}, nil
}

// UpdateBrand 更新品牌
func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brands := model.Brands{}

	if result := global.DB.First(&brands, req.Id); result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}

	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}

	result := global.DB.Save(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}
