package controller

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/linkeunid/go-api/internal/model"
	"github.com/linkeunid/go-api/internal/service"
	"github.com/linkeunid/go-api/pkg/middleware"
	"github.com/linkeunid/go-api/pkg/pagination"
	"github.com/linkeunid/go-api/pkg/response"
	"go.uber.org/zap"
)

// Animal handles animal requests
type Animal struct {
	logger  *zap.Logger
	service service.AnimalService
}

// NewAnimal creates a new Animal controller instance
func NewAnimal(logger *zap.Logger, service service.AnimalService) *Animal {
	return &Animal{
		logger:  logger,
		service: service,
	}
}

// RegisterRoutes registers all routes for the animal controller
func (a *Animal) RegisterRoutes(r chi.Router) {
	r.Route("/animals", func(r chi.Router) {
		r.Get("/", a.GetAnimals)
		r.Post("/", a.CreateAnimal)
		r.Get("/{animalID}", a.GetAnimal)
		r.Put("/{animalID}", a.UpdateAnimal)
		r.Delete("/{animalID}", a.DeleteAnimal)
	})
}

// GetAnimals returns all animals
// @Summary Get all animals
// @Description Get a paginated list of all animals
// @Tags animals
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param sort query string false "Sort field (id, name, species, age, created_at, updated_at)"
// @Param direction query string false "Sort direction (asc, desc)"
// @Success 200 {object} response.APIResponse{data=pagination.PagedData{items=[]model.Animal}}
// @Failure 500 {object} response.APIResponse
// @Router /animals [get]
func (a *Animal) GetAnimals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse pagination parameters from the request
	params := pagination.NewParams(r)

	// Extract sort parameters and add to context
	queryParams := make(map[string]string)
	queryParams["sort"] = r.URL.Query().Get("sort")
	queryParams["direction"] = r.URL.Query().Get("direction")

	// Create a new context with query parameters
	ctxWithParams := context.WithValue(ctx, "queryParams", queryParams)

	// Get paginated animals
	animals, params, err := a.service.GetAllPaginated(ctxWithParams, params)
	if err != nil {
		a.logger.Error("Failed to get animals", zap.Error(err))
		response.InternalServerError(w, r, err)
		return
	}

	response.Paginated(w, r, animals, params, "Animals retrieved successfully")
}

// GetAnimal returns a specific animal by ID
// @Summary Get an animal by ID
// @Description Get an animal by its ID
// @Tags animals
// @Accept json
// @Produce json
// @Param animalID path string true "Animal ID"
// @Success 200 {object} response.APIResponse{data=model.Animal}
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /animals/{animalID} [get]
func (a *Animal) GetAnimal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	animalID := chi.URLParam(r, "animalID")

	animal, err := a.service.GetByID(ctx, animalID)
	if err != nil {
		if err == service.ErrAnimalNotFound {
			response.NotFound(w, r, "Animal not found")
			return
		}
		a.logger.Error("Failed to get animal", zap.String("id", animalID), zap.Error(err))
		response.InternalServerError(w, r, err)
		return
	}

	response.Success(w, r, animal, "Animal retrieved successfully")
}

// CreateAnimal creates a new animal
// @Summary Create a new animal
// @Description Create a new animal with the provided details
// @Tags animals
// @Accept json
// @Produce json
// @Param animal body model.AnimalCreateRequest true "Animal object to be created"
// @Success 201 {object} response.APIResponse{data=model.Animal}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /animals [post]
func (a *Animal) CreateAnimal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var animal model.Animal

	// Validate and decode the request
	if !middleware.HandleValidateRequest(w, r, &animal) {
		return
	}

	if err := a.service.Create(ctx, &animal); err != nil {
		if err == service.ErrInvalidAnimalData {
			response.BadRequest(w, r, "Invalid animal data", err)
			return
		}
		a.logger.Error("Failed to create animal", zap.Error(err))
		response.InternalServerError(w, r, err)
		return
	}

	response.Created(w, r, animal, "Animal created successfully")
}

// UpdateAnimal updates an existing animal
// @Summary Update an animal
// @Description Update an existing animal by its ID
// @Tags animals
// @Accept json
// @Produce json
// @Param animalID path string true "Animal ID"
// @Param animal body model.AnimalUpdateRequest true "Updated animal object"
// @Success 200 {object} response.APIResponse{data=model.Animal}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /animals/{animalID} [put]
func (a *Animal) UpdateAnimal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	animalID := chi.URLParam(r, "animalID")

	var animal model.Animal

	// Validate and decode the request
	if !middleware.HandleValidateRequest(w, r, &animal) {
		return
	}

	if err := a.service.Update(ctx, animalID, &animal); err != nil {
		switch err {
		case service.ErrAnimalNotFound:
			response.NotFound(w, r, "Animal not found")
		case service.ErrInvalidAnimalData:
			response.BadRequest(w, r, "Invalid animal data", err)
		default:
			a.logger.Error("Failed to update animal", zap.String("id", animalID), zap.Error(err))
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.Success(w, r, animal, "Animal updated successfully")
}

// DeleteAnimal deletes an animal
// @Summary Delete an animal
// @Description Delete an animal by its ID
// @Tags animals
// @Accept json
// @Produce json
// @Param animalID path string true "Animal ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /animals/{animalID} [delete]
func (a *Animal) DeleteAnimal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	animalID := chi.URLParam(r, "animalID")

	if err := a.service.Delete(ctx, animalID); err != nil {
		switch err {
		case service.ErrAnimalNotFound:
			response.NotFound(w, r, "Animal not found")
		case service.ErrInvalidAnimalData:
			response.BadRequest(w, r, "Invalid animal ID", err)
		default:
			a.logger.Error("Failed to delete animal", zap.String("id", animalID), zap.Error(err))
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.NoContent(w, r)
}
