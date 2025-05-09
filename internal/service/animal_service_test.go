package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/linkeunid/go-api/internal/model"
	"github.com/linkeunid/go-api/internal/repository"
	"github.com/linkeunid/go-api/pkg/config"
	"github.com/linkeunid/go-api/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockAnimalRepository is a mock implementation of the repository.AnimalRepository interface
type MockAnimalRepository struct {
	mock.Mock
}

func (m *MockAnimalRepository) FindAll(ctx context.Context) (repository.AnimalCollectionResult, error) {
	args := m.Called(ctx)
	return args.Get(0).(repository.AnimalCollectionResult), args.Error(1)
}

func (m *MockAnimalRepository) FindAllPaginated(ctx context.Context, params pagination.Params) (repository.AnimalCollectionResult, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(repository.AnimalCollectionResult), args.Error(1)
}

func (m *MockAnimalRepository) FindByID(ctx context.Context, id uint64) (repository.AnimalResult, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repository.AnimalResult), args.Error(1)
}

func (m *MockAnimalRepository) Create(ctx context.Context, animal *model.Animal) error {
	args := m.Called(ctx, animal)
	return args.Error(0)
}

func (m *MockAnimalRepository) Update(ctx context.Context, animal *model.Animal) error {
	args := m.Called(ctx, animal)
	return args.Error(0)
}

