package model

import (
	protoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
)

func SessionIDFromUsecaseToProto(sessionID string) *protoModel.SessionID {
	return &protoModel.SessionID{
		SessionId: sessionID,
	}
}

func NothingFromUsecaseToProto() *protoModel.Nothing {
	return &protoModel.Nothing{Dummy: true}
}

func UserIDFromUsecaseToProto(userID int64) *protoModel.UserID {
	return &protoModel.UserID{
		Id: userID,
	}
}
