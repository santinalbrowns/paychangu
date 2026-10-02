package paychangu

import (
	"context"
	"errors"
	"net/url"
)

// GetBalance retrieves wallet balance for a currency (defaults to MWK when empty).
func (p *PayChangu) GetBalance(currency string) (*BalanceData, error) {
	return p.GetBalanceContext(context.Background(), currency)
}

// GetBalanceContext retrieves wallet balance with context.
func (p *PayChangu) GetBalanceContext(ctx context.Context, currency string) (*BalanceData, error) {
	q := url.Values{}
	if currency != "" {
		q.Set("currency", currency)
	}
	var response BalanceResponse
	if err := p.get(ctx, queryPath("/wallet-balance", q), &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response.Data, nil
}
