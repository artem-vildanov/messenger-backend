package dto

import "messenger/internal/domain/models"

type AuthRequest struct {
	Username string `json:"username" validate:"min=3,max=30"`
	Password string `json:"password" validate:"min=3,max=30"`
}

func (r *AuthRequest) ToDomain() *models.AuthModel {
	return &models.AuthModel{
		Username: r.Username,
		Password: r.Password,
	}
}
