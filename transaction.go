package paychangu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// The payChangu struct represents a client
// for the PayChangu API. It holds an API key
// required for authenticating requests.
type payChangu struct {
	// secretkey is the secret API secretkey for
	// authentication with the PayChangu API.
	secretkey string
}

// The New function initializes
// a new instance of the payChangu client.
//
// secretKey (string): The secret API key used to authenticate with PayChangu.
//
// A pointer to a new payChangu instance, configured with the provided API key.
func New(secretKey string) *payChangu {
	return &payChangu{secretkey: secretKey}
}

// The InitiatePayment method sends a payment initiation request to the
// PayChangu API. It marshals the request data to JSON, sends it as a POST
// request, and parses the response.
//
// Parameters:
//
// request (Request): The payment request payload, containing necessary
// details such as the amount, currency, and customer information.
//
// Returns:
//
// *Response: A pointer to a Response struct containing details about the initiated payment.
//
// error: An error, if one occurred during the request. This can something return a
// JSON object but this implemention only return it as a string
//
// Example Usage
//
//	// Field appears in JSON as key "myName".
//	client 	:= transaction.New("your_secret_key")
//	req 	:= transaction.Request{Amount: 100, Currency: "MWK", FirstName: "John", ...}
//	resp, err := client.InitiatePayment(req)
//	if err != nil {
//		log.Fatalf("Payment initiation failed: %v", err)
//	}
//	fmt.Printf("Payment successful, redirect to: %s\n", resp.Data.CheckoutURL)
func (p *payChangu) InitiatePayment(request Request) (*Response, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.paychangu.com/payment", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.secretkey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bo, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(bo))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		return nil, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// The VerifyPayment method sends a request to verify the status
// of a specific payment using its transaction reference (txRef)
//
// Parameters:
//
// txRef (string): The unique transaction reference of the payment to be verified.
//
// Returns:
//
// *VerifyPaymentResponse: A pointer to a VerifyPaymentResponse struct
// containing the verification details of the payment.
//
// error: An error, if one occurred during the verification.
//
// Example Usage:
//
//	client := transaction.New("your_secret_key")
//	verifyResp, err := client.VerifyPayment("TX12345ABC")
//	if err != nil {
//		log.Fatalf("Payment verification failed: %v", err)
//	}
//	fmt.Printf("Payment status: %s\n", verifyResp.Data.Status)
func (p *payChangu) VerifyPayment(txRef string) (*VerifyPaymentResponse, error) {
	url := fmt.Sprintf("https://api.paychangu.com/verify-payment/%s", txRef)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.secretkey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var response Error
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(response.Message)
	}

	var response VerifyPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, errors.New(response.Message)
	}

	return &response, nil
}

// GetMobileMoneyOperators retrieves a list of supported mobile money operators.
//
// Returns:
//
// []MobileMoneyOperator: A slice of supported mobile money operators.
//
// error: An error, if one occurred during the request.
//
// Example Usage:
//
//	client := paychangu.New("your_secret_key")
//	operators, err := client.GetMobileMoneyOperators()
//	if err != nil {
//	    log.Fatalf("Failed to get mobile money operators: %v", err)
//	}
//	for _, op := range operators {
//	    fmt.Printf("Operator: %s (Ref ID: %s)\n", op.Name, op.RefID)
//	}
func (p *payChangu) GetMobileMoneyOperators() ([]MobileMoneyOperator, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.paychangu.com/mobile-money", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.secretkey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bo, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response body: %w", err)
		}
		// Attempt to unmarshal into general Error struct for consistent error messages
		var apiErr Error
		if jsonErr := json.Unmarshal(bo, &apiErr); jsonErr == nil && apiErr.Message != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, apiErr.Message)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bo))
	}

	var response MobileMoneyOperatorsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Status != "success" {
		return nil, errors.New(response.Message)
	}

	return response.Data, nil
}

