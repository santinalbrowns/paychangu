package paychangu

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (*PayChangu, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	client := New("test-secret", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	return client, server
}

func assertAuth(t *testing.T, r *http.Request) {
	t.Helper()
	if got := r.Header.Get("Authorization"); got != "Bearer test-secret" {
		t.Fatalf("Authorization = %q, want Bearer test-secret", got)
	}
}

func TestInitiatePayment(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuth(t, r)
		if r.Method != http.MethodPost || r.URL.Path != "/payment" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatal(err)
		}
		if req.TxRef != "TX-1" || req.Amount != 100 {
			t.Fatalf("unexpected body: %+v", req)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"event":"checkout","checkout_url":"https://pay.example/c","data":{"tx_ref":"TX-1","currency":"MWK","amount":100,"mode":"live","status":"pending"}}}`))
	})
	defer server.Close()

	resp, err := client.InitiatePayment(Request{Amount: 100, Currency: "MWK", TxRef: "TX-1", FirstName: "Ada"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data.CheckoutURL != "https://pay.example/c" {
		t.Fatalf("checkout url = %q", resp.Data.CheckoutURL)
	}
}

func TestVerifyPayment(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuth(t, r)
		if r.URL.Path != "/verify-payment/TX-1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"tx_ref":"TX-1","status":"success","amount":100,"currency":"MWK"}}`))
	})
	defer server.Close()

	resp, err := client.VerifyPayment("TX-1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data.TxRef != "TX-1" {
		t.Fatalf("tx_ref = %q", resp.Data.TxRef)
	}
}

func TestGetMobileMoneyOperators(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuth(t, r)
		if r.URL.Path != "/mobile-money" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":[{"id":1,"name":"Airtel Money","ref_id":"abc","short_code":"airtel"}]}`))
	})
	defer server.Close()

	ops, err := client.GetMobileMoneyOperators()
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 || ops[0].Name != "Airtel Money" {
		t.Fatalf("ops = %+v", ops)
	}
}

func TestGetMobileMoneyPayoutDetailsPath(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuth(t, r)
		if r.URL.Path != "/mobile-money/payments/charge-1/details" {
			t.Fatalf("path = %s, want slash before details", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"charge_id":"charge-1","status":"success","amount":50,"currency":"MK"}}`))
	})
	defer server.Close()

	details, err := client.GetMobileMoneyPayoutDetails("charge-1")
	if err != nil {
		t.Fatal(err)
	}
	if details.ChargeID != "charge-1" {
		t.Fatalf("charge_id = %q", details.ChargeID)
	}
}

func TestInitiateMobileMoneyPayout(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mobile-money/payouts/initialize" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"transaction":{"charge_id":"p1","status":"success","amount":100}}}`))
	})
	defer server.Close()

	resp, err := client.InitiateMobileMoneyPayout(MobileMoneyPayoutRequest{
		Mobile: "265999", MobileMoneyOperatorRefID: "op", Amount: 100, ChargeID: "p1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data.Transaction.ChargeID != "p1" {
		t.Fatalf("charge_id = %q", resp.Data.Transaction.ChargeID)
	}
}

func TestGetSupportedBanksEscapesCurrency(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/direct-charge/payouts/supported-banks" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("currency") != "MWK" {
			t.Fatalf("currency = %q", r.URL.Query().Get("currency"))
		}
		_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":[{"uuid":"u1","name":"NBM"}]}`))
	})
	defer server.Close()

	banks, err := client.GetSupportedBanks("MWK")
	if err != nil {
		t.Fatal(err)
	}
	if len(banks) != 1 || banks[0].UUID != "u1" {
		t.Fatalf("banks = %+v", banks)
	}
}

func TestInitiateBankPayoutAmountString(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}
		amount, ok := payload["amount"].(string)
		if !ok || amount != "1000.50" {
			t.Fatalf("amount = %#v", payload["amount"])
		}
		_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"transaction":{"charge_id":"b1","status":"success","amount":1000.5}}}`))
	})
	defer server.Close()

	resp, err := client.InitiateBankPayout(BankPayoutRequest{
		PayoutMethod: "bank_transfer", BankUUID: "u", Amount: 1000.5, ChargeID: "b1",
		BankAccountName: "Ada", BankAccountNumber: "123",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data.Transaction.ChargeID != "b1" {
		t.Fatalf("charge_id = %q", resp.Data.Transaction.ChargeID)
	}
}

func TestGetBankPayoutDetailsSuccessfulStatus(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/direct-charge/payouts/b1/details" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"status":"successful","message":"ok","data":{"charge_id":"b1","status":"success","amount":10}}`))
	})
	defer server.Close()

	details, err := client.GetBankPayoutDetails("b1")
	if err != nil {
		t.Fatal(err)
	}
	if details.ChargeID != "b1" {
		t.Fatalf("charge_id = %q", details.ChargeID)
	}
}

