package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user from the user service.
type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
}

// Company represents a company.
type Company struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
}

// PaymentMethod represents the method of payment for transactions
type PaymentMethod string

const (
	PaymentMethodCash         PaymentMethod = "CASH"
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentMethodMortgage     PaymentMethod = "MORTGAGE"
	PaymentMethodCreditCard   PaymentMethod = "CREDIT_CARD"
	PaymentMethodCheck        PaymentMethod = "CHECK"
	PaymentMethodCrypto       PaymentMethod = "CRYPTO"
)

// TransactionStatus represents the status of a financial transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusRefunded  TransactionStatus = "REFUNDED"
)

// LeaseStatus represents the status of a lease contract
type LeaseStatus string

const (
	LeaseStatusActive     LeaseStatus = "ACTIVE"
	LeaseStatusTerminated LeaseStatus = "TERMINATED"
	LeaseStatusPending    LeaseStatus = "PENDING"
	LeaseStatusExpired    LeaseStatus = "EXPIRED"
)

// PaymentFrequency represents how often rent is paid
type PaymentFrequency string

const (
	PaymentFrequencyMonthly   PaymentFrequency = "MONTHLY"
	PaymentFrequencyQuarterly PaymentFrequency = "QUARTERLY"
	PaymentFrequencyYearly    PaymentFrequency = "YEARLY"
	PaymentFrequencyWeekly    PaymentFrequency = "WEEKLY"
)

// SaleTransaction represents a property sale transaction
type SaleTransaction struct {
	ID            uuid.UUID             `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BuildingID    uuid.UUID             `gorm:"type:uuid;not null" json:"building_id"`
	BuyerID       uuid.UUID             `gorm:"type:uuid;not null" json:"buyer_id"`
	SellerID      uuid.UUID             `gorm:"type:uuid;not null" json:"seller_id"`
	AgentID       *uuid.UUID            `gorm:"type:uuid" json:"agent_id"`
	Price         float64               `gorm:"type:decimal(15,2);not null" json:"price"`
	PaymentMethod PaymentMethod         `gorm:"type:varchar(20);not null" json:"payment_method"`
	Status        TransactionStatus     `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	Commission    float64               `gorm:"type:decimal(15,2);default:0" json:"commission"`
	TaxAmount     float64               `gorm:"type:decimal(15,2);default:0" json:"tax_amount"`
	Fees          float64               `gorm:"type:decimal(15,2);default:0" json:"fees"`
	TotalAmount   float64               `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	Notes         string                `gorm:"type:text" json:"notes"`
	Documents     []TransactionDocument `gorm:"foreignKey:TransactionID" json:"documents"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	CompletedAt   *time.Time            `json:"completed_at"`

	// Relationships
	Building Building `gorm:"foreignKey:BuildingID" json:"building"`
	Buyer    User     `gorm:"foreignKey:BuyerID" json:"buyer"`
	Seller   User     `gorm:"foreignKey:SellerID" json:"seller"`
	Agent    *User    `gorm:"foreignKey:AgentID" json:"agent"`
}

// LeaseContract represents a property lease agreement
type LeaseContract struct {
	ID                uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PropertyID        uuid.UUID        `gorm:"type:uuid;not null" json:"property_id"`
	TenantID          uuid.UUID        `gorm:"type:uuid;not null" json:"tenant_id"`
	LandlordID        uuid.UUID        `gorm:"type:uuid;not null" json:"landlord_id"`
	AgentID           *uuid.UUID       `gorm:"type:uuid" json:"agent_id"`
	DurationMonths    int              `gorm:"not null" json:"duration_months"`
	StartDate         time.Time        `gorm:"not null" json:"start_date"`
	EndDate           time.Time        `gorm:"not null" json:"end_date"`
	MonthlyRent       float64          `gorm:"type:decimal(15,2);not null" json:"monthly_rent"`
	DepositAmount     float64          `gorm:"type:decimal(15,2);not null" json:"deposit_amount"`
	PaymentFrequency  PaymentFrequency `gorm:"type:varchar(20);not null;default:'MONTHLY'" json:"payment_frequency"`
	Status            LeaseStatus      `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	ContractFileURL   string           `gorm:"type:text" json:"contract_file_url"`
	ContractFileKey   string           `gorm:"type:text" json:"contract_file_key"` // Encrypted file key
	UtilitiesIncluded bool             `gorm:"default:false" json:"utilities_included"`
	PetAllowed        bool             `gorm:"default:false" json:"pet_allowed"`
	Furnished         bool             `gorm:"default:false" json:"furnished"`
	LateFeeAmount     float64          `gorm:"type:decimal(15,2);default:0" json:"late_fee_amount"`
	GracePeriodDays   int              `gorm:"default:5" json:"grace_period_days"`
	Notes             string           `gorm:"type:text" json:"notes"`
	Payments          []LeasePayment   `gorm:"foreignKey:LeaseID" json:"payments"`
	Documents         []LeaseDocument  `gorm:"foreignKey:LeaseID" json:"documents"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	SignedAt          *time.Time       `json:"signed_at"`
	TerminatedAt      *time.Time       `json:"terminated_at"`

	// Relationships
	Property Property `gorm:"foreignKey:PropertyID" json:"property"`
	Tenant   User     `gorm:"foreignKey:TenantID" json:"tenant"`
	Landlord User     `gorm:"foreignKey:LandlordID" json:"landlord"`
	Agent    *User    `gorm:"foreignKey:AgentID" json:"agent"`
}

