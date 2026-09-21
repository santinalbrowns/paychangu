package paychangu

import "time"

// Request initiates a standard checkout payment.
type Request struct {
	Amount        float32 `json:"amount"`
	Currency      string  `json:"currency"`
	Email         string  `json:"email"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	CallbackURL   string  `json:"callback_url"`
	ReturnURL     string  `json:"return_url"`
	TxRef         string  `json:"tx_ref"`
	Customization struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"customization"`
	Meta struct {
		UUID     string `json:"uuid"`
		Response string `json:"response"`
	} `json:"meta"`
}

// Response is returned after initiating a standard checkout payment.
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		Event       string `json:"event"`
		CheckoutURL string `json:"checkout_url"`
		Data        struct {
			TxRef    string  `json:"tx_ref"`
			Currency string  `json:"currency"`
			Amount   float64 `json:"amount"`
			Mode     string  `json:"mode"`
			Status   string  `json:"status"`
		} `json:"data"`
	} `json:"data"`
}

// Error is a generic API error payload.
type Error struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Status
}

// VerifyPaymentResponse wraps verified payment details.
type VerifyPaymentResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    PaymentDetails `json:"data"`
}

// PaymentDetails describes a verified standard checkout payment.
type PaymentDetails struct {
	EventType     string               `json:"event_type"`
	TxRef         string               `json:"tx_ref"`
	Mode          string               `json:"mode"`
	Type          string               `json:"type"`
	Status        string               `json:"status"`
	Attempts      int                  `json:"number_of_attempts"`
	Reference     string               `json:"reference"`
	Currency      string               `json:"currency"`
	Amount        float64              `json:"amount"`
	Charges       float64              `json:"charges"`
	Customization Customization        `json:"customization"`
	Meta          interface{}          `json:"meta"`
	Authorization PaymentAuthorization `json:"authorization"`
	Customer      CustomerInfo         `json:"customer"`
	Logs          []PaymentLog         `json:"logs"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// Customization customizes checkout display.
type Customization struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
}

// PaymentAuthorization holds authorization channel details.
type PaymentAuthorization struct {
	Channel      string `json:"channel"`
	CardNumber   string `json:"card_number"`
	Expiry       string `json:"expiry"`
	Brand        string `json:"brand"`
	Provider     string `json:"provider"`
	MobileNumber string `json:"mobile_number"`
	CompletedAt  string `json:"completed_at"`
}

// CustomerInfo is basic customer contact info.
type CustomerInfo struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// PaymentLog is a transaction log entry.
type PaymentLog struct {
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// MobileMoneyOperator is a supported mobile money operator.
type MobileMoneyOperator struct {
	ID                  int     `json:"id"`
	Name                string  `json:"name"`
	RefID               string  `json:"ref_id"`
	LiveMode            int     `json:"live_mode"`
	ShortCode           string  `json:"short_code"`
	Logo                *string `json:"logo"`
	OperatorFee         string  `json:"operator_fee"`
	PaymentPercentFee   string  `json:"payment_percent_fee"`
	PaymentFiatFee      *string `json:"payment_fiat_fee"`
	PayoutPercentFee    *string `json:"payout_percent_fee"`
	PayoutFiatFee       *string `json:"payout_fiat_fee"`
	SupportsWithdrawals bool    `json:"supports_withdrawals"`
	SupportedCountry    struct {
		Name     string `json:"name"`
		Currency string `json:"currency"`
	} `json:"supported_country"`
}

// MobileMoneyOperatorsResponse is the operators list response.
type MobileMoneyOperatorsResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    []MobileMoneyOperator `json:"data"`
}

// MobileMoneyPayoutRequest initiates a mobile money payout.
type MobileMoneyPayoutRequest struct {
	Mobile                   string  `json:"mobile"`
	MobileMoneyOperatorRefID string  `json:"mobile_money_operator_ref_id"`
	Amount                   float64 `json:"amount"`
	ChargeID                 string  `json:"charge_id"`
	Email                    string  `json:"email,omitempty"`
	FirstName                string  `json:"first_name,omitempty"`
	LastName                 string  `json:"last_name,omitempty"`
	TransactionStatus        string  `json:"transaction_status,omitempty"`
}

// TransactionCharges holds fee details.
type TransactionCharges struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

// MobileMoneyInfo describes the operator used on a charge/payout.
type MobileMoneyInfo struct {
	Name    string `json:"name"`
	RefID   string `json:"ref_id"`
	Country string `json:"country"`
}

// PayoutTransactionDetails describes a mobile money payout transaction.
type PayoutTransactionDetails struct {
	ChargeID           string             `json:"charge_id"`
	RefID              string             `json:"ref_id"`
	TransID            *string            `json:"trans_id"`
	Currency           string             `json:"currency"`
	Amount             float64            `json:"amount"`
	FirstName          *string            `json:"first_name"`
	LastName           *string            `json:"last_name"`
	Email              *string            `json:"email"`
	Type               string             `json:"type"`
	TraceID            *string            `json:"trace_id"`
	Status             string             `json:"status"`
	Mobile             string             `json:"mobile"`
	Attempts           int                `json:"attempts"`
	Mode               string             `json:"mode"`
	CreatedAt          time.Time          `json:"created_at"`
	CompletedAt        time.Time          `json:"completed_at"`
	EventType          string             `json:"event_type"`
	MobileMoney        MobileMoneyInfo    `json:"mobile_money"`
	TransactionCharges TransactionCharges `json:"transaction_charges"`
	Customer           *interface{}       `json:"customer"`
}

// MobileMoneyPayoutResponse is returned after initiating a MoMo payout.
type MobileMoneyPayoutResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Transaction PayoutTransactionDetails `json:"transaction"`
	} `json:"data"`
}

// MobileMoneyPayoutErrorResponse is a validation-style error for payouts.
type MobileMoneyPayoutErrorResponse struct {
	Status  string              `json:"status"`
	Data    interface{}         `json:"data"`
	Message map[string][]string `json:"message"`
}

// GetMobileMoneyPayoutDetailsResponse wraps payout details.
type GetMobileMoneyPayoutDetailsResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    PayoutTransactionDetails `json:"data"`
}

// Bank is a supported bank for payouts.
type Bank struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// BanksResponse wraps supported banks.
type BanksResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []Bank `json:"data"`
}

// BankPayoutRequest initiates a bank payout.
type BankPayoutRequest struct {
	PayoutMethod      string  `json:"payout_method"`
	BankUUID          string  `json:"bank_uuid"`
	Amount            float64 `json:"amount"`
	ChargeID          string  `json:"charge_id"`
	BankAccountName   string  `json:"bank_account_name"`
	BankAccountNumber string  `json:"bank_account_number"`
	Email             string  `json:"email,omitempty"`
	FirstName         string  `json:"first_name,omitempty"`
	LastName          string  `json:"last_name,omitempty"`
}

// RecipientAccountDetails is the bank recipient on a payout.
type RecipientAccountDetails struct {
	BankUUID      string `json:"bank_uuid"`
	BankName      string `json:"bank_name"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
}

