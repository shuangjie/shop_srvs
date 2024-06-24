package handler

import (
	"srvs/userop_srv/proto"
)

type UserOpServer struct {
	proto.UnimplementedMessageServer
	proto.UnimplementedAddressServer
	proto.UnimplementedUserFavServer
}
