package paychangu

import "time"

// The Request struct is used to initiate a
// payment request with PayChangu.
// It includes essential details such as
// the transaction amount, customer details,
// and callback URLs.
type Request struct {
	// Amount specifies the transaction
	// amount in the selected currency.
	// Example: 100.50
	Amount float32 `json:"amount"`

	// Currency defines the currency code
	// for the transaction, e.g., 'MWK' or 'USD'.
	// Example: "USD"
	Currency string `json:"currency"`

	// Email is an optional field for the
	// customer's email address, used for notifications.
	// Example: "customer@example.com"
	Email string `json:"email"`

	// FirstName is the required first name of the customer.
	// Example: "John"
	FirstName string `json:"first_name"`

	// LastName is the optional last name of the customer.
	// Example: "Doe"
	LastName string `json:"last_name"`

	// CallbackURL is the URL to redirect the
	// customer after a successful payment.
	// Example: "https://example.com/callback"
	CallbackURL string `json:"callback_url"`

	// ReturnURL is the URL to redirect the
	// customer after a failed transaction.
	// Example: "https://example.com/return"
	ReturnURL string `json:"return_url"`

	// TxRef is a unique transaction reference
	// that must be unique for each request.
	// Example: "TX12345ABC"
	TxRef string `json:"tx_ref"`

	// Customization provides a title and
	// description for the payment, shown on the checkout page.
	Customization struct {
		// Title is the title for the payment,
		// shown to the customer.
		Title string `json:"title"`

		// Description gives a brief description of the payment.
		Description string `json:"description"`
	} `json:"customization"`

	// Meta allows additional data to be passed
	// with the transaction, such as a unique
	// customer identifier.
	Meta struct {
		// UUID is a unique identifier
		// associated with the transaction.
		UUID string `json:"uuid"`

		// Response can store any custom
		// information for tracking.
		Response string `json:"response"`
	} `json:"meta"`
}

// The PayChanguResponse struct represents a
// successful response from the PayChangu service,
// including details such as the checkout URL
// for customer redirection and the transaction specifics.
type Response struct {
	// Message gives a general message
	// about the transaction request result.
	Message string `json:"message"`

	// Status indicates the status of
	// the request, typically "success".
	Status string `json:"status"`

	// Data holds further details about
	// the transaction, including the checkout URL.
	Data struct {
		// Event specifies the event type of the transaction.
		Event string `json:"event"`

		// CheckoutURL is the URL to which the customer
		// should be redirected to complete the payment.
		CheckoutURL string `json:"checkout_url"`

		// Data contains details such as transaction
		// reference, currency, and amount.
		Data struct {
			// TxRef is the unique transaction
			// reference from the request.
			TxRef string `json:"tx_ref"`

			// Currency indicates the transaction currency.
			Currency string `json:"currency"`

			// Amount specifies the transaction
			// amount in the given currency.
			Amount float64 `json:"amount"`

			// Mode describes the payment mode, e.g., "online".
			Mode string `json:"mode"`

			// Status reflects the current status
			// of the transaction, e.g., "pending".
			Status string `json:"status"`
		} `json:"data"`
	} `json:"data"`
}

// The Error struct is used to capture errors
// returned by the PayChangu API.
type Error struct {
	// Status indicates the response
	// status, typically "error".
	Status string `json:"status"`

	// Message provides a detailed
	// error message from the API.
	Message string `json:"message"`
}

// The VerifyPaymentResponse struct represents
// the response returned by PayChangu upon
// verification of a payment, containing the PaymentDetails.
type VerifyPaymentResponse struct {
	// Status indicates the response status,
	// typically "success" or "error".
	Status string `json:"status"`

	// Message provides a description
	// of the response result.
	Message string `json:"message"`

	// Data contains detailed information
	// about the verified payment.
	Data PaymentDetails `json:"data"`
}

// PaymentDetails provides comprehensive information
// about a verified payment, including customer
// information, transaction attempts, and authorization details.
type PaymentDetails struct {
	// EventType describes the type of
	// event, e.g., "payment_success".
	EventType string `json:"event_type"`

	// TxRef is the unique transaction reference.
	TxRef string `json:"tx_ref"`

	// Mode describes the payment mode, e.g., "online".
	Mode string `json:"mode"`

	// Type describes the type of payment.
	Type string `json:"type"`

	// Status represents the payment status,
	// e.g., "completed".
	Status string `json:"status"`

	// Attempts indicates the number of
	// attempts made for this payment.
	Attempts int `json:"number_of_attempts"`

	// Reference is an internal or external
	// reference for the transaction.
	Reference string `json:"reference"`

	// Currency of the transaction, e.g., "USD".
	Currency string `json:"currency"`

	// Amount charged in the transaction.
	Amount float64 `json:"amount"`

	// Charges represents the fees
	// associated with the transaction.
	Charges float64 `json:"charges"`

	// Customization provides display
	// customization details for the transaction.
	Customization Customization `json:"customization"`

	// Meta stores any additional metadata
	// associated with the transaction.
	Meta interface{} `json:"meta"`

	// Authorization contains payment authorization details.
	Authorization PaymentAuthorization `json:"authorization"`

	// Customer holds customer information
	// associated with the transaction.
	Customer CustomerInfo `json:"customer"`

	// Logs provides a list of logs related
	// to the payment processing steps.
	Logs []PaymentLog `json:"logs"`

	// CreatedAt is the timestamp of when
	// the payment was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the timestamp of the
	// last update to the payment.
	UpdatedAt time.Time `json:"updated_at"`
}

