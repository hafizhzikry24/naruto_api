package character

import (
	"errors"

	"my-gin-app/models"

	"github.com/gosimple/slug"
)

type Service interface {
	CreateCharacter(character *models.Character) error
	GetCharacterBySlug(slug string) (*models.Character, error)
	UpdateCharacter(slug string, updatedData *models.Character) error
	DeleteCharacter(slug string) error
	ListCharacters(page int, limit int) ([]models.Character, int64, error)
	SearchCharacters(name string) ([]models.Character, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateCharacter(character *models.Character) error {
	// Validasi atau logika bisnis tambahan dapat ditambahkan di sini
	character.Slug = slug.Make(character.Name)
	return s.repo.Create(character)
}

func (s *service) GetCharacterBySlug(slug string) (*models.Character, error) {
	return s.repo.FindBySlug(slug)
}

func (s *service) UpdateCharacter(slug string, updatedData *models.Character) error {
	existingCharacter, err := s.repo.FindBySlug(slug)
	if err != nil {
		return err
	}

	// Update slug jika nama berubah
	if updatedData.Name != "" && updatedData.Name != existingCharacter.Name {
		existingCharacter.Slug = slug.Make(updatedData.Name)
	}

	// Update fields
	if updatedData.Name != "" {
		existingCharacter.Name = updatedData.Name
	}
	if updatedData.Images != nil {
		existingCharacter.Images = updatedData.Images
	}
	if updatedData.Personal != (models.Personal{}) {
		existingCharacter.Personal = updatedData.Personal
	}
	if updatedData.Rank != (models.Rank{}) {
		existingCharacter.Rank = updatedData.Rank
	}
	if updatedData.Debut != (models.Debut{}) {
		existingCharacter.Debut = updatedData.Debut
	}
	if updatedData.Jutsu != nil {
		existingCharacter.Jutsu = updatedData.Jutsu
	}

	update := map[string]interface{}{
		"name":     existingCharacter.Name,
		"slug":     existingCharacter.Slug,
		"images":   existingCharacter.Images,
		"personal": existingCharacter.Personal,
		"rank":     existingCharacter.Rank,
		"debut":    existingCharacter.Debut,
		"jutsu":    existingCharacter.Jutsu,
	}

	return s.repo.UpdateBySlug(slug, update)
}

func (s *service) DeleteCharacter(slug string) error {
	return s.repo.DeleteBySlug(slug)
}

func (s *service) ListCharacters(page int, limit int) ([]models.Character, int64, error) {
	var skip int64
	var lim int64
	if limit > 0 && page > 0 {
		skip = int64((page - 1) * limit)
		lim = int64(limit)
	}

	characters, err := s.repo.ListCharacters(skip, lim)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.CountCharacters()
	if err != nil {
		return nil, 0, err
	}

	return characters, count, nil
}

func (s *service) SearchCharacters(name string) ([]models.Character, error) {
	if name == "" {
		return nil, errors.New("name query parameter is required")
	}

	// Implementasi pencarian langsung di repository bisa lebih efisien
	// Namun untuk contoh ini, kita gunakan ListCharacters dan filter manual
	characters, err := s.repo.ListCharacters(0, 0)
	if err != nil {
		return nil, err
	}

	var filtered []models.Character
	for _, c := range characters {
		if containsIgnoreCase(c.Name, name) {
			filtered = append(filtered, c)
		}
	}

	if len(filtered) == 0 {
		return nil, errors.New("no characters found")
	}

	return filtered, nil
}

func containsIgnoreCase(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(str) > 0 && (containsIgnoreCase(str[1:], substr) || str[:len(substr)] == substr))
}