func (m *MockAnimalRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestAnimalServiceImpl_GetAll(t *testing.T) {
	// Create test logger
	logger, _ := zap.NewDevelopment()

	// Create test config
	cfg := &config.Config{}

	// Create test data
	animals := []model.Animal{
		{
			ID:          1,
			Name:        "Fluffy",
			Species:     "Cat",
			Age:         3,
			Description: "A fluffy cat",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "Rex",
			Species:     "Dog",
			Age:         5,
			Description: "A loyal dog",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Define test cases
	tests := []struct {
		name             string
		mockSetup        func(mockRepo *MockAnimalRepository)
		expectedResponse AnimalCollectionResponse
		expectedError    error
	}{
		{
			name: "Success",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				mockRepo.On("FindAll", mock.Anything).Return(repository.AnimalCollectionResult{
					Data: animals,
					CacheInfo: &repository.CacheInfo{
						Status:  "miss",
						Enabled: true,
					},
				}, nil)
			},
			expectedResponse: AnimalCollectionResponse{
				Data: animals,
				CacheInfo: &repository.CacheInfo{
					Status:  "miss",
					Enabled: true,
				},
			},
			expectedError: nil,
		},
		{
			name: "RepositoryError",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				mockRepo.On("FindAll", mock.Anything).Return(repository.AnimalCollectionResult{}, errors.New("database error"))
			},
			expectedResponse: AnimalCollectionResponse{},
			expectedError:    errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := new(MockAnimalRepository)

			// Setup the mock expectations
			tt.mockSetup(mockRepo)

			// Create service with mock repository
			service := NewAnimalService(cfg, logger, mockRepo)

			// Call the method being tested
			result, err := service.GetAll(context.Background())

			// Assert the error
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Assert the response
			assert.Equal(t, len(tt.expectedResponse.Data), len(result.Data))

			// Verify that repository method was called
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAnimalServiceImpl_GetAllPaginated(t *testing.T) {
	// Create test logger
	logger, _ := zap.NewDevelopment()

	// Create test config
	cfg := &config.Config{}

	// Create test data
	animals := []model.Animal{
		{
			ID:          1,
			Name:        "Fluffy",
			Species:     "Cat",
			Age:         3,
			Description: "A fluffy cat",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Define pagination params
	params := pagination.Params{
		Page:       1,
		Limit:      10,
		TotalItems: 1,
		TotalPages: 1,
	}

	// Define test cases
	tests := []struct {
		name             string
		mockSetup        func(mockRepo *MockAnimalRepository)
		expectedResponse AnimalCollectionResponse
		expectedError    error
	}{
		{
			name: "Success",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				mockRepo.On("FindAllPaginated", mock.Anything, params).Return(repository.AnimalCollectionResult{
					Data: animals,
					Pagination: &pagination.Params{
						Page:       1,
						Limit:      10,
						TotalItems: 1,
						TotalPages: 1,
					},
					CacheInfo: &repository.CacheInfo{
						Status:  "miss",
						Enabled: true,
					},
				}, nil)
			},
			expectedResponse: AnimalCollectionResponse{
				Data: animals,
				Pagination: &pagination.Params{
					Page:       1,
					Limit:      10,
					TotalItems: 1,
					TotalPages: 1,
				},
				CacheInfo: &repository.CacheInfo{
					Status:  "miss",
					Enabled: true,
				},
			},
			expectedError: nil,
		},
		{
			name: "RepositoryError",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				mockRepo.On("FindAllPaginated", mock.Anything, params).Return(repository.AnimalCollectionResult{}, errors.New("database error"))
			},
			expectedResponse: AnimalCollectionResponse{},
			expectedError:    errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := new(MockAnimalRepository)

			// Setup the mock expectations
			tt.mockSetup(mockRepo)

			// Create service with mock repository
			service := NewAnimalService(cfg, logger, mockRepo)

			// Call the method being tested
			result, err := service.GetAllPaginated(context.Background(), params)

			// Assert the error
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Assert the response
			assert.Equal(t, len(tt.expectedResponse.Data), len(result.Data))

			// Verify that repository method was called
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAnimalServiceImpl_GetByID(t *testing.T) {
	// Create test logger
	logger, _ := zap.NewDevelopment()

	// Create test config
	cfg := &config.Config{}

	// Create test data
	animal := model.Animal{
		ID:          1,
		Name:        "Fluffy",
		Species:     "Cat",
		Age:         3,
		Description: "A fluffy cat",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Define test cases
	tests := []struct {
		name             string
		animalID         string
		mockSetup        func(mockRepo *MockAnimalRepository)
		expectedResponse AnimalResponse
		expectedError    error
	}{
		{
			name:     "Success",
			animalID: "1",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(repository.AnimalResult{
					Data: &animal,
					CacheInfo: &repository.CacheInfo{
						Status:  "miss",
						Enabled: true,
					},
				}, nil)
			},
			expectedResponse: AnimalResponse{
				Data: &animal,
				CacheInfo: &repository.CacheInfo{
					Status:  "miss",
					Enabled: true,
				},
			},
			expectedError: nil,
		},
		{
			name:     "NotFound",
			animalID: "999",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				mockRepo.On("FindByID", mock.Anything, uint64(999)).Return(repository.AnimalResult{
					Data:      nil,
					CacheInfo: nil,
				}, nil)
			},
			expectedResponse: AnimalResponse{},
			expectedError:    ErrAnimalNotFound,
		},
		{
			name:     "InvalidID",
			animalID: "invalid",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for invalid ID
			},
			expectedResponse: AnimalResponse{},
			expectedError:    ErrInvalidAnimalID,
		},
		{
			name:     "EmptyID",
			animalID: "",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for empty ID
			},
			expectedResponse: AnimalResponse{},
			expectedError:    ErrInvalidAnimalData,
		},
		{
			name:     "RepositoryError",
			animalID: "1",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(repository.AnimalResult{}, errors.New("database error"))
			},
			expectedResponse: AnimalResponse{},
			expectedError:    errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := new(MockAnimalRepository)

			// Setup the mock expectations
			tt.mockSetup(mockRepo)

			// Create service with mock repository
			service := NewAnimalService(cfg, logger, mockRepo)

			// Call the method being tested
			result, err := service.GetByID(context.Background(), tt.animalID)

			// Assert the error
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError == ErrAnimalNotFound || tt.expectedError == ErrInvalidAnimalData || tt.expectedError == ErrInvalidAnimalID {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse.Data.ID, result.Data.ID)
				assert.Equal(t, tt.expectedResponse.Data.Name, result.Data.Name)
			}

			// Verify that repository method was called
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAnimalServiceImpl_Create(t *testing.T) {
	// Create test logger
	logger, _ := zap.NewDevelopment()

	// Create test config
	cfg := &config.Config{}

	// Define test cases
	tests := []struct {
		name          string
		animal        *model.Animal
		mockSetup     func(mockRepo *MockAnimalRepository)
		expectedError error
	}{
		{
			name: "Success",
			animal: &model.Animal{
				Name:        "Fluffy",
				Species:     "Cat",
				Age:         3,
				Description: "A fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Animal")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "NilAnimal",
			animal: nil,
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for nil animal
			},
			expectedError: ErrInvalidAnimalData,
		},
		{
			name: "EmptyName",
			animal: &model.Animal{
				Name:        "",
				Species:     "Cat",
				Age:         3,
				Description: "A fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for invalid animal
			},
			expectedError: ErrInvalidAnimalData,
		},
		{
			name: "EmptySpecies",
			animal: &model.Animal{
				Name:        "Fluffy",
				Species:     "",
				Age:         3,
				Description: "A fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for invalid animal
			},
			expectedError: ErrInvalidAnimalData,
		},
		{
			name: "RepositoryError",
			animal: &model.Animal{
				Name:        "Fluffy",
				Species:     "Cat",
				Age:         3,
				Description: "A fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Animal")).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := new(MockAnimalRepository)

			// Setup the mock expectations
			tt.mockSetup(mockRepo)

			// Create service with mock repository
			service := NewAnimalService(cfg, logger, mockRepo)

			// Call the method being tested
			err := service.Create(context.Background(), tt.animal)

			// Assert the error
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError == ErrInvalidAnimalData {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify that repository method was called
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAnimalServiceImpl_Update(t *testing.T) {
	// Create test logger
	logger, _ := zap.NewDevelopment()

	// Create test config
	cfg := &config.Config{}

	// Create test data
	existingAnimal := model.Animal{
		ID:          1,
		Name:        "Fluffy",
		Species:     "Cat",
		Age:         3,
		Description: "A fluffy cat",
		CreatedAt:   time.Now().Add(-24 * time.Hour), // Created yesterday
		UpdatedAt:   time.Now().Add(-12 * time.Hour), // Updated 12 hours ago
	}

	// Define test cases
	tests := []struct {
		name          string
		animalID      string
		animal        *model.Animal
		mockSetup     func(mockRepo *MockAnimalRepository)
		expectedError error
	}{
		{
			name:     "Success",
			animalID: "1",
			animal: &model.Animal{
				Name:        "Fluffy Updated",
				Species:     "Cat",
				Age:         4,
				Description: "An updated fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// First call to FindByID to check if the animal exists
				mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(repository.AnimalResult{
					Data: &existingAnimal,
				}, nil)

				// Second call to Update to update the animal
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Animal")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "NotFound",
			animalID: "999",
			animal: &model.Animal{
				Name:        "Fluffy Updated",
				Species:     "Cat",
				Age:         4,
				Description: "An updated fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// Animal not found
				mockRepo.On("FindByID", mock.Anything, uint64(999)).Return(repository.AnimalResult{
					Data: nil,
				}, nil)
			},
			expectedError: ErrAnimalNotFound,
		},
		{
			name:     "EmptyID",
			animalID: "",
			animal: &model.Animal{
				Name:        "Fluffy Updated",
				Species:     "Cat",
				Age:         4,
				Description: "An updated fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for empty ID
			},
			expectedError: ErrInvalidAnimalData,
		},
		{
			name:     "InvalidID",
			animalID: "invalid",
			animal: &model.Animal{
				Name:        "Fluffy Updated",
				Species:     "Cat",
				Age:         4,
				Description: "An updated fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for invalid ID
			},
			expectedError: ErrInvalidAnimalID,
		},
		{
			name:     "NilAnimal",
			animalID: "1",
			animal:   nil,
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for nil animal
			},
			expectedError: ErrInvalidAnimalData,
		},
		{
			name:     "EmptyName",
			animalID: "1",
			animal: &model.Animal{
				Name:        "",
				Species:     "Cat",
				Age:         4,
				Description: "An updated fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for invalid animal
			},
			expectedError: ErrInvalidAnimalData,
		},
		{
			name:     "EmptySpecies",
			animalID: "1",
			animal: &model.Animal{
				Name:        "Fluffy Updated",
				Species:     "",
				Age:         4,
				Description: "An updated fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for invalid animal
			},
			expectedError: ErrInvalidAnimalData,
		},
		{
			name:     "FindByIDError",
			animalID: "1",
			animal: &model.Animal{
				Name:        "Fluffy Updated",
				Species:     "Cat",
				Age:         4,
				Description: "An updated fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// Error during FindByID
				mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(repository.AnimalResult{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:     "UpdateError",
			animalID: "1",
			animal: &model.Animal{
				Name:        "Fluffy Updated",
				Species:     "Cat",
				Age:         4,
				Description: "An updated fluffy cat",
			},
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// First call to FindByID succeeds
				mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(repository.AnimalResult{
					Data: &existingAnimal,
				}, nil)

				// Second call to Update fails
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Animal")).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := new(MockAnimalRepository)

			// Setup the mock expectations
			tt.mockSetup(mockRepo)

			// Create service with mock repository
			service := NewAnimalService(cfg, logger, mockRepo)

			// Call the method being tested
			err := service.Update(context.Background(), tt.animalID, tt.animal)

			// Assert the error
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError == ErrAnimalNotFound || tt.expectedError == ErrInvalidAnimalData || tt.expectedError == ErrInvalidAnimalID {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify that repository method was called
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAnimalServiceImpl_Delete(t *testing.T) {
	// Create test logger
	logger, _ := zap.NewDevelopment()

	// Create test config
	cfg := &config.Config{}

	// Create test data
	existingAnimal := model.Animal{
		ID:          1,
		Name:        "Fluffy",
		Species:     "Cat",
		Age:         3,
		Description: "A fluffy cat",
	}

	// Define test cases
	tests := []struct {
		name          string
		animalID      string
		mockSetup     func(mockRepo *MockAnimalRepository)
		expectedError error
	}{
		{
			name:     "Success",
			animalID: "1",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// First call to FindByID to check if the animal exists
				mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(repository.AnimalResult{
					Data: &existingAnimal,
				}, nil)

				// Second call to Delete to delete the animal
				mockRepo.On("Delete", mock.Anything, uint64(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "NotFound",
			animalID: "999",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// Animal not found
				mockRepo.On("FindByID", mock.Anything, uint64(999)).Return(repository.AnimalResult{
					Data: nil,
				}, nil)
			},
			expectedError: ErrAnimalNotFound,
		},
		{
			name:     "EmptyID",
			animalID: "",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for empty ID
			},
			expectedError: ErrInvalidAnimalData,
		},
		{
			name:     "InvalidID",
			animalID: "invalid",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// No repository call expected for invalid ID
			},
			expectedError: ErrInvalidAnimalID,
		},
		{
			name:     "FindByIDError",
			animalID: "1",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// Error during FindByID
				mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(repository.AnimalResult{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:     "DeleteError",
			animalID: "1",
			mockSetup: func(mockRepo *MockAnimalRepository) {
				// First call to FindByID succeeds
				mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(repository.AnimalResult{
					Data: &existingAnimal,
				}, nil)

				// Second call to Delete fails
				mockRepo.On("Delete", mock.Anything, uint64(1)).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := new(MockAnimalRepository)

			// Setup the mock expectations
			tt.mockSetup(mockRepo)

			// Create service with mock repository
			service := NewAnimalService(cfg, logger, mockRepo)

			// Call the method being tested
			err := service.Delete(context.Background(), tt.animalID)

			// Assert the error
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError == ErrAnimalNotFound || tt.expectedError == ErrInvalidAnimalData || tt.expectedError == ErrInvalidAnimalID {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify that repository method was called
			mockRepo.AssertExpectations(t)
		})
	}
}
