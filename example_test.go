package paychangu_test

import (
	"fmt"

	"github.com/santinalbrowns/paychangu"
)

func ExampleNew() {
	client := paychangu.New("your_secret_key")
	_ = client
}

func ExamplePayChangu_InitiatePayment() {
	client := paychangu.New("your_secret_key")

	resp, err := client.InitiatePayment(paychangu.Request{
		Amount:      10500,
		Currency:    "MWK",
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john@example.com",
		CallbackURL: "https://yourapp.com/success",
		ReturnURL:   "https://yourapp.com/cancel",
		TxRef:       "TX-123",
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(resp.Data.CheckoutURL)
}
