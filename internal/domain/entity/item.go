package entity

import (
	"errors"
	"strings"
	"time"
)

type Item struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	Brand         string    `json:"brand"`
	PurchasePrice int       `json:"purchase_price"`
	PurchaseDate  string    `json:"purchase_date"` // YYYY-MM-DD 形式
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// カテゴリー定義
var ValidCategories = []string{"時計", "バッグ", "ジュエリー", "靴", "その他"}

func NewItem(name, category, brand string, purchasePrice int, purchaseDate string) (*Item, error) {
	item := &Item{
		Name:          strings.TrimSpace(name),
		Category:      strings.TrimSpace(category),
		Brand:         strings.TrimSpace(brand),
		PurchasePrice: purchasePrice,
		PurchaseDate:  strings.TrimSpace(purchaseDate),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := item.Validate(); err != nil {
		return nil, err
	}

	return item, nil
}

// アイテムフィールドのバリデーション
func (i *Item) Validate() error {
	var errs []string

	if i.Name == "" {
		errs = append(errs, "name is required")
	} else if len(i.Name) > 100 {
		errs = append(errs, "name must be 100 characters or less")
	}

	if i.Category == "" {
		errs = append(errs, "category is required")
	} else if !isValidCategory(i.Category) {
		errs = append(errs, "category must be one of: 時計, バッグ, ジュエリー, 靴, その他")
	}

	if i.Brand == "" {
		errs = append(errs, "brand is required")
	} else if len(i.Brand) > 100 {
		errs = append(errs, "brand must be 100 characters or less")
	}

	if i.PurchasePrice < 0 {
		errs = append(errs, "purchase_price must be 0 or greater")
	}

	if i.PurchaseDate == "" {
		errs = append(errs, "purchase_date is required")
	} else if !isValidDateFormat(i.PurchaseDate) {
		errs = append(errs, "purchase_date must be in YYYY-MM-DD format")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

// アイテムフィールドのアップデート
func (i *Item) Update(name, category, brand string, purchasePrice int, purchaseDate string) error {
	i.Name = strings.TrimSpace(name)
	i.Category = strings.TrimSpace(category)
	i.Brand = strings.TrimSpace(brand)
	i.PurchasePrice = purchasePrice
	i.PurchaseDate = strings.TrimSpace(purchaseDate)
	i.UpdatedAt = time.Now()

	return i.Validate()
}

// カテゴリーのバリデーション
func isValidCategory(category string) bool {
	for _, valid := range ValidCategories {
		if category == valid {
			return true
		}
	}
	return false
}

// デート形式のバリデーション
func isValidDateFormat(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// カテゴリーの取得
func GetValidCategories() []string {
	return ValidCategories
}
