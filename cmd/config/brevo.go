package config

import (
	brevo "github.com/getbrevo/brevo-go/lib"
)

var BrevoClient *brevo.APIClient

func setupBrevoClient(configuration *Configuration) {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", configuration.Brevo.ApiKey)
	BrevoClient = brevo.NewAPIClient(cfg)
}

func GetBrevoClient() *brevo.APIClient {
	return BrevoClient
}
