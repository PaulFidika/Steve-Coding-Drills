package main

import (
	"context"
	chat "steve-fidika/2024-10-12/server/pb/chat"
)

type chatServer struct {
	chat.UnimplementedChatServiceServer
}

func (c *chatServer) SayHello(context context.Context, msg *chat.Message) (*chat.Message, error) {
	return msg, nil
	// return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}