// InitiateMobileMoneyPayout sends a mobile money payout request to the PayChangu API.
//
// Parameters:
//
// request (MobileMoneyPayoutRequest): The payout request payload.
//
// Returns:
//
// *MobileMoneyPayoutResponse: A pointer to a MobileMoneyPayoutResponse struct containing payout details.
//
// error: An error, if one occurred during the request. This can include detailed validation errors.
//
// Example Usage:
//
//	client := paychangu.New("your_secret_key")
//	payoutReq := paychangu.MobileMoneyPayoutRequest{
//	    Mobile:                    "265888123456",
//	    MobileMoneyOperatorRefID:  "27494cb5-ba9e-437f-a114-4e7a7686bcca", // TNM Mpamba ref_id
//	    Amount:                    1000.50,
//	    ChargeID:                  "MM_PAYOUT_12345",
//	    Email:                     "recipient@example.com",
//	    FirstName:                 "Jane",
//	    LastName:                  "Doe",
//	    // TransactionStatus: "successful", // Use for sandbox testing
//	}
//	payoutResp, err := client.InitiateMobileMoneyPayout(payoutReq)
//	if err != nil {
//	    log.Fatalf("Mobile money payout failed: %v", err)
//	}
//	fmt.Printf("Mobile Money Payout Initiated. Ref ID: %s, Status: %s\n", payoutResp.Data.Transaction.RefID, payoutResp.Data.Transaction.Status)
func (p *payChangu) InitiateMobileMoneyPayout(request MobileMoneyPayoutRequest) (*MobileMoneyPayoutResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.paychangu.com/mobile-money/payouts/initialize", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.secretkey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // API typically returns 200 for successful initiation, 400 for errors
		bo, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var apiErr MobileMoneyPayoutErrorResponse
		if jsonErr := json.Unmarshal(bo, &apiErr); jsonErr == nil && apiErr.Status == "failed" {
			// For validation errors, the "message" field is a map
			var errorMessages []string
			if apiErr.Message != nil {
				for field, messages := range apiErr.Message {
					for _, msg := range messages {
						errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", field, msg))
					}
				}
			}
			if len(errorMessages) > 0 {
				return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errors.Join(errors.New("validation failed"), errors.New(string(bo))).Error())
			}
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bo))
		}

		// Fallback for other non-200 statuses or unexpected error formats
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bo))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response MobileMoneyPayoutResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, errors.New(response.Message)
	}

	return &response, nil
}

// GetMobileMoneyPayoutDetails retrieves the details of a specific mobile money payout.
//
// Parameters:
//
// chargeID (string): The unique charge ID of the mobile money payout to retrieve.
//
// Returns:
//
// *PayoutTransactionDetails: A pointer to a PayoutTransactionDetails struct containing the detailed information about the payout.
//
// error: An error, if one occurred during the request.
//
// Example Usage:
//
//	client := paychangu.New("your_secret_key")
//	chargeID := "MY_PAYOUT_TX_001" // The chargeID used when initiating the payout
//	payoutDetails, err := client.GetMobileMoneyPayoutDetails(chargeID)
//	if err != nil {
//	    log.Fatalf("Failed to get mobile money payout details: %v", err)
//	}
//	fmt.Printf("Payout Details for Charge ID %s: Status: %s, Amount: %.2f %s\n",
//	    payoutDetails.ChargeID, payoutDetails.Status, payoutDetails.Amount, payoutDetails.Currency)
func (p *payChangu) GetMobileMoneyPayoutDetails(chargeID string) (*PayoutTransactionDetails, error) {
	url := fmt.Sprintf("https://api.paychangu.com/mobile-money/payments/%sdetails", chargeID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.secretkey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bo, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var apiErr Error // Using the general Error struct for non-200 responses
		if jsonErr := json.Unmarshal(bo, &apiErr); jsonErr == nil && apiErr.Message != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, apiErr.Message)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bo))
	}

	var response GetMobileMoneyPayoutDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, errors.New(response.Message)
	}

	return &response.Data, nil
}

// GetSupportedBanks retrieves a list of banks supported for direct charge payouts for a given currency.
//
// Parameters:
//
// currency (string): The currency code for which to retrieve supported banks (e.g., "MWK", "USD").
//
// Returns:
//
// []SupportedBank: A slice of supported bank details.
//
// error: An error, if one occurred during the request.
//
// Example Usage:
//
//	client := paychangu.New("your_secret_key")
//	banks, err := client.GetSupportedBanks("MWK")
//	if err != nil {
//	    log.Fatalf("Failed to get supported banks: %v", err)
//	}
//	for _, bank := range banks {
//	    fmt.Printf("Bank: %s (UUID: %s)\n", bank.Name, bank.UUID)
//	}
func (p *payChangu) GetSupportedBanks(currency string) ([]Bank, error) {
	url := fmt.Sprintf("https://api.paychangu.com/direct-charge/payouts/supported-banks?currency=%s", currency)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.secretkey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bo, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var apiErr Error // Using the general Error struct for non-200 responses
		if jsonErr := json.Unmarshal(bo, &apiErr); jsonErr == nil && apiErr.Message != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, apiErr.Message)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bo))
	}

	var response BanksResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, errors.New(response.Message)
	}

	return response.Data, nil
}

