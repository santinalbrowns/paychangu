# PayChangu Go SDK Documentation

The **PayChangu Go SDK** provides a robust and seamless integration with the [PayChangu platform](https://paychangu.readme.io/reference/welcome) for Go developers. It simplifies payment collection, transaction verification, and payout processing (mobile money and bank transfers) within Go applications.

---

## Key Features

- **Payment Collection**: Accept payments via mobile money or card.
- **Transaction Verification**: Confirm payment status with ease.
- **Mobile Money Payouts**: Disburse funds to mobile wallets.
- **Bank Payouts**: Transfer funds to bank accounts.
- **Operator & Bank Lookup**: Retrieve supported mobile money operators and banks.
- **Custom Metadata**: Attach transaction-specific metadata.

---

## Getting Started

### Prerequisites

- Go 1.16 or later
- PayChangu account and secret key (available from your [PayChangu dashboard](https://paychangu.readme.io))

### Installation

Install the SDK using:

```bash
go get github.com/santinalbrowns/paychangu
```

### Import the Package

```go
import "github.com/santinalbrowns/paychangu"
```

### Initialize the Client

You need your **PayChangu secret key**.

```go
client := paychangu.New("your_secret_key")
```

## Accepting Payments

### Prepare Payment Request

```go
request := paychangu.Request{
    Amount:    10500,
    Currency:  "MWK",
    FirstName: "John",
    LastName:  "Doe",
    Email:     "john@example.com",
    CallbackURL: "https://yourapp.com/success",
    ReturnURL:   "https://yourapp.com/failure",
    TxRef:       "TX-123456",
    Customization: struct {
        Title       string `json:"title"`
        Description string `json:"description"`
    }{
        Title:       "Order Payment",
        Description: "Payment for electronics",
    },
    Meta: struct {
        UUID     string `json:"uuid"`
        Response string `json:"response"`
    }{
        UUID:     "user-001",
        Response: "tracking-data",
    },
}
```

### Initiate the Payment

```go
response, err := client.InitiatePayment(request)
if err != nil {
    log.Fatalf("Payment initiation failed: %v", err)
}
fmt.Println("Checkout URL:", response.Data.CheckoutURL)
```

### Field Descriptions

| Field                       | Required | Description                           |
| --------------------------- | -------- | ------------------------------------- |
| `Amount`                    | Yes      | Amount to charge in selected currency |
| `Currency`                  | Yes      | Currency code (e.g., "MWK", "USD")    |
| `FirstName`                 | Yes      | Customer’s first name                 |
| `LastName`                  | No       | Customer’s last name                  |
| `Email`                     | No       | Customer’s email                      |
| `CallbackURL`               | Yes      | Redirect URL after success            |
| `ReturnURL`                 | Yes      | Redirect URL after failure/cancel     |
| `TxRef`                     | Yes      | Unique transaction reference          |
| `Customization.title`       | Yes      | Title on checkout screen              |
| `Customization.description` | Yes      | Description on checkout screen        |
| `Meta`                      | No       | Extra metadata (e.g., user ID)        |

## Verifying a Payment

Use this step to confirm if a payment was successful.

```go
verification, err := client.VerifyPayment("TX-123456")
if err != nil {
    log.Fatalf("Verification failed: %v", err)
}
fmt.Println("Payment Status:", verification.Data.Status)
```

## Mobile Money Payouts

### Step 1: Fetch Mobile Money Operators

```go
operators, err := client.GetMobileMoneyOperators()
if err != nil {
    log.Fatalf("Operator fetch failed: %v", err)
}
for _, op := range operators {
    fmt.Println("Name:", op.Name)
    fmt.Println("RefID:", op.RefID) // Needed for payouts
}
```

### Step 2: Initiate Mobile Money Payout

```go
request := paychangu.MobileMoneyPayoutRequest{
    Mobile: "0881234567",
    Amount: 5000,
    MobileMoneyOperatorRefID: "27494cb5-ba9e-437f-a114-4e7a7686bcca",
    ChargeID: fmt.Sprintf("PAYOUT-%d", time.Now().UnixNano()),
    Email: "jane@example.com",
    FirstName: "Jane",
    LastName: "Doe",
}

response, err := client.InitiateMobileMoneyPayout(request)
if err != nil {
    panic(err)
}
fmt.Println("Status:", response.Status)
fmt.Println("Charge ID:", response.Data.Transaction.ChargeID)
```

#### Field Descriptions

| Field                            | Required | Description                       |
| -------------------------------- | -------- | --------------------------------- |
| `Mobile`                         | Yes      | Recipient phone number            |
| `Amount`                         | Yes      | Amount to send                    |
| `MobileMoneyOperatorRefID`       | Yes      | RefID from operator list          |
| `ChargeID`                       | Yes      | Unique identifier for this payout |
| `Email`, `FirstName`, `LastName` | No       | Optional recipient info           |
| `TransactionStatus`              | No       | For sandbox test statuses         |

### Fetch Payout Details

```go
details, err := client.GetMobileMoneyPayoutDetails("PAYOUT-123456")
fmt.Println("Status:", details.Status)
fmt.Println("Amount:", details.Amount)
```

## Bank Payouts

### Step 1: Fetch Supported Banks

```go
banks, err := client.GetSupportedBanks("MWK")
for _, b := range banks {
    fmt.Printf("Bank: %s, UUID: %s\n", b.Name, b.UUID)
}
```

### Step 2: Initiate Bank Payout

```go
bankPayout := paychangu.BankPayoutRequest{
    PayoutMethod: "bank_transfer",
    BankUUID: "82310dd1-ec9b-4fe7-a32c-2f262ef08681",
    Amount: 50000,
    ChargeID: "BANK_PAYOUT_XYZ789",
    BankAccountName: "John Doe",
    BankAccountNumber: "1000000010",
    Email: "john.doe@example.com",
}

resp, err := client.InitiateBankPayout(bankPayout)
fmt.Println("Charge ID:", resp.Data.Transaction.ChargeID)
fmt.Println("Status:", resp.Data.Transaction.Status)
```

### Fetch Bank Payout Details

```go
bankDetails, err := client.GetBankPayoutDetails("BANK_PAYOUT_XYZ789")
fmt.Println("Status:", bankDetails.Status)
fmt.Println("Bank:", bankDetails.RecipientAccountDetails.BankName)
```

## Project Structure

```bash
project/
├── main.go
├── go.mod
└── README.md
```

## Error Handling

Errors returned include detailed messages for debugging:

```go
_, err := client.InitiatePayment(req)
if err != nil {
    fmt.Println("Something went wrong:", err)
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request to improve the SDK
