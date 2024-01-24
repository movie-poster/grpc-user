package config

import (
	"flag"
	handler_user "grpc-user/cmd/handler/user"
	repository_user "grpc-user/internal/domain/repository/implement/user"
	"grpc-user/internal/domain/usecase"
	"grpc-user/internal/infra/proto/user"
	"grpc-user/internal/utils"

	"google.golang.org/grpc"
)

func init() {
	var configPath = ""
	configPath = *flag.String("config", "", "")

	if configPath == "" {
		configPath = "../data/config.yml"
	}

	setConfiguration(configPath)
}

func setConfiguration(configPath string) {
	Setup(configPath)

}

func Run(s *grpc.Server, configPath string) *grpc.Server {
	utils.SetupLoggerZap()
	conf := GetConfig()
	setupDB(conf)
	setupBrevoClient(conf)

	brevoService := usecase.BrevoService(GetBrevoClient())

	user.RegisterUserCrudServer(s, handler_user.NewServerUser(repository_user.UserRepository(DB, brevoService)))
	return s

}