// BankPayoutTransactionDetails describes a bank payout.
type BankPayoutTransactionDetails struct {
	ChargeID                string                  `json:"charge_id"`
	RefID                   string                  `json:"ref_id"`
	TransID                 *string                 `json:"trans_id"`
	Currency                string                  `json:"currency"`
	Amount                  float64                 `json:"amount"`
	FirstName               *string                 `json:"first_name"`
	LastName                *string                 `json:"last_name"`
	Email                   *string                 `json:"email"`
	Type                    string                  `json:"type"`
	TraceID                 *string                 `json:"trace_id"`
	Status                  string                  `json:"status"`
	Mobile                  string                  `json:"mobile"`
	Attempts                int                     `json:"attempts"`
	Mode                    string                  `json:"mode"`
	CreatedAt               time.Time               `json:"created_at"`
	CompletedAt             *time.Time              `json:"completed_at"`
	EventType               string                  `json:"event_type"`
	TransactionCharges      TransactionCharges      `json:"transaction_charges"`
	RecipientAccountDetails RecipientAccountDetails `json:"recipient_account_details"`
}

// BankPayoutResponse is returned after initiating a bank payout.
type BankPayoutResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Transaction BankPayoutTransactionDetails `json:"transaction"`
	} `json:"data"`
}

