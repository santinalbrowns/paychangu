package paychangu

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// InitiatePayment starts a standard checkout payment.
func (p *PayChangu) InitiatePayment(request Request) (*Response, error) {
	return p.InitiatePaymentContext(context.Background(), request)
}

// InitiatePaymentContext starts a standard checkout payment with context.
func (p *PayChangu) InitiatePaymentContext(ctx context.Context, request Request) (*Response, error) {
	var response Response
	if err := p.do(ctx, http.MethodPost, "/payment", request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// VerifyPayment verifies a standard checkout payment by tx_ref.
func (p *PayChangu) VerifyPayment(txRef string) (*VerifyPaymentResponse, error) {
	return p.VerifyPaymentContext(context.Background(), txRef)
}

// VerifyPaymentContext verifies a standard checkout payment with context.
func (p *PayChangu) VerifyPaymentContext(ctx context.Context, txRef string) (*VerifyPaymentResponse, error) {
	var response VerifyPaymentResponse
	path := fmt.Sprintf("/verify-payment/%s", txRef)
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response, nil
}
