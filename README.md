# PayChangu Go SDK

Go client for [PayChangu](https://developer.paychangu.com).

## Install

```bash
go get github.com/santinalbrowns/paychangu
```

```go
import "github.com/santinalbrowns/paychangu"
```

## Quick start

```go
client := paychangu.New("your_secret_key")

// Checkout
resp, err := client.InitiatePayment(paychangu.Request{
    Amount: 10500, Currency: "MWK",
    FirstName: "John", LastName: "Doe",
    Email: "john@example.com",
    CallbackURL: "https://yourapp.com/success",
    ReturnURL:   "https://yourapp.com/cancel",
    TxRef:       "TX-123",
})
fmt.Println(resp.Data.CheckoutURL)

status, err := client.VerifyPayment("TX-123")
```

## Common APIs

**Balance**
```go
bal, err := client.GetBalance("MWK")
```

**Direct MoMo charge**
```go
charge, err := client.ChargeMobileMoney(paychangu.ChargeMobileMoneyRequest{
    Mobile: "265888123456",
    MobileMoneyOperatorRefID: "operator-ref-id",
    Amount: "1000", ChargeID: "CHARGE-1",
})
_, err = client.VerifyDirectCharge("CHARGE-1")
```

**Payouts**
```go
ops, err := client.GetMobileMoneyOperators()
_, err = client.InitiateMobileMoneyPayout(paychangu.MobileMoneyPayoutRequest{
    Mobile: "0881234567", Amount: 5000,
    MobileMoneyOperatorRefID: ops[0].RefID,
    ChargeID: "PAYOUT-1",
})

banks, err := client.GetSupportedBanks("MWK")
_, err = client.InitiateBankPayout(paychangu.BankPayoutRequest{
    PayoutMethod: "bank_transfer", BankUUID: banks[0].UUID,
    Amount: 50000, ChargeID: "BANK-1",
    BankAccountName: "John Doe", BankAccountNumber: "1000000010",
})
```

**Card / bills / Connect / virtual accounts** — see Go docs on the exported methods (`ChargeCard`, `PayBill`, `CreateConnectLink`, `CreateCustomer`, …).

## Options

```go
client := paychangu.New("key",
    paychangu.WithHTTPClient(myHTTPClient),
    paychangu.WithBaseURL("https://api.paychangu.com"),
)
```

## Layout

Single Go package at the module root (`package paychangu`). Split across files by domain (`payment.go`, `payout.go`, `card.go`, …) — import stays one path.

## Full API docs

https://developer.paychangu.com/llms.txt
