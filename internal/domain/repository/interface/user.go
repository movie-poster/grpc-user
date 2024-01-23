package irepository

import (
	"grpc-user/internal/domain/entity"
	objectvalue "grpc-user/internal/domain/object-value"
)

type IUserCrud interface {
	Insert(user *entity.User, originalPassword string) *objectvalue.ResponseValue
	Update(*entity.User) *objectvalue.ResponseValue
	Delete(ID uint64) *objectvalue.ResponseValue
	VerifyUser(nickname string) *objectvalue.ResponseValue
	ResetPassword(entity.RecoverPassword) *objectvalue.ResponseValue
	ChangePassword(token, password, nickName string) *objectvalue.ResponseValue
	CheckToken(token string) *objectvalue.ResponseValue
}
