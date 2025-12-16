package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"Aicon-assignment/internal/domain/entity"
	domainErrors "Aicon-assignment/internal/domain/errors"
)

// MockItemRepository はtestify/mockを使用したモックリポジトリ
type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) FindAll(ctx context.Context) ([]*entity.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Item), args.Error(1)
}

func (m *MockItemRepository) FindByID(ctx context.Context, id int64) (*entity.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemRepository) Create(ctx context.Context, item *entity.Item) (*entity.Item, error) {
	args := m.Called(ctx, item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemRepository) Update(ctx context.Context, item *entity.Item) (*entity.Item, error) {
	args := m.Called(ctx, item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemRepository) GetSummaryByCategory(ctx context.Context) (map[string]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

func TestNewItemUsecase(t *testing.T) {
	mockRepo := new(MockItemRepository)
	usecase := NewItemUsecase(mockRepo)

	assert.NotNil(t, usecase)
}

func TestItemUsecase_GetAllItems(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockItemRepository)
		expectedCount int
		expectedErr   error
	}{
		{
			name: "正常系: 複数のアイテムを取得",
			setupMock: func(mockRepo *MockItemRepository) {
				item1, _ := entity.NewItem("時計1", "時計", "ROLEX", 1000000, "2023-01-01")
				item2, _ := entity.NewItem("バッグ1", "バッグ", "HERMÈS", 500000, "2023-01-02")
				items := []*entity.Item{item1, item2}
				mockRepo.On("FindAll", mock.Anything).Return(items, nil)
			},
			expectedCount: 2,
			expectedErr:   nil,
		},
		{
			name: "正常系: アイテムが0件",
			setupMock: func(mockRepo *MockItemRepository) {
				items := []*entity.Item{}
				mockRepo.On("FindAll", mock.Anything).Return(items, nil)
			},
			expectedCount: 0,
			expectedErr:   nil,
		},
		{
			name: "異常系: データベースエラー",
			setupMock: func(mockRepo *MockItemRepository) {
				mockRepo.On("FindAll", mock.Anything).Return(([]*entity.Item)(nil), domainErrors.ErrDatabaseError)
			},
			expectedCount: 0,
			expectedErr:   domainErrors.ErrDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockItemRepository)
			tt.setupMock(mockRepo)
			usecase := NewItemUsecase(mockRepo)

			ctx := context.Background()
			items, err := usecase.GetAllItems(ctx)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedErr)
				mockRepo.AssertExpectations(t)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, items, tt.expectedCount)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestItemUsecase_GetItemByID(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		setupMock   func(*MockItemRepository)
		expectError bool
		expectedErr error
	}{
		{
			name: "正常系: 存在するアイテムを取得",
			id:   1,
			setupMock: func(mockRepo *MockItemRepository) {
				item, _ := entity.NewItem("時計1", "時計", "ROLEX", 1000000, "2023-01-01")
				item.ID = 1
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return(item, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 存在しないアイテム",
			id:   999,
			setupMock: func(mockRepo *MockItemRepository) {
				mockRepo.On("FindByID", mock.Anything, int64(999)).Return((*entity.Item)(nil), domainErrors.ErrItemNotFound)
			},
			expectError: true,
			expectedErr: domainErrors.ErrItemNotFound,
		},
		{
			name: "異常系: 無効なID（0以下）",
			id:   0,
			setupMock: func(mockRepo *MockItemRepository) {
				// FindByIDは呼ばれない
			},
			expectError: true,
			expectedErr: domainErrors.ErrInvalidInput,
		},
		{
			name: "異常系: データベースエラー",
			id:   1,
			setupMock: func(mockRepo *MockItemRepository) {
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return((*entity.Item)(nil), domainErrors.ErrDatabaseError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockItemRepository)
			tt.setupMock(mockRepo)
			usecase := NewItemUsecase(mockRepo)

			ctx := context.Background()
			item, err := usecase.GetItemByID(ctx, tt.id)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
				assert.Nil(t, item)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, item)
				assert.Equal(t, tt.id, item.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestItemUsecase_CreateItem(t *testing.T) {
	tests := []struct {
		name        string
		input       CreateItemInput
		setupMock   func(*MockItemRepository)
		expectError bool
		expectedErr error
	}{
		{
			name: "正常系: 有効なアイテムを作成",
			input: CreateItemInput{
				Name:          "ロレックス デイトナ",
				Category:      "時計",
				Brand:         "ROLEX",
				PurchasePrice: 1500000,
				PurchaseDate:  "2023-01-15",
			},
			setupMock: func(mockRepo *MockItemRepository) {
				createdItem, _ := entity.NewItem("ロレックス デイトナ", "時計", "ROLEX", 1500000, "2023-01-15")
				createdItem.ID = 1
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Item")).Return(createdItem, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 無効な入力（名前が空）",
			input: CreateItemInput{
				Name:          "",
				Category:      "時計",
				Brand:         "ROLEX",
				PurchasePrice: 1500000,
				PurchaseDate:  "2023-01-15",
			},
			setupMock: func(mockRepo *MockItemRepository) {
				// Createは呼ばれない
			},
			expectError: true,
			expectedErr: domainErrors.ErrInvalidInput,
		},
		{
			name: "異常系: 無効なカテゴリー",
			input: CreateItemInput{
				Name:          "アイテム",
				Category:      "無効なカテゴリー",
				Brand:         "ブランド",
				PurchasePrice: 100000,
				PurchaseDate:  "2023-01-15",
			},
			setupMock: func(mockRepo *MockItemRepository) {
				// Createは呼ばれない
			},
			expectError: true,
			expectedErr: domainErrors.ErrInvalidInput,
		},
		{
			name: "異常系: データベースエラー",
			input: CreateItemInput{
				Name:          "アイテム",
				Category:      "時計",
				Brand:         "ブランド",
				PurchasePrice: 100000,
				PurchaseDate:  "2023-01-15",
			},
			setupMock: func(mockRepo *MockItemRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Item")).Return((*entity.Item)(nil), domainErrors.ErrDatabaseError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockItemRepository)
			tt.setupMock(mockRepo)
			usecase := NewItemUsecase(mockRepo)

			ctx := context.Background()
			item, err := usecase.CreateItem(ctx, tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
				assert.Nil(t, item)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, item)
				assert.Equal(t, tt.input.Name, item.Name)
				assert.Equal(t, tt.input.Category, item.Category)
				assert.Equal(t, tt.input.Brand, item.Brand)
				assert.Equal(t, tt.input.PurchasePrice, item.PurchasePrice)
				assert.Equal(t, tt.input.PurchaseDate, item.PurchaseDate)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestItemUsecase_DeleteItem(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		setupMock   func(*MockItemRepository)
		expectError bool
		expectedErr error
	}{
		{
			name: "正常系: 存在するアイテムを削除",
			id:   1,
			setupMock: func(mockRepo *MockItemRepository) {
				item, _ := entity.NewItem("時計1", "時計", "ROLEX", 1000000, "2023-01-01")
				item.ID = 1
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return(item, nil)
				mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 存在しないアイテム",
			id:   999,
			setupMock: func(mockRepo *MockItemRepository) {
				mockRepo.On("FindByID", mock.Anything, int64(999)).Return((*entity.Item)(nil), domainErrors.ErrItemNotFound)
			},
			expectError: true,
			expectedErr: domainErrors.ErrItemNotFound,
		},
		{
			name: "異常系: 無効なID（0以下）",
			id:   0,
			setupMock: func(mockRepo *MockItemRepository) {
				// FindByIDは呼ばれない
			},
			expectError: true,
			expectedErr: domainErrors.ErrInvalidInput,
		},
		{
			name: "異常系: FindByIDでデータベースエラー",
			id:   1,
			setupMock: func(mockRepo *MockItemRepository) {
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return((*entity.Item)(nil), domainErrors.ErrDatabaseError)
			},
			expectError: true,
		},
		{
			name: "異常系: Deleteでデータベースエラー",
			id:   1,
			setupMock: func(mockRepo *MockItemRepository) {
				item, _ := entity.NewItem("時計1", "時計", "ROLEX", 1000000, "2023-01-01")
				item.ID = 1
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return(item, nil)
				mockRepo.On("Delete", mock.Anything, int64(1)).Return(domainErrors.ErrDatabaseError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockItemRepository)
			tt.setupMock(mockRepo)
			usecase := NewItemUsecase(mockRepo)

			ctx := context.Background()
			err := usecase.DeleteItem(ctx, tt.id)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestItemUsecase_UpdateItem(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		input       UpdateItemInput
		setupMock   func(*MockItemRepository)
		expectError bool
		expectedErr error
	}{
		{
			name: "正常系: 全フィールド更新",
			id:   1,
			input: UpdateItemInput{
				Name:          stringPtr("更新された商品名"),
				Brand:         stringPtr("更新されたブランド"),
				PurchasePrice: intPtr(2000000),
			},
			setupMock: func(mockRepo *MockItemRepository) {
				existingItem, _ := entity.NewItem("既存商品", "時計", "既存ブランド", 1000000, "2023-01-01")
				existingItem.ID = 1
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return(existingItem, nil)

				updatedItem, _ := entity.NewItem("更新された商品名", "時計", "更新されたブランド", 2000000, "2023-01-01")
				updatedItem.ID = 1
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Item")).Return(updatedItem, nil)
			},
			expectError: false,
		},
		{
			name: "正常系: nameのみ更新",
			id:   1,
			input: UpdateItemInput{
				Name: stringPtr("新しい商品名"),
			},
			setupMock: func(mockRepo *MockItemRepository) {
				existingItem, _ := entity.NewItem("既存商品", "時計", "既存ブランド", 1000000, "2023-01-01")
				existingItem.ID = 1
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return(existingItem, nil)

				updatedItem, _ := entity.NewItem("新しい商品名", "時計", "既存ブランド", 1000000, "2023-01-01")
				updatedItem.ID = 1
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Item")).Return(updatedItem, nil)
			},
			expectError: false,
		},
		{
			name: "正常系: brandのみ更新",
			id:   1,
			input: UpdateItemInput{
				Brand: stringPtr("新しいブランド"),
			},
			setupMock: func(mockRepo *MockItemRepository) {
				existingItem, _ := entity.NewItem("既存商品", "時計", "既存ブランド", 1000000, "2023-01-01")
				existingItem.ID = 1
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return(existingItem, nil)

				updatedItem, _ := entity.NewItem("既存商品", "時計", "新しいブランド", 1000000, "2023-01-01")
				updatedItem.ID = 1
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Item")).Return(updatedItem, nil)
			},
			expectError: false,
		},
		{
			name: "正常系: purchase_priceのみ更新",
			id:   1,
			input: UpdateItemInput{
				PurchasePrice: intPtr(3000000),
			},
			setupMock: func(mockRepo *MockItemRepository) {
				existingItem, _ := entity.NewItem("既存商品", "時計", "既存ブランド", 1000000, "2023-01-01")
				existingItem.ID = 1
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return(existingItem, nil)

				updatedItem, _ := entity.NewItem("既存商品", "時計", "既存ブランド", 3000000, "2023-01-01")
				updatedItem.ID = 1
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Item")).Return(updatedItem, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 無効なID（0以下）",
			id:   0,
			input: UpdateItemInput{
				Name: stringPtr("更新された商品名"),
			},
			setupMock: func(mockRepo *MockItemRepository) {
				// FindByIDは呼ばれない
			},
			expectError: true,
			expectedErr: domainErrors.ErrInvalidInput,
		},
		{
			name: "異常系: 存在しないアイテム",
			id:   999,
			input: UpdateItemInput{
				Name: stringPtr("更新された商品名"),
			},
			setupMock: func(mockRepo *MockItemRepository) {
				mockRepo.On("FindByID", mock.Anything, int64(999)).Return((*entity.Item)(nil), domainErrors.ErrItemNotFound)
			},
			expectError: true,
			expectedErr: domainErrors.ErrItemNotFound,
		},
		{
			name: "異常系: バリデーションエラー（空のname）",
			id:   1,
			input: UpdateItemInput{
				Name: stringPtr(""),
			},
			setupMock: func(mockRepo *MockItemRepository) {
				existingItem, _ := entity.NewItem("既存商品", "時計", "既存ブランド", 1000000, "2023-01-01")
				existingItem.ID = 1
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return(existingItem, nil)
			},
			expectError: true,
			expectedErr: domainErrors.ErrInvalidInput,
		},
		{
			name: "異常系: バリデーションエラー（負のpurchase_price）",
			id:   1,
			input: UpdateItemInput{
				PurchasePrice: intPtr(-100),
			},
			setupMock: func(mockRepo *MockItemRepository) {
				existingItem, _ := entity.NewItem("既存商品", "時計", "既存ブランド", 1000000, "2023-01-01")
				existingItem.ID = 1
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return(existingItem, nil)
			},
			expectError: true,
			expectedErr: domainErrors.ErrInvalidInput,
		},
		{
			name: "異常系: FindByIDでデータベースエラー",
			id:   1,
			input: UpdateItemInput{
				Name: stringPtr("更新された商品名"),
			},
			setupMock: func(mockRepo *MockItemRepository) {
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return((*entity.Item)(nil), domainErrors.ErrDatabaseError)
			},
			expectError: true,
		},
		{
			name: "異常系: Updateでデータベースエラー",
			id:   1,
			input: UpdateItemInput{
				Name: stringPtr("更新された商品名"),
			},
			setupMock: func(mockRepo *MockItemRepository) {
				existingItem, _ := entity.NewItem("既存商品", "時計", "既存ブランド", 1000000, "2023-01-01")
				existingItem.ID = 1
				mockRepo.On("FindByID", mock.Anything, int64(1)).Return(existingItem, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Item")).Return((*entity.Item)(nil), domainErrors.ErrDatabaseError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockItemRepository)
			tt.setupMock(mockRepo)
			usecase := NewItemUsecase(mockRepo)

			ctx := context.Background()
			item, err := usecase.UpdateItem(ctx, tt.id, tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
				assert.Nil(t, item)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, item)

				if tt.input.Name != nil {
					assert.Equal(t, *tt.input.Name, item.Name)
				}
				if tt.input.Brand != nil {
					assert.Equal(t, *tt.input.Brand, item.Brand)
				}
				if tt.input.PurchasePrice != nil {
					assert.Equal(t, *tt.input.PurchasePrice, item.PurchasePrice)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// ヘルパー関数
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func TestItemUsecase_GetCategorySummary(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*MockItemRepository)
		expectedTotal      int
		expectedWatchCount int
		expectedBagCount   int
		expectError        bool
	}{
		{
			name: "正常系: 複数カテゴリーのアイテムがある場合",
			setupMock: func(mockRepo *MockItemRepository) {
				summary := map[string]int{
					"時計":  2,
					"バッグ": 1,
				}
				mockRepo.On("GetSummaryByCategory", mock.Anything).Return(summary, nil)
			},
			expectedTotal:      3,
			expectedWatchCount: 2,
			expectedBagCount:   1,
			expectError:        false,
		},
		{
			name: "正常系: アイテムが0件の場合",
			setupMock: func(mockRepo *MockItemRepository) {
				summary := map[string]int{}
				mockRepo.On("GetSummaryByCategory", mock.Anything).Return(summary, nil)
			},
			expectedTotal:      0,
			expectedWatchCount: 0,
			expectedBagCount:   0,
			expectError:        false,
		},
		{
			name: "異常系: データベースエラー",
			setupMock: func(mockRepo *MockItemRepository) {
				mockRepo.On("GetSummaryByCategory", mock.Anything).Return((map[string]int)(nil), domainErrors.ErrDatabaseError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockItemRepository)
			tt.setupMock(mockRepo)
			usecase := NewItemUsecase(mockRepo)

			ctx := context.Background()
			summary, err := usecase.GetCategorySummary(ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, summary)
				mockRepo.AssertExpectations(t)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, summary)

			assert.Equal(t, tt.expectedTotal, summary.Total)
			assert.Equal(t, tt.expectedWatchCount, summary.Categories["時計"])
			assert.Equal(t, tt.expectedBagCount, summary.Categories["バッグ"])

			// すべてのカテゴリーがレスポンスに含まれているかチェック
			expectedCategories := []string{"時計", "バッグ", "ジュエリー", "靴", "その他"}
			for _, category := range expectedCategories {
				assert.Contains(t, summary.Categories, category)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