// Offer represents promotional offers for properties or companies
type Offer struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title           string     `gorm:"type:varchar(255);not null" json:"title"`
	Description     string     `gorm:"type:text;not null" json:"description"`
	DiscountPercent int        `gorm:"not null" json:"discount_percent"`
	DiscountAmount  float64    `gorm:"type:decimal(15,2)" json:"discount_amount"`
	StartDate       time.Time  `gorm:"not null" json:"start_date"`
	EndDate         time.Time  `gorm:"not null" json:"end_date"`
	BuildingID      *uuid.UUID `gorm:"type:uuid" json:"building_id"`
	CompanyID       *uuid.UUID `gorm:"type:uuid" json:"company_id"`
	ImageURL        string     `gorm:"type:text" json:"image_url"`
	TermsConditions string     `gorm:"type:text" json:"terms_conditions"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	MaxUses         int        `gorm:"default:0" json:"max_uses"` // 0 = unlimited
	CurrentUses     int        `gorm:"default:0" json:"current_uses"`
	MinAmount       float64    `gorm:"type:decimal(15,2);default:0" json:"min_amount"`
	MaxAmount       float64    `gorm:"type:decimal(15,2);default:0" json:"max_amount"`
	Code            string     `gorm:"type:varchar(50);unique" json:"code"` // Promo code
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relationships
	Building *Building  `gorm:"foreignKey:BuildingID" json:"building"`
	Company  *Company   `gorm:"foreignKey:CompanyID" json:"company"`
	Uses     []OfferUse `gorm:"foreignKey:OfferID" json:"uses"`
}

// OfferUse tracks when offers are used
type OfferUse struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OfferID  uuid.UUID `gorm:"type:uuid;not null" json:"offer_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Amount   float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Discount float64   `gorm:"type:decimal(15,2);not null" json:"discount"`
	UsedAt   time.Time `json:"used_at"`

	// Relationships
	Offer Offer `gorm:"foreignKey:OfferID" json:"offer"`
	User  User  `gorm:"foreignKey:UserID" json:"user"`
}

// TransactionDocument represents documents related to sale transactions
type TransactionDocument struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TransactionID uuid.UUID `gorm:"type:uuid;not null" json:"transaction_id"`
	DocumentType  string    `gorm:"type:varchar(100);not null" json:"document_type"`
	FileName      string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileURL       string    `gorm:"type:text;not null" json:"file_url"`
	FileKey       string    `gorm:"type:text;not null" json:"file_key"` // Encrypted file key
	FileSize      int64     `json:"file_size"`
	UploadedAt    time.Time `json:"uploaded_at"`

	// Relationships
	Transaction SaleTransaction `gorm:"foreignKey:TransactionID" json:"transaction"`
}

// LeaseDocument represents documents related to lease contracts
type LeaseDocument struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LeaseID      uuid.UUID `gorm:"type:uuid;not null" json:"lease_id"`
	DocumentType string    `gorm:"type:varchar(100);not null" json:"document_type"`
	FileName     string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileURL      string    `gorm:"type:text;not null" json:"file_url"`
	FileKey      string    `gorm:"type:text;not null" json:"file_key"` // Encrypted file key
	FileSize     int64     `json:"file_size"`
	UploadedAt   time.Time `json:"uploaded_at"`

	// Relationships
	Lease LeaseContract `gorm:"foreignKey:LeaseID" json:"lease"`
}

// LeasePayment represents rent payments for lease contracts
type LeasePayment struct {
	ID            uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LeaseID       uuid.UUID         `gorm:"type:uuid;not null" json:"lease_id"`
	Amount        float64           `gorm:"type:decimal(15,2);not null" json:"amount"`
	PaymentDate   time.Time         `gorm:"not null" json:"payment_date"`
	DueDate       time.Time         `gorm:"not null" json:"due_date"`
	Status        TransactionStatus `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	PaymentMethod PaymentMethod     `gorm:"type:varchar(20);not null" json:"payment_method"`
	LateFee       float64           `gorm:"type:decimal(15,2);default:0" json:"late_fee"`
	Notes         string            `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`

	// Relationships
	Lease LeaseContract `gorm:"foreignKey:LeaseID" json:"lease"`
}

// FinancialReport represents financial summaries
type FinancialReport struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ReportType       string    `gorm:"type:varchar(50);not null" json:"report_type"` // "SALES", "RENTALS", "REVENUE"
	PeriodStart      time.Time `gorm:"not null" json:"period_start"`
	PeriodEnd        time.Time `gorm:"not null" json:"period_end"`
	TotalSales       float64   `gorm:"type:decimal(15,2);default:0" json:"total_sales"`
	TotalRentals     float64   `gorm:"type:decimal(15,2);default:0" json:"total_rentals"`
	TotalRevenue     float64   `gorm:"type:decimal(15,2);default:0" json:"total_revenue"`
	TotalCommissions float64   `gorm:"type:decimal(15,2);default:0" json:"total_commissions"`
	TotalTaxes       float64   `gorm:"type:decimal(15,2);default:0" json:"total_taxes"`
	TotalFees        float64   `gorm:"type:decimal(15,2);default:0" json:"total_fees"`
	ReportData       string    `gorm:"type:jsonb" json:"report_data"` // Detailed breakdown
	GeneratedAt      time.Time `json:"generated_at"`
	GeneratedBy      uuid.UUID `gorm:"type:uuid;not null" json:"generated_by"`

	// Relationships
	GeneratedByUser User `gorm:"foreignKey:GeneratedBy" json:"generated_by_user"`
}
