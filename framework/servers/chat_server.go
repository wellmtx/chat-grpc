package servers

import (
	"chat/framework/pb"
	"fmt"
)

type ChatServer struct {
	channel map[string][]chan *pb.Message
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		channel: make(map[string][]chan *pb.Message),
	}
}

func (chatServer *ChatServer) JoinChannel(ch *pb.Channel, msgStream pb.ChatService_JoinChannelServer) error {
	msgChannel := make(chan *pb.Message)

	chatServer.channel[ch.Name] = append(chatServer.channel[ch.Name], msgChannel)

	for {
		select {
		case <-msgStream.Context().Done():
			return nil
		case msg := <-msgChannel:
			fmt.Printf("GO ROUTINE (got message): %v \n", msg)
			msgStream.Send(msg)
		}
	}
}

func (chatServer *ChatServer) SendMessage(msgStream pb.ChatService_SendMessageServer) error {
	msg, err := msgStream.Recv()

	if err != nil {
		return err
	}

	ack := pb.MessageAck{Status: "SENT"}

	msgStream.SendAndClose(&ack)

	go func() {
		streams := chatServer.channel[msg.Channel.Name]

		for _, msgChan := range streams {
			msgChan <- msg
		}
	}()

	return nil
}
