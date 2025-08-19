package convert

import (
	"github.com/Gilf4/grpcChat/chat/internal/domain/models"
	chatv1 "github.com/Gilf4/grpcChat/protos/gen/go/chat/v1"
)

func ChatToProto(chat models.Chat) *chatv1.Chat {
	return &chatv1.Chat{
		Id:   chat.ID,
		Name: chat.Name,
	}
}

func ProtoToChat(proto *chatv1.Chat) *models.Chat {
	return &models.Chat{
		ID:   proto.Id,
		Name: proto.Name,
	}
}

func ToProtoChatList(chats []models.Chat) []*chatv1.Chat {
	res := make([]*chatv1.Chat, 0, len(chats))
	for _, c := range chats {
		res = append(res, ChatToProto(c))
	}
	return res
}
