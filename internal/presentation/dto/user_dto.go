package dto

import "messenger/internal/domain/models"

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func NewUserResponse(userModel *models.UserModel) *UserResponse {
	return &UserResponse{
		Id:       userModel.Id,
		Username: userModel.Username,
	}
}

func NewMultipleUsersResponse(usersModels []*models.UserModel) []*UserResponse {
	usersResponses := make([]*UserResponse, 0, len(usersModels))
	for _, userModel := range usersModels {
		usersResponses = append(usersResponses, NewUserResponse(userModel))
	}
	return usersResponses
}
