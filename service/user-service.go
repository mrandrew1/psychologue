package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/mrandrew1/psychologue/dto"
	"github.com/mrandrew1/psychologue/entity"
	"github.com/mrandrew1/psychologue/repository"
)

type UserService interface {
	Update(user dto.UserUpdateDTO) entity.User
	Get(userID string) entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) Update(user dto.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("failed to map %v: ", err)
	}
	updatedUser := service.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (service *userService) Get(userID string) entity.User {
	return service.userRepository.GetUser(userID)
}
