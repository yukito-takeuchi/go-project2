package usecase

import (
	"context"
	"fmt"

	"Aicon-assignment/internal/domain/entity"
	domainErrors "Aicon-assignment/internal/domain/errors"
)

type ItemUsecase interface {
	GetAllItems(ctx context.Context) ([]*entity.Item, error)
	GetItemByID(ctx context.Context, id int64) (*entity.Item, error)
	CreateItem(ctx context.Context, input CreateItemInput) (*entity.Item, error)
	UpdateItem(ctx context.Context, id int64, input UpdateItemInput) (*entity.Item, error)
	DeleteItem(ctx context.Context, id int64) error
	GetCategorySummary(ctx context.Context) (*CategorySummary, error)
}

type CreateItemInput struct {
	Name          string `json:"name"`
	Category      string `json:"category"`
	Brand         string `json:"brand"`
	PurchasePrice int    `json:"purchase_price"`
	PurchaseDate  string `json:"purchase_date"`
}

type UpdateItemInput struct {
	Name          *string `json:"name"`
	Brand         *string `json:"brand"`
	PurchasePrice *int    `json:"purchase_price"`
}

type CategorySummary struct {
	Categories map[string]int `json:"categories"`
	Total      int            `json:"total"`
}

type itemUsecase struct {
	itemRepo ItemRepository
}

func NewItemUsecase(itemRepo ItemRepository) ItemUsecase {
	return &itemUsecase{
		itemRepo: itemRepo,
	}
}

func (u *itemUsecase) GetAllItems(ctx context.Context) ([]*entity.Item, error) {
	items, err := u.itemRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve items: %w", err)
	}

	return items, nil
}

func (u *itemUsecase) GetItemByID(ctx context.Context, id int64) (*entity.Item, error) {
	if id <= 0 {
		return nil, domainErrors.ErrInvalidInput
	}

	item, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		if domainErrors.IsNotFoundError(err) {
			return nil, domainErrors.ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to retrieve item: %w", err)
	}

	return item, nil
}

func (u *itemUsecase) CreateItem(ctx context.Context, input CreateItemInput) (*entity.Item, error) {
	// バリデーションして、新しいエンティティを作成
	item, err := entity.NewItem(
		input.Name,
		input.Category,
		input.Brand,
		input.PurchasePrice,
		input.PurchaseDate,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domainErrors.ErrInvalidInput, err.Error())
	}

	createdItem, err := u.itemRepo.Create(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return createdItem, nil
}

func (u *itemUsecase) UpdateItem(ctx context.Context, id int64, input UpdateItemInput) (*entity.Item, error) {
	if id <= 0 {
		return nil, domainErrors.ErrInvalidInput
	}

	// 既存アイテムを取得
	existingItem, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		if domainErrors.IsNotFoundError(err) {
			return nil, domainErrors.ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to retrieve item: %w", err)
	}

	// 指定されたフィールドのみ更新
	if input.Name != nil {
		existingItem.Name = *input.Name
	}
	if input.Brand != nil {
		existingItem.Brand = *input.Brand
	}
	if input.PurchasePrice != nil {
		existingItem.PurchasePrice = *input.PurchasePrice
	}

	// バリデーション
	if err := existingItem.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", domainErrors.ErrInvalidInput, err.Error())
	}

	// 更新を実行
	updatedItem, err := u.itemRepo.Update(ctx, existingItem)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return updatedItem, nil
}

func (u *itemUsecase) DeleteItem(ctx context.Context, id int64) error {
	if id <= 0 {
		return domainErrors.ErrInvalidInput
	}

	_, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		if domainErrors.IsNotFoundError(err) {
			return domainErrors.ErrItemNotFound
		}
		return fmt.Errorf("failed to check item existence: %w", err)
	}

	err = u.itemRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

func (u *itemUsecase) GetCategorySummary(ctx context.Context) (*CategorySummary, error) {
	categoryCounts, err := u.itemRepo.GetSummaryByCategory(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get category summary: %w", err)
	}

	// 合計計算
	total := 0
	for _, count := range categoryCounts {
		total += count
	}

	summary := make(map[string]int)
	for _, category := range entity.GetValidCategories() {
		if count, exists := categoryCounts[category]; exists {
			summary[category] = count
		} else {
			summary[category] = 0
		}
	}

	return &CategorySummary{
		Categories: summary,
		Total:      total,
	}, nil
}
