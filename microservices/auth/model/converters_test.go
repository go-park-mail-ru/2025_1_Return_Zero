package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	protoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
)

func TestConvertUserIDToProto(t *testing.T) {
	userID := int64(42)
	expectedProtoUserID := &protoModel.UserID{Id: userID}

	assert.Equal(t, expectedProtoUserID, UserIDFromUsecaseToProto(userID))
}

func TestConvertSessionIDToProto(t *testing.T) {
	sessionID := "session-123"
	expectedProtoSessionID := &protoModel.SessionID{SessionId: sessionID}

	assert.Equal(t, expectedProtoSessionID, SessionIDFromUsecaseToProto(sessionID))
}

func TestConvertNothingToProto(t *testing.T) {
	nothing := NothingFromUsecaseToProto()

	assert.NotNil(t, nothing)
	assert.True(t, nothing.Dummy)
}