package paychangu

import (
	"context"
	"fmt"
	"net/url"
)

// GetBillers lists all available billers.
func (p *PayChangu) GetBillers() (*BillsAPIResponse, error) {
	return p.GetBillersContext(context.Background())
}

// GetBillersContext lists billers with context.
func (p *PayChangu) GetBillersContext(ctx context.Context) (*BillsAPIResponse, error) {
	var response BillsAPIResponse
	if err := p.get(ctx, "/bills/getBillers", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetBillerDetails fetches details for a biller (e.g. "escom").
func (p *PayChangu) GetBillerDetails(billerID string) (*BillsAPIResponse, error) {
	return p.GetBillerDetailsContext(context.Background(), billerID)
}

// GetBillerDetailsContext fetches biller details with context.
func (p *PayChangu) GetBillerDetailsContext(ctx context.Context, billerID string) (*BillsAPIResponse, error) {
	var response BillsAPIResponse
	path := fmt.Sprintf("/bills/getBillers/%s", url.PathEscape(billerID))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ValidateBill validates a bill account before payment.
func (p *PayChangu) ValidateBill(request ValidateBillRequest) (*BillsAPIResponse, error) {
	return p.ValidateBillContext(context.Background(), request)
}

// ValidateBillContext validates a bill with context.
func (p *PayChangu) ValidateBillContext(ctx context.Context, request ValidateBillRequest) (*BillsAPIResponse, error) {
	var response BillsAPIResponse
	if err := p.post(ctx, "/bills/validate", request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// PayBill pays a bill.
func (p *PayChangu) PayBill(request PayBillRequest) (*BillsAPIResponse, error) {
	return p.PayBillContext(context.Background(), request)
}

// PayBillContext pays a bill with context.
func (p *PayChangu) PayBillContext(ctx context.Context, request PayBillRequest) (*BillsAPIResponse, error) {
	var response BillsAPIResponse
	if err := p.post(ctx, "/bills/pay", request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// BuyAirtime recharges airtime for TNM or Airtel.
func (p *PayChangu) BuyAirtime(request AirtimeRequest) (*BillsAPIResponse, error) {
	return p.BuyAirtimeContext(context.Background(), request)
}

// BuyAirtimeContext buys airtime with context.
func (p *PayChangu) BuyAirtimeContext(ctx context.Context, request AirtimeRequest) (*BillsAPIResponse, error) {
	var response BillsAPIResponse
	if err := p.post(ctx, "/bills/buy-airtime", request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetBillTransaction fetches bill transaction details by reference.
func (p *PayChangu) GetBillTransaction(reference string) (*BillsAPIResponse, error) {
	return p.GetBillTransactionContext(context.Background(), reference)
}

// GetBillTransactionContext fetches bill transaction details with context.
func (p *PayChangu) GetBillTransactionContext(ctx context.Context, reference string) (*BillsAPIResponse, error) {
	var response BillsAPIResponse
	path := fmt.Sprintf("/bills/getTransactions/%s", url.PathEscape(reference))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetBillStatistics fetches bill payment statistics.
func (p *PayChangu) GetBillStatistics() (*BillsAPIResponse, error) {
	return p.GetBillStatisticsContext(context.Background())
}

// GetBillStatisticsContext fetches bill statistics with context.
func (p *PayChangu) GetBillStatisticsContext(ctx context.Context) (*BillsAPIResponse, error) {
	var response BillsAPIResponse
	if err := p.get(ctx, "/bills/getStatistics", &response); err != nil {
		return nil, err
	}
	return &response, nil
}
