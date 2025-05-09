package model

import (
	"fmt"
	"time"

	"github.com/linkeunid/go-api/pkg/validator"
)

// Animal represents an animal entity
type Animal struct {
	ID          uint64    `json:"id" gorm:"primaryKey;type:bigint unsigned;autoIncrement"`
	Name        string    `json:"name" validate:"required,min=2,max=100,animalname" gorm:"type:varchar(100);not null;index:idx_animal_name" example:"Fluffy"`
	Species     string    `json:"species" validate:"required,min=2,max=100" gorm:"type:varchar(100);not null;index:idx_animal_species" example:"Cat"`
	Age         int       `json:"age" validate:"gte=0,lte=200" gorm:"type:int;index:idx_animal_age" example:"3"`
	Description string    `json:"description" validate:"omitempty,max=1000" gorm:"type:text" example:"A friendly cat with white fur"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;index:idx_animal_created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;index:idx_animal_updated_at"`
}

// AnimalCreateRequest represents a request body example for creating a new animal
// @name AnimalCreateRequest
type AnimalCreateRequest struct {
	Name        string `json:"name" example:"Fluffy"`
	Species     string `json:"species" example:"Cat"`
	Age         int    `json:"age" example:"3"`
	Description string `json:"description" example:"A friendly cat with white fur"`
}

// AnimalUpdateRequest represents a request body example for updating an animal
// @name AnimalUpdateRequest
type AnimalUpdateRequest struct {
	Name        string `json:"name" example:"Fluffy"`
	Species     string `json:"species" example:"Cat"`
	Age         int    `json:"age" example:"3"`
	Description string `json:"description" example:"A friendly cat with white fur"`
}

// TableName returns the table name for the Animal model
func (Animal) TableName() string {
	return "animals"
}

// CacheEnabled returns whether this model should be cached
func (Animal) CacheEnabled() bool {
	return true
}

// CacheTTL returns the time-to-live for this model in cache
func (Animal) CacheTTL() time.Duration {
	return 30 * time.Minute
}

// CacheKey returns a unique key for this model instance
func (a Animal) CacheKey() string {
	return fmt.Sprintf("animal:%d", a.ID)
}

// Validate performs validation on the Animal model
func (a Animal) Validate() []validator.ValidationError {
	return validator.Validate(a)
}
