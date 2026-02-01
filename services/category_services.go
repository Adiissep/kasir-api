package services

import (
	"fmt"
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.GetAllCategories()
}

func (s *CategoryService) CreateCategory(c *models.Category) error {
	if c.Name == "" {
		return fmt.Errorf("category name is required")
	}
	return s.repo.CreateCategory(c)
}

func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *CategoryService) UpdateCategory(c *models.Category) error {
	if c.ID == 0 {
		return fmt.Errorf("invalid category ID")
	}
	return s.repo.UpdateCategory(c)
}

func (s *CategoryService) DeleteCategory(id int) error {
	return s.repo.DeleteCategory(id)
}
