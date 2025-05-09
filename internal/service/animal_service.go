package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/linkeunid/go-api/internal/model"
	"github.com/linkeunid/go-api/internal/repository"
	"github.com/linkeunid/go-api/pkg/config"
	"github.com/linkeunid/go-api/pkg/pagination"
	"go.uber.org/zap"
)

var (
	// ErrAnimalNotFound is returned when an animal cannot be found
	ErrAnimalNotFound = errors.New("animal not found")

	// ErrInvalidAnimalData is returned when animal data is invalid
	ErrInvalidAnimalData = errors.New("invalid animal data")

	// ErrInvalidAnimalID is returned when animal ID is invalid
	ErrInvalidAnimalID = errors.New("invalid animal ID")
)

// AnimalService defines the interface for animal operations
type AnimalService interface {
	GetAll(ctx context.Context) ([]model.Animal, error)
	GetAllPaginated(ctx context.Context, params pagination.Params) ([]model.Animal, pagination.Params, error)
	GetByID(ctx context.Context, id string) (*model.Animal, error)
	Create(ctx context.Context, animal *model.Animal) error
	Update(ctx context.Context, id string, animal *model.Animal) error
	Delete(ctx context.Context, id string) error
}

// AnimalServiceImpl implements AnimalService
type AnimalServiceImpl struct {
	logger     *zap.Logger
	config     *config.Config
	repository repository.AnimalRepository
}

// NewAnimalService creates a new animal service
func NewAnimalService(
	cfg *config.Config,
	logger *zap.Logger,
	repository repository.AnimalRepository,
) AnimalService {
	return &AnimalServiceImpl{
		logger:     logger,
		config:     cfg,
		repository: repository,
	}
}

// GetAll retrieves all animals
func (s *AnimalServiceImpl) GetAll(ctx context.Context) ([]model.Animal, error) {
	// Add a timeout to the context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repository.FindAll(ctx)
}

// GetAllPaginated retrieves paginated animals
func (s *AnimalServiceImpl) GetAllPaginated(ctx context.Context, params pagination.Params) ([]model.Animal, pagination.Params, error) {
	// Add a timeout to the context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repository.FindAllPaginated(ctx, params)
}

// GetByID retrieves an animal by ID
func (s *AnimalServiceImpl) GetByID(ctx context.Context, id string) (*model.Animal, error) {
	if id == "" {
		return nil, ErrInvalidAnimalData
	}

	// Convert string ID to uint64
	numericID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		s.logger.Error("Invalid animal ID format", zap.String("id", id), zap.Error(err))
		return nil, ErrInvalidAnimalID
	}

	// Add a timeout to the context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	animal, err := s.repository.FindByID(ctx, numericID)
	if err != nil {
		return nil, err
	}

	if animal == nil {
		return nil, ErrAnimalNotFound
	}

	return animal, nil
}

// Create creates a new animal
func (s *AnimalServiceImpl) Create(ctx context.Context, animal *model.Animal) error {
	if animal == nil || animal.Name == "" || animal.Species == "" {
		return ErrInvalidAnimalData
	}

	// Add a timeout to the context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repository.Create(ctx, animal)
}

// Update updates an existing animal
func (s *AnimalServiceImpl) Update(ctx context.Context, id string, animal *model.Animal) error {
	if id == "" || animal == nil || animal.Name == "" || animal.Species == "" {
		return ErrInvalidAnimalData
	}

	// Convert string ID to uint64
	numericID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		s.logger.Error("Invalid animal ID format", zap.String("id", id), zap.Error(err))
		return ErrInvalidAnimalID
	}

	// Ensure the ID in the path matches the animal ID
	animal.ID = numericID

	// Add a timeout to the context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if the animal exists
	existing, err := s.repository.FindByID(ctx, numericID)
	if err != nil {
		return err
	}

	if existing == nil {
		return ErrAnimalNotFound
	}

	// Preserve created_at timestamp
	animal.CreatedAt = existing.CreatedAt

	return s.repository.Update(ctx, animal)
}

// Delete removes an animal
func (s *AnimalServiceImpl) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidAnimalData
	}

	// Convert string ID to uint64
	numericID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		s.logger.Error("Invalid animal ID format", zap.String("id", id), zap.Error(err))
		return ErrInvalidAnimalID
	}

	// Add a timeout to the context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if the animal exists
	existing, err := s.repository.FindByID(ctx, numericID)
	if err != nil {
		return err
	}

	if existing == nil {
		return ErrAnimalNotFound
	}

	return s.repository.Delete(ctx, numericID)
}
