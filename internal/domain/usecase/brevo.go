package usecase

import (
	"context"
	"grpc-user/internal/constant"
	"grpc-user/internal/domain/entity"
	sb "grpc-user/internal/domain/service"

	brevo "github.com/getbrevo/brevo-go/lib"
)

type senderBrevo struct {
	client *brevo.APIClient
}

func BrevoService(brevoClient *brevo.APIClient) sb.IBrevoSender {
	return &senderBrevo{
		client: brevoClient,
	}
}

func (r *senderBrevo) SendToNewUser(user entity.User, password string) error {
	ctx := context.Background()

	var params any = map[string]any{
		"NOMBRE":        user.Name,
		"TEMP_PASSWORD": password,
	}

	_, _, err := r.client.TransactionalEmailsApi.SendTransacEmail(ctx, brevo.SendSmtpEmail{
		To: []brevo.SendSmtpEmailTo{
			{
				Email: user.Email,
				Name:  user.Name,
			},
		},
		Params:     &params,
		TemplateId: constant.TemplateSendToNewUser,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *senderBrevo) ResetPasswordWithUsername(token, username, email, name string) error {
	ctx := context.Background()

	var params any = map[string]any{
		"TOKEN":    token,
		"USERNAME": username,
		"NAME":     name,
	}

	_, _, err := r.client.TransactionalEmailsApi.SendTransacEmail(ctx, brevo.SendSmtpEmail{
		To: []brevo.SendSmtpEmailTo{
			{
				Email: email,
				Name:  name,
			},
		},
		Params:     &params,
		TemplateId: constant.TemplatePasswordResetUsername,
	})
	if err != nil {
		return err
	}

	return nil
}
