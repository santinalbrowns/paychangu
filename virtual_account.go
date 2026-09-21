package paychangu

import (
	"context"
	"fmt"
	"net/url"
)

// CreateCustomer creates a customer for USD virtual accounts.
func (p *PayChangu) CreateCustomer(request VirtualCustomerRequest) (*VirtualCustomerResponse, error) {
	return p.CreateCustomerContext(context.Background(), request)
}

// CreateCustomerContext creates a virtual-account customer with context.
func (p *PayChangu) CreateCustomerContext(ctx context.Context, request VirtualCustomerRequest) (*VirtualCustomerResponse, error) {
	var response VirtualCustomerResponse
	if err := p.post(ctx, "/virtual-account/api/customers/create", request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetCustomers lists virtual-account customers.
func (p *PayChangu) GetCustomers(page, perPage string) (*VirtualCustomersResponse, error) {
	return p.GetCustomersContext(context.Background(), page, perPage)
}

// GetCustomersContext lists customers with context.
func (p *PayChangu) GetCustomersContext(ctx context.Context, page, perPage string) (*VirtualCustomersResponse, error) {
	q := url.Values{}
	if page != "" {
		q.Set("page", page)
	}
	if perPage != "" {
		q.Set("per_page", perPage)
	}
	var response VirtualCustomersResponse
	if err := p.get(ctx, queryPath("/virtual-account/api/customers", q), &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetCustomer fetches a single virtual-account customer.
func (p *PayChangu) GetCustomer(customerID string) (*VirtualCustomerResponse, error) {
	return p.GetCustomerContext(context.Background(), customerID)
}

// GetCustomerContext fetches a customer with context.
func (p *PayChangu) GetCustomerContext(ctx context.Context, customerID string) (*VirtualCustomerResponse, error) {
	var response VirtualCustomerResponse
	path := fmt.Sprintf("/virtual-account/api/customers/%s", url.PathEscape(customerID))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// UpdateCustomer updates a virtual-account customer.
func (p *PayChangu) UpdateCustomer(customerID string, request VirtualCustomerRequest) (*VirtualCustomerResponse, error) {
	return p.UpdateCustomerContext(context.Background(), customerID, request)
}

// UpdateCustomerContext updates a customer with context.
func (p *PayChangu) UpdateCustomerContext(ctx context.Context, customerID string, request VirtualCustomerRequest) (*VirtualCustomerResponse, error) {
	var response VirtualCustomerResponse
	path := fmt.Sprintf("/virtual-account/api/customers/%s", url.PathEscape(customerID))
	if err := p.put(ctx, path, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// DeleteCustomer deletes a virtual-account customer.
func (p *PayChangu) DeleteCustomer(customerID string) error {
	return p.DeleteCustomerContext(context.Background(), customerID)
}

// DeleteCustomerContext deletes a customer with context.
func (p *PayChangu) DeleteCustomerContext(ctx context.Context, customerID string) error {
	path := fmt.Sprintf("/virtual-account/api/customers/%s", url.PathEscape(customerID))
	return p.delete(ctx, path, nil)
}

// CreateUSAccount creates a US virtual account for a customer.
// Docs expose this as GET /virtual-account/api/customers/{customerId}/virtual-account.
func (p *PayChangu) CreateUSAccount(customerID string) (*VirtualAccountResponse, error) {
	return p.CreateUSAccountContext(context.Background(), customerID)
}

// CreateUSAccountContext creates a US account with context.
func (p *PayChangu) CreateUSAccountContext(ctx context.Context, customerID string) (*VirtualAccountResponse, error) {
	var response VirtualAccountResponse
	path := fmt.Sprintf("/virtual-account/api/customers/%s/virtual-account", url.PathEscape(customerID))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// DeactivateUSAccount deactivates a US virtual account.
func (p *PayChangu) DeactivateUSAccount(customerID string) (*VirtualAccountResponse, error) {
	return p.DeactivateUSAccountContext(context.Background(), customerID)
}

// DeactivateUSAccountContext deactivates a US account with context.
func (p *PayChangu) DeactivateUSAccountContext(ctx context.Context, customerID string) (*VirtualAccountResponse, error) {
	var response VirtualAccountResponse
	path := fmt.Sprintf("/virtual-account/api/customers/%s/virtual-account/deactivate", url.PathEscape(customerID))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ReactivateUSAccount reactivates a previously deactivated US account.
func (p *PayChangu) ReactivateUSAccount(customerID string) (*VirtualAccountResponse, error) {
	return p.ReactivateUSAccountContext(context.Background(), customerID)
}

// ReactivateUSAccountContext reactivates a US account with context.
func (p *PayChangu) ReactivateUSAccountContext(ctx context.Context, customerID string) (*VirtualAccountResponse, error) {
	var response VirtualAccountResponse
	path := fmt.Sprintf("/virtual-account/api/customers/%s/virtual-account/reactivate", url.PathEscape(customerID))
	if err := p.post(ctx, path, map[string]any{}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetUSAccountActivity fetches activity history for a US virtual account.
func (p *PayChangu) GetUSAccountActivity(customerID string) (*VirtualAccountResponse, error) {
	return p.GetUSAccountActivityContext(context.Background(), customerID)
}

// GetUSAccountActivityContext fetches US account activity with context.
func (p *PayChangu) GetUSAccountActivityContext(ctx context.Context, customerID string) (*VirtualAccountResponse, error) {
	var response VirtualAccountResponse
	path := fmt.Sprintf("/virtual-account/api/customers/%s/virtual-account/activities", url.PathEscape(customerID))
	if err := p.get(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