// GetBankPayoutDetailsResponse wraps bank payout details.
type GetBankPayoutDetailsResponse struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    BankPayoutTransactionDetails `json:"data"`
}

// BalanceData is the wallet balance payload.
type BalanceData struct {
	Environment       string      `json:"environment"`
	Currency          string      `json:"currency"`
	MainBalance       string      `json:"main_balance"`
	CollectionBalance interface{} `json:"collection_balance"`
}

// BalanceResponse wraps wallet balances.
type BalanceResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    BalanceData `json:"data"`
}

// ChargeMobileMoneyRequest initiates a direct MoMo collection charge.
type ChargeMobileMoneyRequest struct {
	Mobile                   string `json:"mobile"`
	MobileMoneyOperatorRefID string `json:"mobile_money_operator_ref_id"`
	Amount                   string `json:"amount"`
	ChargeID                 string `json:"charge_id"`
	Email                    string `json:"email,omitempty"`
	FirstName                string `json:"first_name,omitempty"`
	LastName                 string `json:"last_name,omitempty"`
}

// DirectChargeDetails describes a direct charge / MoMo payment.
type DirectChargeDetails struct {
	Amount             float64               `json:"amount"`
	ChargeID           string                `json:"charge_id"`
	RefID              string                `json:"ref_id"`
	TransID            *string               `json:"trans_id"`
	FirstName          *string               `json:"first_name"`
	LastName           *string               `json:"last_name"`
	Email              *string               `json:"email"`
	Type               string                `json:"type"`
	TraceID            *string               `json:"trace_id"`
	Status             string                `json:"status"`
	Mobile             string                `json:"mobile"`
	Attempts           int                   `json:"attempts"`
	Currency           string                `json:"currency"`
	Mode               string                `json:"mode"`
	CreatedAt          time.Time             `json:"created_at"`
	CompletedAt        *time.Time            `json:"completed_at"`
	EventType          string                `json:"event_type"`
	MobileMoney        *MobileMoneyInfo      `json:"mobile_money,omitempty"`
	TransactionCharges *TransactionCharges   `json:"transaction_charges,omitempty"`
	Authorization      *PaymentAuthorization `json:"authorization,omitempty"`
	Logs               []PaymentLog          `json:"logs,omitempty"`
}

// ChargeMobileMoneyResponse is returned after initiating a MoMo charge.
type ChargeMobileMoneyResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    DirectChargeDetails `json:"data"`
}

// DirectChargeStatusResponse wraps verify/details for direct charges.
type DirectChargeStatusResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    DirectChargeDetails `json:"data"`
}

// BankTransferRequest initiates a bank transfer collection.
type BankTransferRequest struct {
	Amount                 string `json:"amount"`
	Currency               string `json:"currency"`
	PaymentMethod          string `json:"payment_method"`
	ChargeID               string `json:"charge_id"`
	Email                  string `json:"email,omitempty"`
	FirstName              string `json:"first_name,omitempty"`
	LastName               string `json:"last_name,omitempty"`
	Mobile                 string `json:"mobile,omitempty"`
	CreatePermanentAccount bool   `json:"create_permanent_account,omitempty"`
}

// PaymentAccountDetails is the virtual account for bank transfer.
type PaymentAccountDetails struct {
	BankName                   string `json:"bank_name"`
	AccountNumber              string `json:"account_number"`
	AccountName                string `json:"account_name"`
	AccountExpirationTimestamp int64  `json:"account_expiration_timestamp"`
}

// BankTransferTransaction is the transaction nested in bank transfer responses.
type BankTransferTransaction struct {
	ChargeID           string                  `json:"charge_id"`
	RefID              string                  `json:"ref_id"`
	TransID            *string                 `json:"trans_id"`
	Currency           string                  `json:"currency"`
	Amount             float64                 `json:"amount"`
	FirstName          *string                 `json:"first_name"`
	LastName           *string                 `json:"last_name"`
	Email              *string                 `json:"email"`
	Type               string                  `json:"type"`
	TraceID            *string                 `json:"trace_id"`
	Status             string                  `json:"status"`
	Mobile             string                  `json:"mobile"`
	Attempts           int                     `json:"attempts"`
	Mode               string                  `json:"mode"`
	CreatedAt          time.Time               `json:"created_at"`
	CompletedAt        *time.Time              `json:"completed_at"`
	EventType          string                  `json:"event_type"`
	TransactionCharges *TransactionCharges     `json:"transaction_charges,omitempty"`
	Authorization      *map[string]interface{} `json:"authorization,omitempty"`
	Logs               []PaymentLog            `json:"logs,omitempty"`
}

