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

// FlowerSeeder seeds flower data
type FlowerSeeder struct {
	db     database.Database
	logger *zap.Logger
	count  int
}

// NewFlowerSeeder creates a new flower seeder
func NewFlowerSeeder(db database.Database, logger *zap.Logger, count int) *FlowerSeeder {
	return &FlowerSeeder{
		db:     db,
		logger: logger,
		count:  count,
	}
}

// GetName returns the name of the seeder
func (s *FlowerSeeder) GetName() string {
	return "flower"
}

// FakerFlower is a struct used to generate flower data with faker tags
type FakerFlower struct {
	Name        string `faker:"flower_name"`
	Species     string `faker:"flower_species"`
	Color       string `faker:"flower_color"`
	Description string `faker:"flower_description"`
	Seasonal    bool   `faker:"-"`
}

// Custom faker providers
func init() {
	// Create a random source using the current time
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Register flower name provider
	_ = faker.AddProvider("flower_name", func(v reflect.Value) (interface{}, error) {
		flowerNames := []string{
			"Rose", "Tulip", "Daisy", "Sunflower", "Lily", "Orchid", "Daffodil",
			"Carnation", "Peony", "Iris", "Chrysanthemum", "Poppy", "Marigold",
			"Hibiscus", "Magnolia", "Lavender", "Dahlia", "Hydrangea", "Jasmine",
			"Bluebell", "Cherry Blossom", "Buttercup", "Forget-me-not", "Dandelion",
		}
		return flowerNames[r.Intn(len(flowerNames))], nil
	})

	// Register flower species provider
	_ = faker.AddProvider("flower_species", func(v reflect.Value) (interface{}, error) {
		species := []string{
			"Rosa", "Tulipa", "Bellis", "Helianthus", "Lilium", "Orchidaceae",
			"Narcissus", "Dianthus", "Paeonia", "Iridaceae", "Chrysanthemum",
			"Papaver", "Tagetes", "Hibiscus", "Magnolia", "Lavandula", "Dahlia",
			"Hydrangea", "Jasminum", "Hyacinthoides", "Prunus", "Ranunculus",
			"Myosotis", "Taraxacum",
		}
		return species[r.Intn(len(species))], nil
	})

	// Register flower color provider
	_ = faker.AddProvider("flower_color", func(v reflect.Value) (interface{}, error) {
		colors := []string{
			"Red", "Pink", "Yellow", "Orange", "Purple", "Blue", "White",
			"Violet", "Indigo", "Cream", "Coral", "Lavender", "Maroon",
			"Fuchsia", "Peach", "Magenta", "Crimson", "Lilac", "Gold", "Burgundy",
		}
		return colors[r.Intn(len(colors))], nil
	})

	// Register flower description provider
	_ = faker.AddProvider("flower_description", func(v reflect.Value) (interface{}, error) {
		descriptions := []string{
			"Beautiful fragrant flower with soft petals",
			"Bold and vibrant with striking colors",
			"Delicate flower with a sweet fragrance",
			"Hardy perennial with long-lasting blooms",
			"Exotic flower with unique features",
			"Perfect for garden borders and beds",
			"Elegant flower that attracts butterflies",
			"Drought-resistant variety with minimal care needs",
			"Showy blooms that make excellent cut flowers",
			"Spreads rapidly with abundant flowers",
			"Rare variety with spectacular blooms",
			"Native wildflower with ecological benefits",
		}
		return descriptions[r.Intn(len(descriptions))], nil
	})
}

// Seed seeds flower data
func (s *FlowerSeeder) Seed(ctx context.Context) error {
	// Check if there are already flowers in the database
	var count int64
	if err := s.db.GetDB().Model(&model.Flower{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count flowers: %w", err)
	}

	// Skip seeding if flowers already exist
	if count > 0 {
		s.logger.Info("Flowers already exist, skipping seeding", zap.Int64("count", count))
		return nil
	}

	// Create flowers
	flowers, err := s.generateFlowers(s.count)
	if err != nil {
		return fmt.Errorf("failed to generate flowers: %w", err)
	}

	s.logger.Info("Seeding flowers", zap.Int("count", len(flowers)))

	// Start a transaction for better data consistency
	tx := s.db.GetDB().Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Insert in batches for better performance
	batchSize := 100
	for i := 0; i < len(flowers); i += batchSize {
		end := i + batchSize
		if end > len(flowers) {
			end = len(flowers)
		}

		batch := flowers[i:end]
		if err := tx.CreateInBatches(batch, len(batch)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to seed flowers batch %d: %w", i/batchSize, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("Successfully seeded flowers", zap.Int("count", len(flowers)))
	return nil
}

// generateFlowers creates a slice of random flower data using faker
func (s *FlowerSeeder) generateFlowers(count int) ([]*model.Flower, error) {
	// Create a random source
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Configure faker
	if err := faker.SetRandomStringLength(10); err != nil {
		return nil, fmt.Errorf("failed to set random string length: %w", err)
	}

	// Generate random flowers
	flowers := make([]*model.Flower, count)
	for i := 0; i < count; i++ {
		// Generate fake flower data using tags
		fakeFlower := FakerFlower{}
		if err := faker.FakeData(&fakeFlower); err != nil {
			return nil, fmt.Errorf("failed to generate fake flower data: %w", err)
		}

		// Generate a creation time within the last year
		createdAt := time.Now().Add(-time.Duration(r.Intn(365)) * 24 * time.Hour)

		// Generate an update time between creation and now
		updateDuration := time.Since(createdAt)
		updateOffset := time.Duration(r.Int63n(int64(updateDuration)))
		updatedAt := createdAt.Add(updateOffset)

		// Create flower model from fake data
		flower := &model.Flower{
			Name:        fakeFlower.Name,
			Species:     fakeFlower.Species,
			Color:       fakeFlower.Color,
			Description: fakeFlower.Description,
			Seasonal:    r.Intn(2) == 1,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		flowers[i] = flower
	}

	return flowers, nil
}
