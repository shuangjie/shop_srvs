package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"srvs/userop_srv/proto"
)

type GRPCTestSuite struct {
	suite.Suite
	conn          *grpc.ClientConn
	addressClient proto.AddressClient
	messageClient proto.MessageClient
	userFavClient proto.UserFavClient
}

func (suite *GRPCTestSuite) SetupSuite() {
	var err error
	suite.conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		suite.T().Fatalf("failed to connect to gRPC server: %v", err)
	}
	suite.addressClient = proto.NewAddressClient(suite.conn)
	suite.messageClient = proto.NewMessageClient(suite.conn)
	suite.userFavClient = proto.NewUserFavClient(suite.conn)
}

func (suite *GRPCTestSuite) TearDownSuite() {
	suite.conn.Close()
}

func (suite *GRPCTestSuite) TestAddressList() {
	rsp, err := suite.addressClient.GetAddressList(context.Background(), &proto.AddressRequest{
		UserId: 1,
	})
	assert.NoError(suite.T(), err)
	for _, addr := range rsp.Data {
		fmt.Println(addr.Address)
	}
}

func (suite *GRPCTestSuite) TestMessageList() {
	rsp, err := suite.messageClient.MessageList(context.Background(), &proto.MessageRequest{
		UserId: 1,
	})
	assert.NoError(suite.T(), err)
	for _, msg := range rsp.Data {
		fmt.Println(msg.Message)
	}
}

func (suite *GRPCTestSuite) TestUserFavList() {
	rsp, err := suite.userFavClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: 1,
	})
	assert.NoError(suite.T(), err)
	for _, fav := range rsp.Data {
		fmt.Printf("UserID: %d, GoodsID: %d\n", fav.UserId, fav.GoodsId)
	}
}

func TestGRPCTestSuite(t *testing.T) {
	suite.Run(t, new(GRPCTestSuite))
}
