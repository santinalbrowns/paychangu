package paychangu

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// GetMobileMoneyOperators lists supported mobile money operators.
func (p *PayChangu) GetMobileMoneyOperators() ([]MobileMoneyOperator, error) {
	return p.GetMobileMoneyOperatorsContext(context.Background())
}

// GetMobileMoneyOperatorsContext lists operators with context.
func (p *PayChangu) GetMobileMoneyOperatorsContext(ctx context.Context) ([]MobileMoneyOperator, error) {
	var response MobileMoneyOperatorsResponse
	if err := p.get(ctx, "/mobile-money", &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return response.Data, nil
}

// InitiateMobileMoneyPayout sends funds to a mobile money wallet.
func (p *PayChangu) InitiateMobileMoneyPayout(request MobileMoneyPayoutRequest) (*MobileMoneyPayoutResponse, error) {
	return p.InitiateMobileMoneyPayoutContext(context.Background(), request)
}

// InitiateMobileMoneyPayoutContext sends a MoMo payout with context.
func (p *PayChangu) InitiateMobileMoneyPayoutContext(ctx context.Context, request MobileMoneyPayoutRequest) (*MobileMoneyPayoutResponse, error) {
	var response MobileMoneyPayoutResponse
	if err := p.post(ctx, "/mobile-money/payouts/initialize", request, &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response, nil
}

// GetMobileMoneyPayoutDetails fetches MoMo payout details by charge ID.
func (p *PayChangu) GetMobileMoneyPayoutDetails(chargeID string) (*PayoutTransactionDetails, error) {
	return p.GetMobileMoneyPayoutDetailsContext(context.Background(), chargeID)
}

// GetMobileMoneyPayoutDetailsContext fetches MoMo payout details with context.
func (p *PayChangu) GetMobileMoneyPayoutDetailsContext(ctx context.Context, chargeID string) (*PayoutTransactionDetails, error) {
	var response GetMobileMoneyPayoutDetailsResponse
	path := fmt.Sprintf("/mobile-money/payments/%s/details", url.PathEscape(chargeID))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response.Data, nil
}

// GetSupportedBanks lists banks available for payouts in a currency.
func (p *PayChangu) GetSupportedBanks(currency string) ([]Bank, error) {
	return p.GetSupportedBanksContext(context.Background(), currency)
}

// GetSupportedBanksContext lists supported banks with context.
func (p *PayChangu) GetSupportedBanksContext(ctx context.Context, currency string) ([]Bank, error) {
	q := url.Values{}
	q.Set("currency", currency)
	var response BanksResponse
	if err := p.get(ctx, queryPath("/direct-charge/payouts/supported-banks", q), &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return response.Data, nil
}

// InitiateBankPayout sends funds to a bank account.
func (p *PayChangu) InitiateBankPayout(request BankPayoutRequest) (*BankPayoutResponse, error) {
	return p.InitiateBankPayoutContext(context.Background(), request)
}

// InitiateBankPayoutContext sends a bank payout with context.
func (p *PayChangu) InitiateBankPayoutContext(ctx context.Context, request BankPayoutRequest) (*BankPayoutResponse, error) {
	payload := struct {
		PayoutMethod      string `json:"payout_method"`
		BankUUID          string `json:"bank_uuid"`
		Amount            string `json:"amount"`
		ChargeID          string `json:"charge_id"`
		BankAccountName   string `json:"bank_account_name"`
		BankAccountNumber string `json:"bank_account_number"`
		Email             string `json:"email,omitempty"`
		FirstName         string `json:"first_name,omitempty"`
		LastName          string `json:"last_name,omitempty"`
	}{
		PayoutMethod:      request.PayoutMethod,
		BankUUID:          request.BankUUID,
		Amount:            fmt.Sprintf("%.2f", request.Amount),
		ChargeID:          request.ChargeID,
		BankAccountName:   request.BankAccountName,
		BankAccountNumber: request.BankAccountNumber,
		Email:             request.Email,
		FirstName:         request.FirstName,
		LastName:          request.LastName,
	}

	var response BankPayoutResponse
	if err := p.post(ctx, "/direct-charge/payouts/initialize", payload, &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response, nil
}

// GetBankPayoutDetails fetches bank payout details by charge ID.
func (p *PayChangu) GetBankPayoutDetails(chargeID string) (*BankPayoutTransactionDetails, error) {
	return p.GetBankPayoutDetailsContext(context.Background(), chargeID)
}

// GetBankPayoutDetailsContext fetches bank payout details with context.
func (p *PayChangu) GetBankPayoutDetailsContext(ctx context.Context, chargeID string) (*BankPayoutTransactionDetails, error) {
	var response GetBankPayoutDetailsResponse
	path := fmt.Sprintf("/direct-charge/payouts/%s/details", url.PathEscape(chargeID))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response.Data, nil
}

// ListBankPayouts lists all bank payouts.
func (p *PayChangu) ListBankPayouts() (*PaginatedBankPayouts, error) {
	return p.ListBankPayoutsContext(context.Background())
}

// ListBankPayoutsContext lists all bank payouts with context.
func (p *PayChangu) ListBankPayoutsContext(ctx context.Context) (*PaginatedBankPayouts, error) {
	var response ListBankPayoutsResponse
	if err := p.get(ctx, "/direct-charge/payouts", &response); err != nil {
		return nil, err
	}
	if !isSuccessStatus(response.Status) {
		return nil, errors.New(response.Message)
	}
	return &response.Data, nil
}
