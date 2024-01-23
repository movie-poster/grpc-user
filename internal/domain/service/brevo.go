package service

import (
	"grpc-user/internal/domain/entity"
)

type NickNameCompany struct {
	NickName    string `json:"NickName"`
	CompanyName string `json:"CompanyName"`
}

type IBrevoSender interface {
	SendToNewUser(user entity.User, password string) error
	ResetPasswordWithUsername(token, username, email, name string) error
}
