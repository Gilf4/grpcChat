package chatgrpc

import (
	"context"
	"slices"
	"sync"

	chatv1 "github.com/Gilf4/grpcChat/protos/gen/go/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Hub struct {
	mu      sync.RWMutex
	streams map[int64][]chan *chatv1.Message
}

func NewHub() *Hub {
	return &Hub{
		streams: make(map[int64][]chan *chatv1.Message),
	}
}

func (h *Hub) Subscribe(chatID int64) <-chan *chatv1.Message {
	ch := make(chan *chatv1.Message, 10)
	h.mu.Lock()
	h.streams[chatID] = append(h.streams[chatID], ch)
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(chatID int64, sub <-chan *chatv1.Message) {
	h.mu.Lock()
	defer h.mu.Unlock()
	subs := h.streams[chatID]
	for i, c := range subs {
		if c == sub {
			h.streams[chatID] = slices.Delete(subs, i, i+1)
			close(c)
			break
		}
	}
}

func (h *Hub) Broadcast(msg *chatv1.Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	subs := h.streams[msg.ChatId]
	for _, ch := range subs {
		select {
		case ch <- msg:
		default:
			continue
		}
	}
}

type serverApi struct {
	chatv1.UnimplementedChatServiceServer

	hub *Hub
	msg chan *chatv1.Message
}

func Register(gRPCServer *grpc.Server) {
	chatv1.RegisterChatServiceServer(gRPCServer, &serverApi{
		msg: make(chan *chatv1.Message, 100),
		hub: NewHub(),
	})
}

func (s *serverApi) ConnectChat(req *chatv1.ConnectChatRequest, stream chatv1.ChatService_ConnectChatServer) error {
	chatID := req.GetId()
	sub := s.hub.Subscribe(chatID)
	defer s.hub.Unsubscribe(chatID, sub)

	for {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended")
		case msg := <-sub:
			if err := stream.Send(msg); err != nil {
				return err
			}
		}
	}
}

func (s *serverApi) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest) (*emptypb.Empty, error) {
	msg := &chatv1.Message{
		ChatId:   req.GetChatId(),
		SenderId: req.GetSenderId(),
		Text:     req.GetText(),
	}

	s.hub.Broadcast(msg)

	return &emptypb.Empty{}, nil
}
