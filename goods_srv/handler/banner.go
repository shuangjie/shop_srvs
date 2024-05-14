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

// BannerList 轮播图列表
func (s *GoodsServer) BannerList(ctx context.Context, req *emptypb.Empty) (rsp *proto.BannerListResponse, err error) {
	bannerListResponse := proto.BannerListResponse{}

	var banners []model.Banner
	result := global.DB.Find(&banners)
	bannerListResponse.Total = int32(result.RowsAffected)

	var bannerResponses []*proto.BannerResponse
	for _, banner := range banners {
		bannerResponses = append(bannerResponses, &proto.BannerResponse{
			Id:    banner.ID,
			Image: banner.Image,
			Url:   banner.Url,
			Index: banner.Index,
		})
	}
	bannerListResponse.Data = bannerResponses

	return &bannerListResponse, nil
}

// CreateBanner 创建轮播图
func (s *GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (rsp *proto.BannerResponse, err error) {
	banner := model.Banner{
		Image: req.Image,
		Url:   req.Url,
		Index: req.Index,
	}
	global.DB.Save(&banner)

	return &proto.BannerResponse{
		Id:    banner.ID,
		Image: banner.Image,
		Url:   banner.Url,
		Index: banner.Index,
	}, nil
}

// DeleteBanner 删除轮播图
func (s *GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (rsp *emptypb.Empty, err error) {
	if result := global.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	return &emptypb.Empty{}, nil
}

// UpdateBanner 更新轮播图
func (s *GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (rsp *emptypb.Empty, err error) {
	var banner model.Banner
	if result := global.DB.First(&banner, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}

	global.DB.Save(&banner)

	return &emptypb.Empty{}, nil
}