// InitiateBankPayout sends a bank payout request to the PayChangu API.
//
// Parameters:
//
// request (BankPayoutRequest): The bank payout payload, including recipient bank details.
//
// Returns:
//
// *BankPayoutResponse: A pointer to a BankPayoutResponse struct containing details about the initiated bank payout.
//
// error: An error, if one occurred during the request. This can include detailed validation errors.
//
// Example Usage:
//
//	client := paychangu.New("your_secret_key")
//	bankPayoutReq := paychangu.BankPayoutRequest{
//	    PayoutMethod:      "bank_transfer", // Always "bank_transfer" for this method
//	    BankUUID:          "82310dd1-ec9b-4fe7-a32c-2f262ef08681", // Example NBM UUID from GetSupportedBanks
//	    Amount:            50000.00,
//	    ChargeID:          "BANK_PAYOUT_XYZ789",
//	    BankAccountName:   "John Doe",
//	    BankAccountNumber: "1000000010",
//	    Email:             "john.doe@example.com",
//	}
//	bankPayoutResp, err := client.InitiateBankPayout(bankPayoutReq)
//	if err != nil {
//	    log.Fatalf("Bank payout failed: %v", err)
//	}
//	fmt.Printf("Bank Payout Initiated. Charge ID: %s, Status: %s\n", bankPayoutResp.Data.Transaction.ChargeID, bankPayoutResp.Data.Transaction.Status)
//	fmt.Printf("Recipient Bank: %s, Account: %s\n", bankPayoutResp.Data.Transaction.RecipientAccountDetails.BankName, bankPayoutResp.Data.Transaction.RecipientAccountDetails.AccountNumber)
func (p *payChangu) InitiateBankPayout(request BankPayoutRequest) (*BankPayoutResponse, error) {
	// The API expects amount as a string, so we need to format it before marshaling
	// We'll create an anonymous struct to handle this, as modifying the original
	// BankPayoutRequest struct's Amount field to string would be less type-safe for users.
	requestPayload := struct {
		PayoutMethod      string `json:"payout_method"`
		BankUUID          string `json:"bank_uuid"`
		Amount            string `json:"amount"` // Marshaled as string
		ChargeID          string `json:"charge_id"`
		BankAccountName   string `json:"bank_account_name"`
		BankAccountNumber string `json:"bank_account_number"`
		Email             string `json:"email,omitempty"`
		FirstName         string `json:"first_name,omitempty"`
		LastName          string `json:"last_name,omitempty"`
	}{
		PayoutMethod:      request.PayoutMethod,
		BankUUID:          request.BankUUID,
		Amount:            fmt.Sprintf("%.2f", request.Amount), // Format float to string with 2 decimal places
		ChargeID:          request.ChargeID,
		BankAccountName:   request.BankAccountName,
		BankAccountNumber: request.BankAccountNumber,
		Email:             request.Email,
		FirstName:         request.FirstName,
		LastName:          request.LastName,
	}

	data, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.paychangu.com/direct-charge/payouts/initialize", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.secretkey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // Assuming 200 OK for success, and other codes for errors
		bo, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response body: %w", err)
		}

		// Try to unmarshal into MobileMoneyPayoutErrorResponse which handles map[string][]string for validation errors
		var apiErr MobileMoneyPayoutErrorResponse
		if jsonErr := json.Unmarshal(bo, &apiErr); jsonErr == nil && apiErr.Status == "failed" {
			// If message is a map (validation error), format it
			var errorMessages []string
			if apiErr.Message != nil {
				for field, messages := range apiErr.Message {
					for _, msg := range messages {
						errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", field, msg))
					}
				}
			}
			if len(errorMessages) > 0 {
				return nil, fmt.Errorf("API error (%d): validation failed: %s", resp.StatusCode, errors.Join(errors.New("validation failed"), errors.New(string(bo))).Error())
			}
		}

		// Fallback to general error struct or raw body if specific unmarshal fails
		var generalErr Error
		if jsonErr := json.Unmarshal(bo, &generalErr); jsonErr == nil && generalErr.Message != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, generalErr.Message)
		}

		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bo))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response BankPayoutResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, errors.New(response.Message)
	}

	return &response, nil
}

// GetBankPayoutDetails retrieves the details of a specific bank payout.
//
// Parameters:
//
// chargeID (string): The unique charge ID of the bank payout to retrieve.
//
// Returns:
//
// *BankPayoutTransactionDetails: A pointer to a BankPayoutTransactionDetails struct containing the detailed information about the payout.
//
// error: An error, if one occurred during the request.
//
// Example Usage:
//
//	client := paychangu.New("your_secret_key")
//	chargeID := "BANK_PAYOUT_XYZ789" // The chargeID used when initiating the bank payout
//	bankPayoutDetails, err := client.GetBankPayoutDetails(chargeID)
//	if err != nil {
//	    log.Fatalf("Failed to get bank payout details: %v", err)
//	}
//	fmt.Printf("Bank Payout Details for Charge ID %s: Status: %s, Amount: %.2f %s\n",
//	    bankPayoutDetails.ChargeID, bankPayoutDetails.Status, bankPayoutDetails.Amount, bankPayoutDetails.Currency)
func (p *payChangu) GetBankPayoutDetails(chargeID string) (*BankPayoutTransactionDetails, error) {
	url := fmt.Sprintf("https://api.paychangu.com/direct-charge/payouts/%s/details", chargeID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.secretkey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bo, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var apiErr Error // Using the general Error struct for non-200 responses
		if jsonErr := json.Unmarshal(bo, &apiErr); jsonErr == nil && apiErr.Message != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, apiErr.Message)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bo))
	}

	var response GetBankPayoutDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// The API returns "successful" instead of "success" for the status field in the top-level response.
	// We should check against both or just rely on HTTP status code if API behavior is consistent.
	// For robustness, checking the specific 'status' in the body is good.
	if response.Status != "successful" { // Note the 'successful' string
		return nil, errors.New(response.Message)
	}

	return &response.Data, nil
}
