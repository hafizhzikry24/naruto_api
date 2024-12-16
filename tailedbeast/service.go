package tailedbeast

import (
	"errors"

	"my-gin-app/models"

	"github.com/gosimple/slug"
)

type Service interface {
	CreateBeast(beast *models.TailedBeast) error
	GetBeastBySlug(slug string) (*models.TailedBeast, error)
	UpdateBeast(slug string, updatedData *models.TailedBeast) error
	DeleteBeast(slug string) error
	ListBeasts(page int, limit int) ([]models.TailedBeast, int64, error)
	SearchBeasts(name string) ([]models.TailedBeast, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateBeast(beast *models.TailedBeast) error {
	beast.Slug = slug.Make(beast.Name)
	return s.repo.Create(beast)
}

func (s *service) GetBeastBySlug(slug string) (*models.TailedBeast, error) {
	return s.repo.FindBySlug(slug)
}

func (s *service) UpdateBeast(slug string, updatedData *models.TailedBeast) error {
	existingBeast, err := s.repo.FindBySlug(slug)
	if err != nil {
		return err
	}

	// Update slug jika nama berubah
	if updatedData.Name != "" && updatedData.Name != existingBeast.Name {
		existingBeast.Slug = slug.Make(updatedData.Name)
	}

	// Update fields
	if updatedData.Name != "" {
		existingBeast.Name = updatedData.Name
	}
	if updatedData.Images != nil {
		existingBeast.Images = updatedData.Images
	}
	if updatedData.Rank != "" {
		existingBeast.Rank = updatedData.Rank
	}
	if updatedData.Abilities != nil {
		existingBeast.Abilities = updatedData.Abilities
	}
	if updatedData.Personality != "" {
		existingBeast.Personality = updatedData.Personality
	}

	update := map[string]interface{}{
		"name":        existingBeast.Name,
		"slug":        existingBeast.Slug,
		"images":      existingBeast.Images,
		"rank":        existingBeast.Rank,
		"abilities":   existingBeast.Abilities,
		"personality": existingBeast.Personality,
	}

	return s.repo.UpdateBySlug(slug, update)
}

func (s *service) DeleteBeast(slug string) error {
	return s.repo.DeleteBySlug(slug)
}

func (s *service) ListBeasts(page int, limit int) ([]models.TailedBeast, int64, error) {
	var skip int64
	var lim int64
	if limit > 0 && page > 0 {
		skip = int64((page - 1) * limit)
		lim = int64(limit)
	}

	beasts, err := s.repo.ListBeasts(skip, lim)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.CountBeasts()
	if err != nil {
		return nil, 0, err
	}

	return beasts, count, nil
}

func (s *service) SearchBeasts(name string) ([]models.TailedBeast, error) {
	if name == "" {
		return nil, errors.New("name query parameter is required")
	}

	beasts, _, err := s.repo.ListBeasts(0, 0)
	if err != nil {
		return nil, err
	}

	var filtered []models.TailedBeast
	for _, b := range beasts {
		if containsIgnoreCase(b.Name, name) {
			filtered = append(filtered, b)
		}
	}

	if len(filtered) == 0 {
		return nil, errors.New("no tailed beasts found")
	}

	return filtered, nil
}

func containsIgnoreCase(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(str) > 0 && (containsIgnoreCase(str[1:], substr) || str[:len(substr)] == substr))
}
