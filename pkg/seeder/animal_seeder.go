package seeder

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/linkeunid/go-api/internal/model"
	"github.com/linkeunid/go-api/pkg/database"
	"go.uber.org/zap"
)

// AnimalSeeder seeds animal data
type AnimalSeeder struct {
	db     database.Database
	logger *zap.Logger
	count  int
}

// NewAnimalSeeder creates a new animal seeder
func NewAnimalSeeder(db database.Database, logger *zap.Logger, count int) *AnimalSeeder {
	return &AnimalSeeder{
		db:     db,
		logger: logger,
		count:  count,
	}
}

// GetName returns the name of the seeder
func (s *AnimalSeeder) GetName() string {
	return "animal"
}

// FakerAnimal is a struct used to generate animal data with faker tags
type FakerAnimal struct {
	Name        string `faker:"pet_name"`
	Species     string `faker:"animal_species"`
	Age         int    `faker:"boundary_start=1,boundary_end=15"`
	Description string `faker:"animal_description"`
}

// Custom faker providers
func init() {
	// Create a random source using the current time
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Register pet name provider
	_ = faker.AddProvider("pet_name", func(v reflect.Value) (interface{}, error) {
		petNames := []string{
			"Bella", "Max", "Luna", "Charlie", "Lucy", "Cooper", "Daisy", "Milo",
			"Sadie", "Rocky", "Molly", "Buddy", "Bailey", "Maggie", "Jack",
			"Lola", "Oliver", "Stella", "Zeus", "Lily", "Duke", "Zoe", "Bentley",
			"Sophie", "Toby", "Chloe", "Dexter", "Penny", "Gus", "Willow",
		}
		return petNames[r.Intn(len(petNames))], nil
	})

	// Register animal species provider
	_ = faker.AddProvider("animal_species", func(v reflect.Value) (interface{}, error) {
		species := []string{
			"Dog", "Cat", "Rabbit", "Hamster", "Guinea Pig", "Parrot", "Goldfish",
			"Turtle", "Snake", "Lizard", "Horse", "Cow", "Pig", "Sheep", "Goat",
			"Chicken", "Duck", "Donkey", "Ferret", "Chinchilla",
		}
		return species[r.Intn(len(species))], nil
	})

	// Register animal description provider
	_ = faker.AddProvider("animal_description", func(v reflect.Value) (interface{}, error) {
		descriptions := []string{
			"Very friendly and playful",
			"A bit shy but very loving",
			"Energetic and loves to run",
			"Calm and well-behaved",
			"Curious and intelligent",
			"Loves cuddles and attention",
			"Independent but affectionate",
			"Protective and loyal",
			"Gentle with children",
			"Loves to play with toys",
			"Very social with other animals",
			"Quiet and observant",
		}
		return descriptions[r.Intn(len(descriptions))], nil
	})
}

// Seed seeds animal data
func (s *AnimalSeeder) Seed(ctx context.Context) error {
	// Check if there are already animals in the database
	var count int64
	if err := s.db.GetDB().Model(&model.Animal{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count animals: %w", err)
	}

	// Skip seeding if animals already exist
	if count > 0 {
		s.logger.Info("Animals already exist, skipping seeding", zap.Int64("count", count))
		return nil
	}

	// Create animals
	animals, err := s.generateAnimals(s.count)
	if err != nil {
		return fmt.Errorf("failed to generate animals: %w", err)
	}

	s.logger.Info("Seeding animals", zap.Int("count", len(animals)))

	// Start a transaction for better data consistency
	tx := s.db.GetDB().Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Insert in batches for better performance
	batchSize := 100
	for i := 0; i < len(animals); i += batchSize {
		end := i + batchSize
		if end > len(animals) {
			end = len(animals)
		}

		batch := animals[i:end]
		if err := tx.CreateInBatches(batch, len(batch)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to seed animals batch %d: %w", i/batchSize, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("Successfully seeded animals", zap.Int("count", len(animals)))
	return nil
}

// generateAnimals creates a slice of random animal data using faker
func (s *AnimalSeeder) generateAnimals(count int) ([]*model.Animal, error) {
	// Create a random source
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Configure faker
	faker.SetRandomStringLength(10)

	// Generate random animals
	animals := make([]*model.Animal, count)
	for i := 0; i < count; i++ {
		// Generate fake animal data using tags
		fakeAnimal := FakerAnimal{}
		if err := faker.FakeData(&fakeAnimal); err != nil {
			return nil, fmt.Errorf("failed to generate fake animal data: %w", err)
		}

		// Generate a creation time within the last year
		createdAt := time.Now().Add(-time.Duration(r.Intn(365)) * 24 * time.Hour)

		// Generate an update time between creation and now
		updateDuration := time.Since(createdAt)
		updateOffset := time.Duration(r.Int63n(int64(updateDuration)))
		updatedAt := createdAt.Add(updateOffset)

		// Create animal model from fake data
		animal := &model.Animal{
			Name:        fakeAnimal.Name,
			Species:     fakeAnimal.Species,
			Age:         fakeAnimal.Age,
			Description: fakeAnimal.Description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		animals[i] = animal
	}

	return animals, nil
}
