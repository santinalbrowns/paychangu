# PayChangu Go SDK

The **[PayChangu](https://paychangu.readme.io/reference/welcome) Go SDK** is a library designed to make it simple for developers to initiate and verify payments using the [PayChangu](https://paychangu.readme.io/reference/welcome) payment API. With this SDK, you can charge customers, receive payment notifications, and verify payment transactions.

## Features

- Initiate Payments: Start a payment by specifying the amount, currency, customer details, and more.
- Verify Payments: Confirm the status of a payment by transaction reference.
- Customizable Metadata: Add extra information to transactions via the Meta field.

## Installation

To install the PayChangu SDK, use the following command:

```bash
go get github.com/santinalbrowns/paychangu
```

## Getting Started

1. Import the Package
In your Go code, import the PayChangu package:

    ```go
    import "github.com/santinalbrowns/paychangu"
    ```

2. Initialize the Client
To use the SDK, create a new PayChangu client by providing your secret API key.

    ```go
    client := paychangu.New("your_secret_key")
    ```

3. To initiate a payment, you need to create a PayChanguRequest with the required details, such as Amount, Currency, TxRef, and customer details.

    ```go
    request := paychangu.PayChanguRequest{
        Amount:    10500,
        Currency:  "MWK",
        FirstName: "John",
        LastName:  "Doe",
        Email:     "<johndoe@example.com>",
        CallbackURL: "<https://yourcallback.url>",
        ReturnURL:   "<https://yourreturn.url>",
        TxRef:       "unique_transaction_reference",
        Customization: struct {
            Title       string `json:"title"`
            Description string `json:"description"`
        }{
            Title:       "Service Payment",
            Description: "Payment for services rendered",
        },
        Meta: struct {
            UUID     string `json:"uuid"`
            Response string `json:"response"`
        }{
            UUID:     "unique_user_identifier",
            Response: "custom_response_data",
        },
    }

    // Initiate the payment
    response, err := client.InitiatePayment(request)
    if err != nil {
        log.Fatalf("Error initiating payment: %v", err)
    }

    fmt.Printf("Payment Initiated. Checkout URL: %s\n", response.Data.CheckoutURL)
    ```

    Fields:

    - Amount (required): Amount to charge the customer.
    - Currency (required): Currency code (e.g., "USD" or "MWK").
    - FirstName (required): Customer's first name.
    - CallbackURL (required): URL to which PayChangu will redirect after payment success.
    - ReturnURL (required): URL to which PayChangu will redirect after payment failure or cancellation.
    - TxRef (required): Unique transaction reference.

    The response will contain a CheckoutURL where the customer can complete the payment.

4. Verify a Payment
Once the payment process is complete, you can verify the status using the transaction reference (TxRef).

    ```go
    txRef := "unique_transaction_reference"

    verificationResponse, err := client.VerifyPayment(txRef)
    if err != nil {
        log.Fatalf("Error verifying payment: %v", err)
    }

    fmt.Printf("Payment Status: %s\n", verificationResponse.Data.Status)
    ````

## Example Project Structure

```bash
project/
│
├── main.go     # Your main application file
└── go.mod      # Go module file with paychangu dependency
```

### Example main.go

```go
package main

import (
    "fmt"
    "log"
    "github.com/santinalbrowns/paychangu"
)

func main() {
    client := paychangu.New("your_secret_key")

    request := paychangu.PayChanguRequest{
        Amount:    10500,
        Currency:  "MWK",
        FirstName: "John",
        LastName:  "Doe",
        Email:     "johndoe@example.com",
        CallbackURL: "https://yourcallback.url/checkout",
        ReturnURL:   "https://yourreturn.url/cancel",
        TxRef:       "unique_transaction_reference",
    }

    response, err := client.InitiatePayment(request)
    if err != nil {
        log.Fatalf("Error initiating payment: %v", err)
    }

    fmt.Printf("Payment Initiated. Checkout URL: %s\n", response.Data.CheckoutURL)
}
```

## Error Handling

The library returns detailed error messages for any failed request. Use these error messages to debug issues or display error messages to users.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request to improve the SDK
