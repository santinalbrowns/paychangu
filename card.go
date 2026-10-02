package paychangu

import (
	"context"
	"fmt"
	"net/url"
)

// ChargeCard charges a card. May return a 3DS auth link.
func (p *PayChangu) ChargeCard(request ChargeCardRequest) (*ChargeCardResponse, error) {
	return p.ChargeCardContext(context.Background(), request)
}

// ChargeCardContext charges a card with context.
func (p *PayChangu) ChargeCardContext(ctx context.Context, request ChargeCardRequest) (*ChargeCardResponse, error) {
	var response ChargeCardResponse
	if err := p.post(ctx, "/charge-card/payments", request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// VerifyCardCharge verifies a card charge by charge ID.
func (p *PayChangu) VerifyCardCharge(chargeID string) (*CardChargeStatusResponse, error) {
	return p.VerifyCardChargeContext(context.Background(), chargeID)
}

// VerifyCardChargeContext verifies a card charge with context.
func (p *PayChangu) VerifyCardChargeContext(ctx context.Context, chargeID string) (*CardChargeStatusResponse, error) {
	var response CardChargeStatusResponse
	path := fmt.Sprintf("/charge-card/verify/%s", url.PathEscape(chargeID))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// RefundCardCharge refunds a card charge by charge ID.
func (p *PayChangu) RefundCardCharge(chargeID string) (*CardChargeStatusResponse, error) {
	return p.RefundCardChargeContext(context.Background(), chargeID)
}

// RefundCardChargeContext refunds a card charge with context.
func (p *PayChangu) RefundCardChargeContext(ctx context.Context, chargeID string) (*CardChargeStatusResponse, error) {
	var response CardChargeStatusResponse
	path := fmt.Sprintf("/charge-card/refund/%s", url.PathEscape(chargeID))
	if err := p.post(ctx, path, map[string]any{}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