// BankTransferResponse is returned after initiating a bank transfer charge.
type BankTransferResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		PaymentAccountDetails PaymentAccountDetails   `json:"payment_account_details"`
		Transaction           BankTransferTransaction `json:"transaction"`
	} `json:"data"`
}

// BankTransferDetailsResponse wraps a single bank transfer transaction.
type BankTransferDetailsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Transaction BankTransferTransaction `json:"transaction"`
	} `json:"data"`
}

// ChargeCardRequest charges a card directly.
type ChargeCardRequest struct {
	CardNumber     string `json:"card_number"`
	Expiry         string `json:"expiry"`
	CVV            string `json:"cvv"`
	CardholderName string `json:"cardholder_name"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	Email          string `json:"email,omitempty"`
	ChargeID       string `json:"charge_id"`
	RedirectURL    string `json:"redirect_url"`
}

// ChargeCardResponse is returned after charging a card.
type ChargeCardResponse struct {
	Success         bool   `json:"success"`
	Requires3DSAuth bool   `json:"requires_3ds_auth"`
	OrderReference  string `json:"orderReference"`
	ThreeDSAuthLink string `json:"3ds_auth_link"`
}

// CardChargeStatusResponse wraps card verify/refund responses.
type CardChargeStatusResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// PaginatedBankPayouts is a paginated list of bank payouts.
type PaginatedBankPayouts struct {
	CurrentPage int                            `json:"current_page"`
	TotalPages  int                            `json:"total_pages"`
	PerPage     int                            `json:"per_page"`
	NextPageURL *string                        `json:"next_page_url"`
	Data        []BankPayoutTransactionDetails `json:"data"`
}

// ListBankPayoutsResponse wraps all bank payouts.
type ListBankPayoutsResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Data    PaginatedBankPayouts `json:"data"`
}

// ValidateBillRequest validates a bill before payment.
type ValidateBillRequest struct {
	Biller      string `json:"biller"`
	Account     string `json:"account"`
	AccountType string `json:"account_type,omitempty"`
	Amount      string `json:"amount,omitempty"`
}

// PayBillRequest pays a bill.
type PayBillRequest struct {
	Biller       string `json:"biller"`
	Account      string `json:"account"`
	Amount       string `json:"amount,omitempty"`
	CustomerName string `json:"customer_name,omitempty"`
	AccountType  string `json:"account_type,omitempty"`
	Reference    string `json:"reference,omitempty"`
}

// AirtimeRequest buys airtime.
type AirtimeRequest struct {
	Phone     string `json:"phone"`
	Amount    string `json:"amount"`
	Reference string `json:"reference,omitempty"`
}

// BillsAPIResponse is a generic bills endpoint wrapper.
type BillsAPIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ConnectLinkRequest generates a Connect authorization URL.
type ConnectLinkRequest struct {
	ClientID      string
	RedirectURI   string
	Scope         string
	Mode          string
	WebhookURL    string
	WebhookSecret string
}

// ConnectLinkResponse holds the generated Connect link.
type ConnectLinkResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		URL string `json:"url"`
	} `json:"data"`
	URL string `json:"url"`
}

// ConnectUserResponse is the connected user info.
type ConnectUserResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// VirtualCustomerRequest creates or updates a virtual-account customer.
type VirtualCustomerRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// VirtualCustomer is a USD virtual-account customer.
type VirtualCustomer struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// VirtualCustomerResponse wraps a single customer.
type VirtualCustomerResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    VirtualCustomer `json:"data"`
}

// VirtualCustomersResponse wraps a customer list.
type VirtualCustomersResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    []VirtualCustomer `json:"data"`
}

// VirtualAccountResponse wraps US virtual account data.
type VirtualAccountResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
