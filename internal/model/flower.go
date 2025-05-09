package model

import (
	"fmt"
	"time"

	"github.com/linkeunid/go-api/pkg/validator"
)

// Flower represents a flower entity
type Flower struct {
	ID          uint64    `json:"id" gorm:"primaryKey;type:bigint unsigned;autoIncrement"`
	Name        string    `json:"name" validate:"required,min=2,max=100" gorm:"type:varchar(100);not null;index:idx_flower_name" example:"Rose"`
	Species     string    `json:"species" validate:"required,min=2,max=100" gorm:"type:varchar(100);not null;index:idx_flower_species" example:"Rosa"`
	Color       string    `json:"color" validate:"required,min=2,max=50" gorm:"type:varchar(50);not null;index:idx_flower_color" example:"Red"`
	Description string    `json:"description" validate:"omitempty,max=1000" gorm:"type:text" example:"A beautiful red rose with thorny stems"`
	Seasonal    bool      `json:"seasonal" gorm:"type:boolean;default:false" example:"true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;index:idx_flower_created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;index:idx_flower_updated_at"`
}

// FlowerCreateRequest represents a request body example for creating a new flower
// @name FlowerCreateRequest
type FlowerCreateRequest struct {
	Name        string `json:"name" example:"Rose"`
	Species     string `json:"species" example:"Rosa"`
	Color       string `json:"color" example:"Red"`
	Description string `json:"description" example:"A beautiful red rose with thorny stems"`
	Seasonal    bool   `json:"seasonal" example:"true"`
}

// FlowerUpdateRequest represents a request body example for updating a flower
// @name FlowerUpdateRequest
type FlowerUpdateRequest struct {
	Name        string `json:"name" example:"Rose"`
	Species     string `json:"species" example:"Rosa"`
	Color       string `json:"color" example:"Red"`
	Description string `json:"description" example:"A beautiful red rose with thorny stems"`
	Seasonal    bool   `json:"seasonal" example:"true"`
}

// TableName returns the table name for the Flower model
func (Flower) TableName() string {
	return "flowers"
}

// CacheEnabled returns whether this model should be cached
func (Flower) CacheEnabled() bool {
	return true
}

// CacheTTL returns the time-to-live for this model in cache
func (Flower) CacheTTL() time.Duration {
	return 30 * time.Minute
}

// CacheKey returns a unique key for this model instance
func (f Flower) CacheKey() string {
	return fmt.Sprintf("flower:%d", f.ID)
}

// Validate performs validation on the Flower model
func (f Flower) Validate() []validator.ValidationError {
	return validator.Validate(f)
}