func TestListBankPayouts(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/direct-charge/payouts" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"current_page":1,"total_pages":1,"per_page":20,"data":[{"charge_id":"x"}]}}`))
	})
	defer server.Close()

	page, err := client.ListBankPayouts()
	if err != nil {
		t.Fatal(err)
	}
	if page.CurrentPage != 1 || len(page.Data) != 1 {
		t.Fatalf("page = %+v", page)
	}
}

func TestGetBalance(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wallet-balance" || r.URL.Query().Get("currency") != "USD" {
			t.Fatalf("bad request %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"environment":"live","currency":"USD","main_balance":"12.00","collection_balance":0}}`))
	})
	defer server.Close()

	bal, err := client.GetBalance("USD")
	if err != nil {
		t.Fatal(err)
	}
	if bal.MainBalance != "12.00" {
		t.Fatalf("balance = %+v", bal)
	}
}

func TestChargeMobileMoneyAndVerify(t *testing.T) {
	var sawCharge, sawVerify, sawDetails bool
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/mobile-money/payments/initialize":
			sawCharge = true
			_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"charge_id":"c1","status":"pending","amount":50}}`))
		case r.URL.Path == "/mobile-money/payments/c1/verify":
			sawVerify = true
			_, _ = w.Write([]byte(`{"status":"successful","message":"ok","data":{"charge_id":"c1","status":"success","amount":50}}`))
		case r.URL.Path == "/mobile-money/payments/c1/details":
			sawDetails = true
			_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"charge_id":"c1","status":"success","amount":50}}`))
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	})
	defer server.Close()

	charged, err := client.ChargeMobileMoney(ChargeMobileMoneyRequest{
		Mobile: "265", MobileMoneyOperatorRefID: "op", Amount: "50", ChargeID: "c1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if charged.ChargeID != "c1" {
		t.Fatal(charged)
	}
	verified, err := client.VerifyDirectCharge("c1")
	if err != nil {
		t.Fatal(err)
	}
	if verified.Status != "success" {
		t.Fatal(verified)
	}
	details, err := client.GetChargeDetails("c1")
	if err != nil {
		t.Fatal(err)
	}
	if details.ChargeID != "c1" || !sawCharge || !sawVerify || !sawDetails {
		t.Fatalf("details=%+v flags=%v/%v/%v", details, sawCharge, sawVerify, sawDetails)
	}
}

func TestBankTransfer(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/direct-charge/payments/initialize":
			body, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(body), `"payment_method":"mobile_bank_transfer"`) {
				t.Fatalf("body = %s", body)
			}
			_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"payment_account_details":{"bank_name":"Centenary","account_number":"1","account_name":"Pay","account_expiration_timestamp":1},"transaction":{"charge_id":"PTC1","status":"pending","amount":1000}}}`))
		case "/direct-charge/transactions/PTC1/details":
			_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"transaction":{"charge_id":"PTC1","status":"success","amount":1000}}}`))
		default:
			t.Fatalf("path = %s", r.URL.Path)
		}
	})
	defer server.Close()

	resp, err := client.InitiateBankTransfer(BankTransferRequest{Amount: "1000", Currency: "MWK", ChargeID: "PTC1"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data.PaymentAccountDetails.BankName != "Centenary" {
		t.Fatal(resp)
	}
	tx, err := client.GetBankTransferDetails("PTC1")
	if err != nil {
		t.Fatal(err)
	}
	if tx.ChargeID != "PTC1" {
		t.Fatal(tx)
	}
}

func TestCardAPIs(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/charge-card/payments":
			_, _ = w.Write([]byte(`{"success":true,"requires_3ds_auth":true,"orderReference":"ord","3ds_auth_link":"https://3ds"}`))
		case r.URL.Path == "/charge-card/verify/c1":
			_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"status":"success"}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/charge-card/refund/c1":
			_, _ = w.Write([]byte(`{"status":"success","message":"refunded","data":{}}`))
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	})
	defer server.Close()

	charged, err := client.ChargeCard(ChargeCardRequest{
		CardNumber: "4242", Expiry: "12/30", CVV: "123", CardholderName: "Ada",
		Amount: "100", Currency: "USD", ChargeID: "c1", RedirectURL: "https://ex",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !charged.Requires3DSAuth || charged.ThreeDSAuthLink == "" {
		t.Fatal(charged)
	}
	if _, err := client.VerifyCardCharge("c1"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.RefundCardCharge("c1"); err != nil {
		t.Fatal(err)
	}
}

func TestBillsAPIs(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"success","message":"ok","data":{"ok":true}}`))
	})
	defer server.Close()

	checks := []struct {
		name string
		fn   func() error
	}{
		{"GetBillers", func() error { _, err := client.GetBillers(); return err }},
		{"GetBillerDetails", func() error { _, err := client.GetBillerDetails("escom"); return err }},
		{"ValidateBill", func() error {
			_, err := client.ValidateBill(ValidateBillRequest{Biller: "escom", Account: "1"})
			return err
		}},
		{"PayBill", func() error {
			_, err := client.PayBill(PayBillRequest{Biller: "escom", Account: "1", Amount: "100"})
			return err
		}},
		{"BuyAirtime", func() error {
			_, err := client.BuyAirtime(AirtimeRequest{Phone: "0999", Amount: "100"})
			return err
		}},
		{"GetBillTransaction", func() error { _, err := client.GetBillTransaction("ref-1"); return err }},
		{"GetBillStatistics", func() error { _, err := client.GetBillStatistics(); return err }},
	}
	for _, c := range checks {
		if err := c.fn(); err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
	}
}

func TestConnectAPIs(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/connect/authorize-link":
			if r.URL.Query().Get("client_id") != "cid" || r.URL.Query().Get("mode") != "live" {
				t.Fatalf("query = %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"status":"success","data":{"url":"https://connect"}}`))
		case r.URL.Path == "/connect/user":
			if r.URL.Query().Get("access_token") != "tok" {
				t.Fatalf("token missing")
			}
			_, _ = w.Write([]byte(`{"status":"success","data":{"id":"u1"}}`))
		case r.URL.Path == "/connect/revoke":
			if r.URL.Query().Get("token") != "tok" {
				t.Fatalf("token missing")
			}
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("path = %s", r.URL.Path)
		}
	})
	defer server.Close()

	link, err := client.CreateConnectLink(ConnectLinkRequest{
		ClientID: "cid", RedirectURI: "https://app/cb", Mode: "live", Scope: "payments:read",
	})
	if err != nil {
		t.Fatal(err)
	}
	if link.Data.URL != "https://connect" {
		t.Fatal(link)
	}
	user, err := client.GetConnectUser("tok")
	if err != nil {
		t.Fatal(err)
	}
	if user.Data["id"] != "u1" {
		t.Fatal(user)
	}
	if err := client.RevokeAccessToken("tok"); err != nil {
		t.Fatal(err)
	}
}

