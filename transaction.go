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

		return nil, fmt.Errorf(string(bo))
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
