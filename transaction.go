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
