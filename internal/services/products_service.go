package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphaelmb/go-bid/internal/store/pgstore"
)

type ProductService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductService(pool *pgxpool.Pool) ProductService {
	return ProductService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductService) CreateProduct(
	ctx context.Context,
	sellerID uuid.UUID,
	productName, description string,
	baseprice float64,
	auctionEnd time.Time,
) (uuid.UUID, error) {
	id, err := ps.queries.CreateProduct(ctx, pgstore.CreateProductParams{
		SellerID:    sellerID,
		ProductName: productName,
		Description: description,
		Baseprice:   baseprice,
		AuctionEnd:  auctionEnd,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

type Product struct {
	SellerID    uuid.UUID `json:"seller_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	Baseprice   float64   `json:"baseprice"`
	AuctionEnd  time.Time `json:"auction_end"`
}

func (ps *ProductService) GetAllProducts(ctx context.Context) ([]Product, error) {
	products, err := ps.queries.GetAllProducts(ctx)
	if err != nil {
		return []Product{}, err
	}

	return productsOutput(products), nil
}

func productsOutput(products []pgstore.Product) []Product {
	var result []Product
	for _, val := range products {
		product := Product{
			SellerID:    val.SellerID,
			ProductName: val.ProductName,
			Description: val.Description,
			Baseprice:   val.Baseprice,
			AuctionEnd:  val.AuctionEnd,
		}
		result = append(result, product)
	}
	return result
}
