package services

import (
  "fiberapi/internal/core/ports"
)

type AdminService struct {
  adminRepository ports.AdminRepository
}

//This line is for get feedback in case we are not implementing the interface correctly
var _ ports.AdminService = (*AdminService)(nil)

func NewAdminService(repository ports.AdminService) *AdminService {
  return &AdminService{
	adminRepository: repository,
  }
}

func (s *AdminService) LoginAdmin(email string, password string) (bool, error) {
  boolean,err := s.adminRepository.LoginAdmin(email, password)
  if err != nil {
	return false,err
  }
  if boolean{
    return true, nil
  }
  return false, nil
  
}

func (s *AdminService) GetDashboard() ([]ports.Offer, []ports.OrderStatus, int, error){
	return s.adminRepository.GetDashboard()
}

func (s *AdminService) PatchStatus(id int, status string) error{
	return s.adminRepository.PatchStatus(id, status)
}
