# PayChangu Go SDK

The **[PayChangu](https://paychangu.readme.io/reference/welcome) Go SDK** is a library designed to make it simple for developers to initiate and verify payments using the [PayChangu](https://paychangu.readme.io/reference/welcome) payment API. With this SDK, you can charge customers, receive payment notifications, and verify payment transactions.

## Features

- **Initiate Payments**: Start a payment by specifying the amount, currency, customer details, and more.
- **Verify Payments**: Confirm the status of a payment by transaction reference.
- **Mobile Money Payouts**: Disburse funds directly to mobile money wallets.
    - Get Supported Mobile Money Operators: Retrieve a list of available mobile money networks.
    - Initiate Mobile Money Payout: Send funds to a mobile money number.
- **Customizable Metadata**: Add extra information to transactions via the Meta field.

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

3. To initiate a payment, you need to create a Request with the required details, such as Amount, Currency, TxRef, and customer details.

    ```go
    request := paychangu.Request{
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

    request := paychangu.Request{
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

    // Initiate Mobile Money Payout Example 
    // You'll need a valid mobile number 
    // and operator ref_id for a real payout.
    // Using dummy values for example purposes.
 
    payoutRequest := paychangu.MobileMoneyPayoutRequest{
        Mobile:                    "+265888123456",
        MobileMoneyOperatorRefID:  "27494cb5-ba9e-437f-a114-4e7a7686bcca",
        Amount:                    100.00,
        ChargeID:                  "payout_ref_" + fmt.Sprintf("%d", time.Now().Unix()),
        Email:                     "recipient@example.com",
        FirstName:                 "Test",
        LastName:                  "User",
    }

    payoutResponse, err := client.InitiateMobileMoneyPayout(payoutRequest)
    if err != nil {
        if payoutErr, ok := err.(*paychangu.MobileMoneyPayoutErrorResponse); ok {
             fmt.Printf("Payout validation errors: %+v\n", payoutErr.Message)
        }
    } else {
        fmt.Printf("Transaction Ref ID: %s\n", payoutResponse.Data.Transaction.RefID)
        fmt.Printf("Transaction Status: %s\n", payoutResponse.Data.Transaction.Status)
    }
}
```

## Payouts

### Mobile Money Payout (MOMO)

1. Get Supported Mobile Money Operators
Before initiating a mobile money payout, you can retrieve a list of supported operators and their details.

    ```go
    operators, err := client.GetMobileMoneyOperators()
    if err != nil {
        log.Fatalf("Error getting mobile money operators: %v", err)
    }

    fmt.Println("Supported Mobile Money Operators:")
    for _, op := range operators {
        fmt.Printf("- %s (Ref ID: %s, Supports Withdrawals: %t)\n", op.Name, op.RefID, op.SupportsWithdrawals)
    }
    ```

2. Initiate a Mobile Money Payout
To send money to a mobile money wallet, you'll need the recipient's mobile number, the operator's ref_id (obtained from GetMobileMoneyOperators), the amount, and a unique charge_id.

    ```go
    payoutRequest := paychangu.MobileMoneyPayoutRequest{
        Mobile:                    "+265888123456",
        MobileMoneyOperatorRefID:  "27494cb5-ba9e-437f-a114-4e7a7686bcca",
        Amount:                    5000.00,
        ChargeID:                  "MY_PAYOUT_TX_001",
        Email:                     "recipient.email@example.com",
        FirstName:                 "Recipient",
        LastName:                  "User",
    }

    payoutResponse, err := client.InitiateMobileMoneyPayou(payoutRequest)
    if err != nil {
        log.Fatalf("Error initiating mobile money payout: %v", err)
    }

    fmt.Printf("Mobile Money Payout initiated successfully!\n")
    fmt.Printf("Transaction Ref ID: %s\n", payoutResponse.Data.Transaction.RefID)
    fmt.Printf("Transaction Status: %s\n", payoutResponse.Data.Transaction.Status)
    ```

## Error Handling

The library returns detailed error messages for any failed request. Use these error messages to debug issues or display error messages to users.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request to improve the SDK
