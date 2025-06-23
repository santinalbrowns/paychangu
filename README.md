# PayChangu Go SDK

The **[PayChangu](https://paychangu.readme.io/reference/welcome) Go SDK** helps developers integrate with the [PayChangu payment platform](https://paychangu.readme.io/reference/welcome). It supports collecting payments from customers and sending payouts via mobile money. The SDK provides a simplified interface to:

- Accept payments from customers via mobile money or card.
- Verify transaction status.
- Send payouts to customers’ mobile wallets.
- Fetch available mobile money operators.
- Customize payment experiences with metadata and branding.

---

## 🔧 Features

- ✅ **Initiate Payments**: Charge a customer via card or mobile money.
- ✅ **Verify Payments**: Confirm transaction status by reference.
- ✅ **Mobile Money Payouts**: Send money to mobile wallets.
- ✅ **Fetch Operators**: List available mobile money operators.
- ✅ **Custom Metadata**: Attach custom identifiers (UUID, user data, etc.).
- ✅ **Branded Checkout**: Show title and description on checkout page.

---

## 🚀 Installation

To install the SDK in your Go project:

```bash
go get github.com/santinalbrowns/paychangu
```

---

## 📦 Getting Started

### 1. Import the Package

```go
import "github.com/santinalbrowns/paychangu"
```

### 2. Initialize the PayChangu Client

You need your **secret key** to authenticate with PayChangu:

```go
client := paychangu.New("your_secret_key")
```

---

## 💰 Accepting Payments

### Step-by-Step Example

```go
request := paychangu.Request{
    Amount:    10500,
    Currency:  "MWK",
    FirstName: "John",
    LastName:  "Doe",
    Email:     "john@example.com",
    CallbackURL: "https://yourapp.com/payment/success",
    ReturnURL:   "https://yourapp.com/payment/failed",
    TxRef:       "TX-12345-ABC",

    Customization: struct {
        Title       string `json:"title"`
        Description string `json:"description"`
    }{
        Title:       "Premium Subscription",
        Description: "Access to all features",
    },

    Meta: struct {
        UUID     string `json:"uuid"`
        Response string `json:"response"`
    }{
        UUID:     "user-abc-001",
        Response: "subscription-payment",
    },
}

response, err := client.InitiatePayment(request)
if err != nil {
    log.Fatalf("Payment error: %v", err)
}

fmt.Println("Checkout URL:", response.Data.CheckoutURL)
```

### 🔍 Fields Explained

| Field         | Required | Description |
|---------------|----------|-------------|
| `Amount`      | ✅        | Amount to charge (in the specified currency) |
| `Currency`    | ✅        | `"MWK"` or `"USD"` |
| `FirstName`   | ✅        | Customer’s first name |
| `LastName`    | ❌        | Customer’s last name |
| `Email`       | ❌        | Customer’s email for receipts |
| `CallbackURL` | ✅        | Redirect URL after successful payment |
| `ReturnURL`   | ✅        | Redirect URL if payment fails or is cancelled |
| `TxRef`       | ✅        | Unique transaction reference (per payment) |
| `Customization.title` | ✅ | Payment title shown on checkout |
| `Customization.description` | ✅ | Payment description |
| `Meta`        | ❌        | Optional extra info (e.g., user ID, reference data) |

---

## ✅ Verifying a Payment

After a customer completes (or cancels) payment, verify its status with the transaction reference:

```go
verification, err := client.VerifyPayment("TX-12345-ABC")
if err != nil {
    log.Fatalf("Verification error: %v", err)
}

fmt.Println("Status:", verification.Data.Status)
fmt.Println("Amount:", verification.Data.Amount)
fmt.Println("Customer:", verification.Data.Customer.Email)
```

---

## 💸 Mobile Money Payouts

You can also send money to customers’ mobile wallets using their number and mobile money provider.

---

### Step 1: Get List of Operators

```go
operators, err := client.GetMobileMoneyOperators()
if err != nil {
    panic(err)
}

for _, op := range operators {
    fmt.Println("Name:", op.Name)
    fmt.Println("RefID:", op.RefID) // Needed for payouts
}
```

---

### Step 2: Send a Payout

```go
request := paychangu.MobileMoneyPayoutRequest{
    Mobile:                   "0881234567",
    Amount:                   5000,
    MobileMoneyOperatorRefID: "27494cb5-ba9e-437f-a114-4e7a7686bcca", // Use RefID from step 1
    ChargeID:                 fmt.Sprintf("PAYOUT-%d", time.Now().UnixNano()),
    Email:                    "jane@example.com",
    FirstName:                "Jane",
    LastName:                 "Doe",
}

response, err := client.InitiateMobileMoneyPayout(request)
if err != nil {
    panic(err)
}

fmt.Println("Payout sent. Status:", response.Status)
fmt.Println("Transaction ID:", response.Data.Transaction.ChargeID)
```

### 📌 Payout Fields Explained

| Field                      | Required | Description |
|---------------------------|----------|-------------|
| `Mobile`                  | ✅        | Recipient’s phone number |
| `Amount`                  | ✅        | Amount to send |
| `MobileMoneyOperatorRefID`| ✅        | RefID from available operators |
| `ChargeID`                | ✅        | Unique identifier for this payout |
| `Email`                   | ❌        | Optional recipient email |
| `FirstName`               | ❌        | Optional first name |
| `LastName`                | ❌        | Optional last name |
| `TransactionStatus`       | ❌        | For mocking responses in sandbox |

---

## 🧪 Example Project Layout

```bash
your-project/
│
├── main.go        # Your app entry point
├── go.mod         # Go module config
└── README.md      # This documentation
```

---

## 🧯 Error Handling

If a request fails, the SDK returns a descriptive error message. Handle it like this:

```go
_, err := client.InitiatePayment(req)
if err != nil {
    fmt.Println("Something went wrong:", err)
}
```

For failed payouts, you can inspect the response for detailed error messages.

---

## 🤝 Contributing

Want to improve the SDK? Fix a bug? Add a feature?  
Please open a pull request or issue at [GitHub Repo](https://github.com/santinalbrowns/paychangu).

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).