func TestVirtualAccountAPIs(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/virtual-account/api/customers/create":
			_, _ = w.Write([]byte(`{"status":"success","data":{"id":"cust-1","email":"a@b.c","first_name":"A","last_name":"B"}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/virtual-account/api/customers":
			_, _ = w.Write([]byte(`{"status":"success","data":[{"id":"cust-1"}]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/virtual-account/api/customers/cust-1":
			_, _ = w.Write([]byte(`{"status":"success","data":{"id":"cust-1"}}`))
		case r.Method == http.MethodPut && r.URL.Path == "/virtual-account/api/customers/cust-1":
			_, _ = w.Write([]byte(`{"status":"success","data":{"id":"cust-1","first_name":"Ada"}}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/virtual-account/api/customers/cust-1":
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/virtual-account/api/customers/cust-1/virtual-account":
			_, _ = w.Write([]byte(`{"status":"success","data":{"account_number":"123"}}`))
		case r.URL.Path == "/virtual-account/api/customers/cust-1/virtual-account/deactivate":
			_, _ = w.Write([]byte(`{"status":"success","data":{"active":false}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/virtual-account/api/customers/cust-1/virtual-account/reactivate":
			_, _ = w.Write([]byte(`{"status":"success","data":{"active":true}}`))
		case r.URL.Path == "/virtual-account/api/customers/cust-1/virtual-account/activities":
			_, _ = w.Write([]byte(`{"status":"success","data":{"items":[]}}`))
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	})
	defer server.Close()

	created, err := client.CreateCustomer(VirtualCustomerRequest{Email: "a@b.c", FirstName: "A", LastName: "B"})
	if err != nil || created.Data.ID != "cust-1" {
		t.Fatalf("create: %+v %v", created, err)
	}
	list, err := client.GetCustomers("1", "10")
	if err != nil || len(list.Data) != 1 {
		t.Fatalf("list: %+v %v", list, err)
	}
	one, err := client.GetCustomer("cust-1")
	if err != nil || one.Data.ID != "cust-1" {
		t.Fatalf("get: %+v %v", one, err)
	}
	upd, err := client.UpdateCustomer("cust-1", VirtualCustomerRequest{FirstName: "Ada", Email: "a@b.c", LastName: "B"})
	if err != nil || upd.Data.FirstName != "Ada" {
		t.Fatalf("update: %+v %v", upd, err)
	}
	if _, err := client.CreateUSAccount("cust-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.DeactivateUSAccount("cust-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.ReactivateUSAccount("cust-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.GetUSAccountActivity("cust-1"); err != nil {
		t.Fatal(err)
	}
	if err := client.DeleteCustomer("cust-1"); err != nil {
		t.Fatal(err)
	}
}

func TestAPIErrorParsing(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"status":"failed","message":"Session has expired"}`))
	})
	defer server.Close()

	_, err := client.GetBalance("MWK")
	if err == nil || !strings.Contains(err.Error(), "Session has expired") {
		t.Fatalf("err = %v", err)
	}
}

func TestValidationErrorParsing(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"status":"failed","message":{"amount":["required"]}}`))
	})
	defer server.Close()

	_, err := client.ChargeMobileMoney(ChargeMobileMoneyRequest{})
	if err == nil || !strings.Contains(err.Error(), "amount: required") {
		t.Fatalf("err = %v", err)
	}
}
