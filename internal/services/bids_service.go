package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphaelmb/go-bid/internal/store/pgstore"
)

var ErrBidIsTooLow = errors.New("the bid value is too low")

type BidService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewBidService(pool *pgxpool.Pool) BidService {
	return BidService{pool: pool, queries: pgstore.New(pool)}
}

func (bs *BidService) PlaceBid(ctx context.Context, productId, bidderId uuid.UUID, amount float64) (pgstore.Bid, error) {
	product, err := bs.queries.GetProductById(ctx, productId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
		return pgstore.Bid{}, fmt.Errorf("error fetching product: %w", err)
	}

	highestBid, err := bs.queries.GetHighestBidByProductId(ctx, productId)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}

	if product.Baseprice >= amount || highestBid.BidAmount >= amount {
		return pgstore.Bid{}, ErrBidIsTooLow
	}

	highestBid, err = bs.queries.CreateBid(ctx, pgstore.CreateBidParams{
		ProductID: productId,
		BidderID:  bidderId,
		BidAmount: amount,
	})
	if err != nil {
		return pgstore.Bid{}, err
	}

	return highestBid, nil
}