// The Customization struct is used to customize
// the appearance of the payment interface.
type Customization struct {
	// Title is a display title for the transaction.
	Title string `json:"title"`

	// Description provides a brief
	// description of the transaction.
	Description string `json:"description"`

	// Logo is an optional URL to the logo
	// displayed on the payment page.
	Logo string `json:"logo"`
}

// The PaymentAuthorization struct captures
// authorization details for the payment,
// such as card information or mobile details.
type PaymentAuthorization struct {
	// Channel specifies the authorization
	// channel, e.g., "card" or "mobile".
	Channel string `json:"channel"`

	// CardNumber shows the masked
	// card number used for authorization.
	CardNumber string `json:"card_number"`

	// Expiry is the expiry date of the card used.
	Expiry string `json:"expiry"`

	// Brand is the card brand, such
	// as "Visa" or "MasterCard".
	Brand string `json:"brand"`

	// Provider is the provider of the
	// payment service, e.g., "PayChangu".
	Provider string `json:"provider"`

	// MobileNumber is the mobile number if
	// payment was made through mobile.
	MobileNumber string `json:"mobile_number"`

	// CompletedAt provides the timestamp
	// when the authorization was completed.
	CompletedAt string `json:"completed_at"`
}

// The CustomerInfo struct captures basic
// information about the customer involved
// in the payment.
type CustomerInfo struct {
	// Email is the customer's email address.
	Email string `json:"email"`

	// FirstName is the customer's first name.
	FirstName string `json:"first_name"`

	// LastName is the customer's last name.
	LastName string `json:"last_name"`
}

// The PaymentLog struct represents a
// log entry related to the transaction,
// useful for debugging and tracking payment events.
type PaymentLog struct {
	// Type indicates the type of log entry, e.g., "info" or "error".
	Type string `json:"type"`

	// Message provides a message detailing the log event.
	Message string `json:"message"`

	// CreatedAt records the time the log entry was created.
	CreatedAt time.Time `json:"created_at"`
}

////////Payout via Mobile Money////////

// MobileMoneyOperator represents a supported mobile money operator.
type MobileMoneyOperator struct {
	ID                  int     `json:"id"`
	Name                string  `json:"name"`
	RefID               string  `json:"ref_id"`
	LiveMode            int     `json:"live_mode"`
	ShortCode           string  `json:"short_code"`
	Logo                *string `json:"logo"` // Can be null
	OperatorFee         string  `json:"operator_fee"`
	PaymentPercentFee   string  `json:"payment_percent_fee"`
	PaymentFiatFee      *string `json:"payment_fiat_fee"`   // Can be null
	PayoutPercentFee    *string `json:"payout_percent_fee"` // Can be null
	PayoutFiatFee       *string `json:"payout_fiat_fee"`    // Can be null
	SupportsWithdrawals bool    `json:"supports_withdrawals"`
	SupportedCountry    struct {
		Name     string `json:"name"`
		Currency string `json:"currency"`
	} `json:"supported_country"`
}

// MobileMoneyOperatorsResponse is the response structure for fetching mobile money operators.
type MobileMoneyOperatorsResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    []MobileMoneyOperator `json:"data"`
}

// MobileMoneyPayoutRequest is the payload for initiating a mobile money payout.
type MobileMoneyPayoutRequest struct {
	Mobile                   string  `json:"mobile"`
	MobileMoneyOperatorRefID string  `json:"mobile_money_operator_ref_id"`
	Amount                   float64 `json:"amount"` // Use float64 for amount
	ChargeID                 string  `json:"charge_id"`
	Email                    string  `json:"email,omitempty"`              // Optional
	FirstName                string  `json:"first_name,omitempty"`         // Optional
	LastName                 string  `json:"last_name,omitempty"`          // Optional
	TransactionStatus        string  `json:"transaction_status,omitempty"` // Optional, sandbox mode only
}

