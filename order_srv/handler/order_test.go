// handler/order_test.go
package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"srvs/order_srv/proto"
	"testing"
)

var orderClient proto.OrderClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	orderClient = proto.NewOrderClient(conn)
}

func TestCreateCartItem(t *testing.T) {
	Init()
	defer conn.Close()

	tests := []struct {
		userId  int32
		goodsId int32
		nums    int32
		wantErr bool
	}{
		{userId: 1, goodsId: 421, nums: 1, wantErr: false},
		{userId: 2, goodsId: 422, nums: 2, wantErr: false},
		//{userId: 3, goodsId: 1, nums: 3, wantErr: true}, // 示例：此处预期会出错
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("用户ID=%d,商品ID=%d,数量=%d", tt.userId, tt.goodsId, tt.nums), func(t *testing.T) {
			rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
				UserId:  tt.userId,
				GoodsId: tt.goodsId,
				Nums:    tt.nums,
			})
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateCartItem 错误 = %v, 预期错误 %v", err, tt.wantErr)
			}
			if err == nil {
				fmt.Println("CreateCartItem 成功，ID:", rsp.Id)
			}
		})
	}
}
