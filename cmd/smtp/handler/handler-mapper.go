package handler

import (
	"brickwall/internal/common"
	"brickwall/internal/provider"
)

var HandlerMapper = map[string]provider.MessageHandler{
	string(common.TopicUserRegistrationEmail):  UserRegistrationEmail,
	string(common.TopicUserResetPasswordEmail): UserResetPasswordEmail,
}
