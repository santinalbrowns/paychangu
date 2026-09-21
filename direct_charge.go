package paychangu

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// ChargeMobileMoney initiates a direct mobile money collection charge.
func (p *PayChangu) ChargeMobileMoney(request ChargeMobileMoneyRequest) (*DirectChargeDetails, error) {
	return p.ChargeMobileMoneyContext(context.Background(), request)
}

// ChargeMobileMoneyContext initiates a MoMo charge with context.
func (p *PayChangu) ChargeMobileMoneyContext(ctx context.Context, request ChargeMobileMoneyRequest) (*DirectChargeDetails, error) {
	var response ChargeMobileMoneyResponse
	if err := p.post(ctx, "/mobile-money/payments/initialize", request, &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response.Data, nil
}

// VerifyDirectCharge verifies a direct MoMo charge by charge ID.
func (p *PayChangu) VerifyDirectCharge(chargeID string) (*DirectChargeDetails, error) {
	return p.VerifyDirectChargeContext(context.Background(), chargeID)
}

// VerifyDirectChargeContext verifies a direct charge with context.
func (p *PayChangu) VerifyDirectChargeContext(ctx context.Context, chargeID string) (*DirectChargeDetails, error) {
	var response DirectChargeStatusResponse
	path := fmt.Sprintf("/mobile-money/payments/%s/verify", url.PathEscape(chargeID))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response.Data, nil
}

// GetChargeDetails fetches direct MoMo charge details by charge ID.
func (p *PayChangu) GetChargeDetails(chargeID string) (*DirectChargeDetails, error) {
	return p.GetChargeDetailsContext(context.Background(), chargeID)
}

// GetChargeDetailsContext fetches charge details with context.
func (p *PayChangu) GetChargeDetailsContext(ctx context.Context, chargeID string) (*DirectChargeDetails, error) {
	var response DirectChargeStatusResponse
	path := fmt.Sprintf("/mobile-money/payments/%s/details", url.PathEscape(chargeID))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response.Data, nil
}

// InitiateBankTransfer creates a bank transfer collection (virtual account).
func (p *PayChangu) InitiateBankTransfer(request BankTransferRequest) (*BankTransferResponse, error) {
	return p.InitiateBankTransferContext(context.Background(), request)
}

// InitiateBankTransferContext initiates a bank transfer charge with context.
func (p *PayChangu) InitiateBankTransferContext(ctx context.Context, request BankTransferRequest) (*BankTransferResponse, error) {
	if request.PaymentMethod == "" {
		request.PaymentMethod = "mobile_bank_transfer"
	}
	var response BankTransferResponse
	if err := p.post(ctx, "/direct-charge/payments/initialize", request, &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response, nil
}

// GetBankTransferDetails fetches a bank transfer transaction by charge ID.
func (p *PayChangu) GetBankTransferDetails(chargeID string) (*BankTransferTransaction, error) {
	return p.GetBankTransferDetailsContext(context.Background(), chargeID)
}

// GetBankTransferDetailsContext fetches bank transfer details with context.
func (p *PayChangu) GetBankTransferDetailsContext(ctx context.Context, chargeID string) (*BankTransferTransaction, error) {
	var response BankTransferDetailsResponse
	path := fmt.Sprintf("/direct-charge/transactions/%s/details", url.PathEscape(chargeID))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response.Data.Transaction, nil
}
