package dynamodb

import (
	"context"
	"fmt"
	"time"

	"ec-sample/services/cart-service/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	cartTableName = "carts"
)

// DynamoDBRepository DynamoDB実装
type DynamoDBRepository struct {
	client *dynamodb.Client
}

// NewDynamoDBRepository コンストラクタ
func NewDynamoDBRepository(client *dynamodb.Client) domain.CartRepository {
	return &DynamoDBRepository{client: client}
}

// GetCart カート取得
func (r *DynamoDBRepository) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(cartTableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var cart domain.Cart
	err = attributevalue.UnmarshalMap(result.Item, &cart)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart: %w", err)
	}

	return &cart, nil
}

// SaveCart カート保存
func (r *DynamoDBRepository) SaveCart(ctx context.Context, cart *domain.Cart) error {
	item, err := attributevalue.MarshalMap(cart)
	if err != nil {
		return fmt.Errorf("failed to marshal cart: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(cartTableName),
		Item:      item,
	})

	return err
}

// AddItem アイテム追加
func (r *DynamoDBRepository) AddItem(ctx context.Context, userID string, item *domain.CartItem) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	if cart == nil {
		cart = &domain.Cart{
			UserID:       userID,
			Items:        []*domain.CartItem{},
			Total:        0,
			ItemCount:    0,
			LastModified: time.Now(),
			ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		}
	}

	cart.Items = append(cart.Items, item)
	return r.SaveCart(ctx, cart)
}

// UpdateItem アイテム更新
func (r *DynamoDBRepository) UpdateItem(ctx context.Context, userID, productID string, quantity int) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	if cart == nil {
		return fmt.Errorf("cart not found")
	}

	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].Quantity = quantity
			break
		}
	}

	return r.SaveCart(ctx, cart)
}

// RemoveItem アイテム削除
func (r *DynamoDBRepository) RemoveItem(ctx context.Context, userID, productID string) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	if cart == nil {
		return fmt.Errorf("cart not found")
	}

	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			break
		}
	}

	return r.SaveCart(ctx, cart)
}

// ClearCart カートクリア
func (r *DynamoDBRepository) ClearCart(ctx context.Context, userID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(cartTableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
	})
	return err
}

// GetCartByID IDでカート取得
func (r *DynamoDBRepository) GetCartByID(ctx context.Context, userID string) (*domain.Cart, error) {
	return r.GetCart(ctx, userID)
}
