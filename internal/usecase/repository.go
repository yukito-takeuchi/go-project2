package usecase

import (
	"context"

	"Aicon-assignment/internal/domain/entity"
)

// ItemRepository defines the interface for item data access
type ItemRepository interface {
	// FindAll retrieves all items
	FindAll(ctx context.Context) ([]*entity.Item, error)

	// FindByID retrieves an item by ID
	FindByID(ctx context.Context, id int64) (*entity.Item, error)

	// Create creates a new item and returns it with the generated ID
	Create(ctx context.Context, item *entity.Item) (*entity.Item, error)

	// Delete deletes an item by ID
	Delete(ctx context.Context, id int64) error

	// GetSummaryByCategory returns item counts grouped by category (bonus feature)
	GetSummaryByCategory(ctx context.Context) (map[string]int, error)
}