// PayoutTransactionDetails represents the details of a payout transaction.
type PayoutTransactionDetails struct {
	ChargeID    string    `json:"charge_id"`
	RefID       string    `json:"ref_id"`
	TransID     *string   `json:"trans_id"` // Can be null
	Currency    string    `json:"currency"`
	Amount      float64   `json:"amount"`
	FirstName   *string   `json:"first_name"` // Can be null
	LastName    *string   `json:"last_name"`  // Can be null
	Email       *string   `json:"email"`      // Can be null
	Type        string    `json:"type"`
	TraceID     *string   `json:"trace_id"` // Can be null
	Status      string    `json:"status"`
	Mobile      string    `json:"mobile"`
	Attempts    int       `json:"attempts"`
	Mode        string    `json:"mode"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	EventType   string    `json:"event_type"`
	MobileMoney struct {
		Name    string `json:"name"`
		RefID   string `json:"ref_id"`
		Country string `json:"country"`
	} `json:"mobile_money"`
	TransactionCharges struct {
		Currency string `json:"currency"`
		Amount   string `json:"amount"` // The API returns this as a string, e.g., "1.7"
	} `json:"transaction_charges"`
	Customer *interface{} `json:"customer"` // Can be null or an object
}

// MobileMoneyPayoutResponse is the response for a successful mobile money payout initialization.
type MobileMoneyPayoutResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Transaction PayoutTransactionDetails `json:"transaction"`
	} `json:"data"`
}

// MobileMoneyPayoutErrorResponse is the response for a failed mobile money payout.
type MobileMoneyPayoutErrorResponse struct {
	Status  string              `json:"status"`
	Data    interface{}         `json:"data"`    // Can be null
	Message map[string][]string `json:"message"` // Detailed error messages
}

// GetMobileMoneyPayoutDetailsResponse is the response structure for fetching mobile money payout details.
type GetMobileMoneyPayoutDetailsResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    PayoutTransactionDetails `json:"data"` // Reusing the existing PayoutTransactionDetails struct
}

// SupportedBank represents a bank supported for payouts.
type Bank struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// SupportedBanksResponse is the response structure for fetching supported banks.
type BanksResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []Bank `json:"data"`
}

// BankPayoutRequest is the payload for initiating a bank payout.
type BankPayoutRequest struct {
	PayoutMethod      string  `json:"payout_method"` // Defaults to "bank_transfer"
	BankUUID          string  `json:"bank_uuid"`
	Amount            float64 `json:"amount"` // Use float64 for amount, will be marshaled to string
	ChargeID          string  `json:"charge_id"`
	BankAccountName   string  `json:"bank_account_name"`
	BankAccountNumber string  `json:"bank_account_number"`
	Email             string  `json:"email,omitempty"`      // Optional
	FirstName         string  `json:"first_name,omitempty"` // Optional
	LastName          string  `json:"last_name,omitempty"`  // Optional
}

// RecipientAccountDetails represents the bank account details of the recipient for a bank payout.
type RecipientAccountDetails struct {
	BankUUID      string `json:"bank_uuid"`
	BankName      string `json:"bank_name"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
}

// BankPayoutTransactionDetails represents the detailed transaction information for a bank payout.
// This structure is similar to PayoutTransactionDetails but includes specific bank recipient details.
type BankPayoutTransactionDetails struct {
	ChargeID           string     `json:"charge_id"`
	RefID              string     `json:"ref_id"`
	TransID            *string    `json:"trans_id"` // Can be null
	Currency           string     `json:"currency"`
	Amount             float64    `json:"amount"`
	FirstName          *string    `json:"first_name"` // Can be null
	LastName           *string    `json:"last_name"`  // Can be null
	Email              *string    `json:"email"`      // Can be null
	Type               string     `json:"type"`
	TraceID            *string    `json:"trace_id"` // Can be null
	Status             string     `json:"status"`
	Mobile             string     `json:"mobile"` // API returns "0" for bank payouts, but still present
	Attempts           int        `json:"attempts"`
	Mode               string     `json:"mode"`
	CreatedAt          time.Time  `json:"created_at"`
	CompletedAt        *time.Time `json:"completed_at"` // Can be null
	EventType          string     `json:"event_type"`
	TransactionCharges struct {
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	} `json:"transaction_charges"`
	RecipientAccountDetails RecipientAccountDetails `json:"recipient_account_details"` // New field for bank payouts
	// Note: mobile_money field from PayoutTransactionDetails is not present here.
	// This makes it distinct from MobileMoneyPayoutTransactionDetails.
}

// BankPayoutResponse is the response for a successful bank payout initialization.
type BankPayoutResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Transaction BankPayoutTransactionDetails `json:"transaction"`
	} `json:"data"`
}

// GetBankPayoutDetailsResponse is the response structure for fetching bank payout details.
type GetBankPayoutDetailsResponse struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    BankPayoutTransactionDetails `json:"data"` // Reusing the existing BankPayoutTransactionDetails struct
}
