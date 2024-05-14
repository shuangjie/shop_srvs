package handler

import (
	"context"
	"srvs/goods_srv/global"
	"srvs/goods_srv/model"
	"srvs/goods_srv/proto"
)

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

/*
// 品牌和轮播图
BrandList(context.Context, *BrandFilterRequest) (*BrandListResponse, error)
CreateBrand(context.Context, *BrandRequest) (*BrandInfoResponse, error)
DeleteBrand(context.Context, *BrandRequest) (*emptypb.Empty, error)
UpdateBrand(context.Context, *BrandRequest) (*emptypb.Empty, error)
*/
