package services

import (
  "errors"
  "fiberapi/internal/core/ports"
)

type UserService struct {
  userRepository ports.UserRepository
}

//This line is for get feedback in case we are not implementing the interface correctly
var _ ports.UserService = (*UserService)(nil)

func NewUserService(repository ports.UserRepository) *UserService {
  return &UserService{
	userRepository: repository,
  }
}

func (s *UserService) Login(email string, password string) (bool, error) {
  boolean,err := s.userRepository.Login(email, password)
  if err != nil {
	return false,err
  }
  if boolean{
    return true, nil
  }
  return false, nil
  
}

func (s *UserService) Register(username string, email string, password string, confirmPass string) error {
  if password != confirmPass {
	return errors.New("the passwords are not equal")
  }
  err := s.userRepository.Register(username,email, password)
  if err != nil {
	return err
  }
  return nil
}

func (s *UserService) GetOffers() ([]ports.Offer, error) {
  return s.userRepository.GetOffers()
}

func (s *UserService) Checkout(order ports.Order) (int, string, error) {
  return s.userRepository.Checkout(order)
}

func (s *UserService) GetStatus(id int) (string, error){
  return s.userRepository.GetStatus(id)
}

func (s *UserService) UpdateOffers(food, medicine map[string]string) error{
  return s.userRepository.UpdateOffers(food, medicine)
}
