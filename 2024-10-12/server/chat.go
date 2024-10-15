package main

import (
	"fmt"
	"io"
	chat "steve-fidika/2024-10-12/server/pb/chat"

	"google.golang.org/grpc"
)

type chatServer struct {
	chat.UnimplementedChatServiceServer
}

// func (c *chatServer) SayHello(context context.Context, msg *chat.Message) (*chat.Message, error) {
// 	return msg, nil
// 	// return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
// }

func (c *chatServer) SayHello(stream grpc.BidiStreamingServer[chat.Message, chat.Message]) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &chat.Message{
			Body: "Hello there! Re: " + in.GetBody(),
		}

		err = stream.Send(reply)
		if err != nil {
			fmt.Println("Unable to send response")
		}
	}
	
	// return status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}