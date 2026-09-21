# PayChangu Go SDK

Go client for the [PayChangu](https://developer.paychangu.com) payment APIs: checkout, direct charge, card, payouts, bills, Connect, and US virtual accounts.

## Requirements

- Go 1.23+
- PayChangu secret key from your [dashboard](https://developer.paychangu.com)

## Install

```bash
go get github.com/mzati-paychangu/paychangu_go_sdk
```

```go
import "github.com/mzati-paychangu/paychangu_go_sdk"
```

## Client

```go
client := paychangu.New("your_secret_key")
```

Optional: `paychangu.WithHTTPClient(...)`, `paychangu.WithBaseURL(...)`.

## Features

- Standard checkout (`InitiatePayment`, `VerifyPayment`)
- Wallet balance (`GetBalance`)
- Direct charge: mobile money, bank transfer, verify/details
- Card charge, verify, refund
- Mobile money & bank payouts (including list bank payouts)
- Billers, validate/pay bill, airtime, bill stats
- PayChangu Connect (authorize link, user, revoke)
- US virtual accounts (customers + account lifecycle)

## Standard checkout

```go
resp, err := client.InitiatePayment(paychangu.Request{
    Amount: 10500, Currency: "MWK", FirstName: "John", LastName: "Doe",
    Email: "john@example.com",
    CallbackURL: "https://yourapp.com/success",
    ReturnURL:   "https://yourapp.com/failure",
    TxRef:       "TX-123456",
})
checkoutURL := resp.Data.CheckoutURL

verification, err := client.VerifyPayment("TX-123456")
```

## Direct charge (mobile money)

```go
charge, err := client.ChargeMobileMoney(paychangu.ChargeMobileMoneyRequest{
    Mobile: "265888123456",
    MobileMoneyOperatorRefID: "20be6c20-adeb-4b5b-a7ba-0769820df4fb",
    Amount: "1000",
    ChargeID: "CHARGE-001",
})
status, err := client.VerifyDirectCharge("CHARGE-001")
details, err := client.GetChargeDetails("CHARGE-001")
```

## Bank transfer collection

```go
transfer, err := client.InitiateBankTransfer(paychangu.BankTransferRequest{
    Amount: "1000", Currency: "MWK", ChargeID: "PTC-001",
})
account := transfer.Data.PaymentAccountDetails
tx, err := client.GetBankTransferDetails("PTC-001")
```

## Card

```go
card, err := client.ChargeCard(paychangu.ChargeCardRequest{
    CardNumber: "4242424242424242", Expiry: "12/30", CVV: "123",
    CardholderName: "John Doe", Amount: "1000", Currency: "USD",
    ChargeID: "CARD-001", RedirectURL: "https://yourapp.com/redirect",
})
_, err = client.VerifyCardCharge("CARD-001")
_, err = client.RefundCardCharge("CARD-001")
```

## Payouts

```go
operators, err := client.GetMobileMoneyOperators()
payout, err := client.InitiateMobileMoneyPayout(paychangu.MobileMoneyPayoutRequest{
    Mobile: "0881234567", Amount: 5000,
    MobileMoneyOperatorRefID: operators[0].RefID,
    ChargeID: "PAYOUT-001",
})
details, err := client.GetMobileMoneyPayoutDetails("PAYOUT-001")

banks, err := client.GetSupportedBanks("MWK")
bankPayout, err := client.InitiateBankPayout(paychangu.BankPayoutRequest{
    PayoutMethod: "bank_transfer", BankUUID: banks[0].UUID,
    Amount: 50000, ChargeID: "BANK-001",
    BankAccountName: "John Doe", BankAccountNumber: "1000000010",
})
bankDetails, err := client.GetBankPayoutDetails("BANK-001")
all, err := client.ListBankPayouts()
```

## Bills

```go
billers, err := client.GetBillers()
_, err = client.ValidateBill(paychangu.ValidateBillRequest{Biller: "escom", Account: "123"})
_, err = client.PayBill(paychangu.PayBillRequest{Biller: "escom", Account: "123", Amount: "5000"})
_, err = client.BuyAirtime(paychangu.AirtimeRequest{Phone: "0999123456", Amount: "1000"})
```

## Connect

```go
link, err := client.CreateConnectLink(paychangu.ConnectLinkRequest{
    ClientID: "your-client-id",
    RedirectURI: "https://yourapp.com/connect/callback",
    Mode: "live",
    Scope: "payments:write payments:read",
})
user, err := client.GetConnectUser(accessToken)
err = client.RevokeAccessToken(accessToken) // uses /connect/revoke
```

## US virtual accounts

```go
customer, err := client.CreateCustomer(paychangu.VirtualCustomerRequest{
    Email: "john@example.com", FirstName: "John", LastName: "Doe",
})
account, err := client.CreateUSAccount(customer.Data.ID)
_, err = client.GetUSAccountActivity(customer.Data.ID)
```

## Balance

```go
balance, err := client.GetBalance("MWK")
```

## Docs

API reference: https://developer.paychangu.com/llms.txt

## Contributing

Open an issue or PR against this repository.
