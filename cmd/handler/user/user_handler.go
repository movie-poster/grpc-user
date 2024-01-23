package handler_user

import (
	"context"
	"errors"
	"grpc-user/internal/constant"
	cryptopassword "grpc-user/internal/crypto"
	"grpc-user/internal/domain/entity"
	irepository "grpc-user/internal/domain/repository/interface"
	"grpc-user/internal/utils"

	pb "grpc-user/internal/infra/proto/user"
)

func NewServerUser(crud irepository.IUserCrud) *server {
	return &server{
		crud: crud,
	}
}

type server struct {
	crud irepository.IUserCrud
	pb.UnimplementedUserCrudServer
}

func (s *server) Insert(context context.Context, user *pb.User) (*pb.ResponseUser, error) {
	password, hashedPassword := utils.CreateTemporalPassword()

	userObject := &entity.User{
		Name:     user.GetName(),
		Email:    user.GetEmail(),
		NickName: user.GetNickName(),
		IsAdmin:  constant.NoIsAdmin,
		State:    constant.ActiveState,
		Surname:  user.GetSurname(),
		Password: hashedPassword,
	}

	response := s.crud.Insert(userObject, password)

	responsePB := pb.ResponseUser{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	if response.IsOk {
		responsePB.User = response.Value.(*pb.User)
	}

	return &responsePB, nil
}

func (s *server) Update(context context.Context, user *pb.User) (*pb.ResponseUser, error) {
	userObject := &entity.User{
		Name:    user.GetName(),
		Surname: user.GetSurname(),
	}

	response := s.crud.Update(userObject)
	responsePB := pb.ResponseUser{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	if response.IsOk {
		responsePB.User = response.Value.(*pb.User)
	}

	return &responsePB, nil
}

func (s *server) Delete(context context.Context, req *pb.RequestByIdUser) (*pb.ResponseUser, error) {
	response := s.crud.Delete(req.GetId())

	responsePB := pb.ResponseUser{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	return &responsePB, nil
}

func (s *server) Login(ctx context.Context, in *pb.LoginUser) (*pb.User, error) {
	response := s.crud.VerifyUser(in.GetNickName())

	if response.IsOk {
		loginResponse := response.Value.(*pb.User)
		if utils.ValidatePassword(in.GetPassword(), loginResponse.Password) {
			return loginResponse, nil
		}
	}

	return nil, errors.New("usuario o contraseña inválidos")
}

func (s *server) ResetPassword(context context.Context, req *pb.ResetPasswordRequest) (*pb.ResponseUser, error) {
	reset := entity.RecoverPassword{
		Token:    req.GetToken(),
		Nickname: req.GetNickname(),
		State:    constant.ActiveState,
	}

	response := s.crud.ResetPassword(reset)
	responsePB := &pb.ResponseUser{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	return responsePB, nil
}

func (s *server) CheckToken(context context.Context, req *pb.CheckTokenRequest) (*pb.ResponseUser, error) {
	response := s.crud.CheckToken(req.GetToken())
	responsePB := &pb.ResponseUser{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	return responsePB, nil
}

func (s *server) ChangePassword(context context.Context, req *pb.ChangePasswordRequest) (*pb.ResponseUser, error) {
	hashedPassword, _ := cryptopassword.HashAndSalt([]byte(req.GetPassword()))

	response := s.crud.ChangePassword(req.GetToken(), hashedPassword, req.GetNickName())
	responsePB := &pb.ResponseUser{
		Title:   response.Title,
		IsOk:    response.IsOk,
		Message: response.Message,
		Status:  response.Status,
	}

	return responsePB, nil
}
