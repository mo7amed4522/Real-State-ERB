package services

import (
	"fmt"
	"my-property/go-service/database"
	"my-property/go-service/models"
	"my-property/go-service/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FinancialService struct {
	db               *gorm.DB
	encryptionService *utils.EncryptionService
	uploadDir        string
}

func NewFinancialService(db *gorm.DB, encryptionService *utils.EncryptionService) *FinancialService {
	uploadDir := "./uploads/financial"
	return &FinancialService{
		db:               db,
		encryptionService: encryptionService,
		uploadDir:        uploadDir,
	}
}

// Sale Transaction Management
func (s *FinancialService) CreateSaleTransaction(input CreateSaleTransactionInput) (*models.SaleTransaction, error) {
	// Calculate totals
	totalAmount := input.Price + input.Commission + input.TaxAmount + input.Fees

	transaction := &models.SaleTransaction{
		BuildingID:    input.BuildingID,
		BuyerID:       input.BuyerID,
		SellerID:      input.SellerID,
		AgentID:       input.AgentID,
		Price:         input.Price,
		PaymentMethod: input.PaymentMethod,
		Status:        models.TransactionStatusPending,
		Commission:    input.Commission,
		TaxAmount:     input.TaxAmount,
		Fees:          input.Fees,
		TotalAmount:   totalAmount,
		Notes:         input.Notes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.db.Create(transaction).Error; err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *FinancialService) UpdateSaleTransactionStatus(transactionID uuid.UUID, status models.TransactionStatus) (*models.SaleTransaction, error) {
	var transaction models.SaleTransaction
	if err := s.db.First(&transaction, transactionID).Error; err != nil {
		return nil, err
	}

	transaction.Status = status
	transaction.UpdatedAt = time.Now()

	if status == models.TransactionStatusCompleted {
		now := time.Now()
		transaction.CompletedAt = &now
	}

	if err := s.db.Save(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *FinancialService) GetSaleTransactions(filters SaleTransactionFilters) ([]models.SaleTransaction, error) {
	var transactions []models.SaleTransaction
	query := s.db.Preload("Building").Preload("Buyer").Preload("Seller").Preload("Agent").Preload("Documents")

	if filters.BuildingID != nil {
		query = query.Where("building_id = ?", *filters.BuildingID)
	}
	if filters.BuyerID != nil {
		query = query.Where("buyer_id = ?", *filters.BuyerID)
	}
	if filters.SellerID != nil {
		query = query.Where("seller_id = ?", *filters.SellerID)
	}
	if filters.AgentID != nil {
		query = query.Where("agent_id = ?", *filters.AgentID)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", *filters.EndDate)
	}

	err := query.Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

// Lease Contract Management
func (s *FinancialService) CreateLeaseContract(input CreateLeaseContractInput) (*models.LeaseContract, error) {
	// Calculate end date based on duration
	endDate := input.StartDate.AddDate(0, input.DurationMonths, 0)

	lease := &models.LeaseContract{
		PropertyID:       input.PropertyID,
		TenantID:         input.TenantID,
		LandlordID:       input.LandlordID,
		AgentID:          input.AgentID,
		DurationMonths:   input.DurationMonths,
		StartDate:        input.StartDate,
		EndDate:          endDate,
		MonthlyRent:      input.MonthlyRent,
		DepositAmount:    input.DepositAmount,
		PaymentFrequency: input.PaymentFrequency,
		Status:           models.LeaseStatusPending,
		UtilitiesIncluded: input.UtilitiesIncluded,
		PetAllowed:       input.PetAllowed,
		Furnished:        input.Furnished,
		LateFeeAmount:    input.LateFeeAmount,
		GracePeriodDays:  input.GracePeriodDays,
		Notes:            input.Notes,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.db.Create(lease).Error; err != nil {
		return nil, err
	}

	return lease, nil
}

func (s *FinancialService) UpdateLeaseStatus(leaseID uuid.UUID, status models.LeaseStatus) (*models.LeaseContract, error) {
	var lease models.LeaseContract
	if err := s.db.First(&lease, leaseID).Error; err != nil {
		return nil, err
	}

	lease.Status = status
	lease.UpdatedAt = time.Now()

	if status == models.LeaseStatusActive {
		now := time.Now()
		lease.SignedAt = &now
	} else if status == models.LeaseStatusTerminated {
		now := time.Now()
		lease.TerminatedAt = &now
	}

	if err := s.db.Save(&lease).Error; err != nil {
		return nil, err
	}

	return &lease, nil
}

func (s *FinancialService) GetLeaseContracts(filters LeaseContractFilters) ([]models.LeaseContract, error) {
	var leases []models.LeaseContract
	query := s.db.Preload("Property").Preload("Tenant").Preload("Landlord").Preload("Agent").Preload("Payments").Preload("Documents")

	if filters.PropertyID != nil {
		query = query.Where("property_id = ?", *filters.PropertyID)
	}
	if filters.TenantID != nil {
		query = query.Where("tenant_id = ?", *filters.TenantID)
	}
	if filters.LandlordID != nil {
		query = query.Where("landlord_id = ?", *filters.LandlordID)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.ActiveOnly {
		query = query.Where("status = ? AND end_date >= ?", models.LeaseStatusActive, time.Now())
	}

	err := query.Order("created_at DESC").Find(&leases).Error
	return leases, err
}

// Lease Payment Management
func (s *FinancialService) CreateLeasePayment(input CreateLeasePaymentInput) (*models.LeasePayment, error) {
	// Check if lease exists and is active
	var lease models.LeaseContract
	if err := s.db.First(&lease, input.LeaseID).Error; err != nil {
		return nil, fmt.Errorf("lease not found")
	}

	if lease.Status != models.LeaseStatusActive {
		return nil, fmt.Errorf("lease is not active")
	}

	// Calculate late fee if payment is late
	lateFee := 0.0
	if input.PaymentDate.After(input.DueDate.AddDate(0, 0, lease.GracePeriodDays)) {
		lateFee = lease.LateFeeAmount
	}

	payment := &models.LeasePayment{
		LeaseID:       input.LeaseID,
		Amount:        input.Amount,
		PaymentDate:   input.PaymentDate,
		DueDate:       input.DueDate,
		Status:        models.TransactionStatusPending,
		PaymentMethod: input.PaymentMethod,
		LateFee:       lateFee,
		Notes:         input.Notes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.db.Create(payment).Error; err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *FinancialService) UpdatePaymentStatus(paymentID uuid.UUID, status models.TransactionStatus) (*models.LeasePayment, error) {
	var payment models.LeasePayment
	if err := s.db.First(&payment, paymentID).Error; err != nil {
		return nil, err
	}

	payment.Status = status
	payment.UpdatedAt = time.Now()

	if err := s.db.Save(&payment).Error; err != nil {
		return nil, err
	}

	return &payment, nil
}

// Offer Management
func (s *FinancialService) CreateOffer(input CreateOfferInput) (*models.Offer, error) {
	// Generate unique promo code if not provided
	if input.Code == "" {
		input.Code = s.generatePromoCode()
	}

	offer := &models.Offer{
		Title:           input.Title,
		Description:     input.Description,
		DiscountPercent: input.DiscountPercent,
		DiscountAmount:  input.DiscountAmount,
		StartDate:       input.StartDate,
		EndDate:         input.EndDate,
		BuildingID:      input.BuildingID,
		CompanyID:       input.CompanyID,
		ImageURL:        input.ImageURL,
		TermsConditions: input.TermsConditions,
		IsActive:        true,
		MaxUses:         input.MaxUses,
		MinAmount:       input.MinAmount,
		MaxAmount:       input.MaxAmount,
		Code:            input.Code,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.db.Create(offer).Error; err != nil {
		return nil, err
	}

	return offer, nil
}

func (s *FinancialService) UseOffer(offerCode string, userID uuid.UUID, amount float64) (*models.OfferUse, error) {
	var offer models.Offer
	if err := s.db.Where("code = ? AND is_active = ?", offerCode, true).First(&offer).Error; err != nil {
		return nil, fmt.Errorf("offer not found or inactive")
	}

	// Check if offer is still valid
	now := time.Now()
	if now.Before(offer.StartDate) || now.After(offer.EndDate) {
		return nil, fmt.Errorf("offer is not valid at this time")
	}

	// Check usage limits
	if offer.MaxUses > 0 && offer.CurrentUses >= offer.MaxUses {
		return nil, fmt.Errorf("offer usage limit reached")
	}

	// Check amount limits
	if offer.MinAmount > 0 && amount < offer.MinAmount {
		return nil, fmt.Errorf("minimum amount not met")
	}
	if offer.MaxAmount > 0 && amount > offer.MaxAmount {
		return nil, fmt.Errorf("maximum amount exceeded")
	}

	// Calculate discount
	discount := 0.0
	if offer.DiscountPercent > 0 {
		discount = amount * float64(offer.DiscountPercent) / 100.0
	} else {
		discount = offer.DiscountAmount
	}

	// Create offer use record
	offerUse := &models.OfferUse{
		OfferID:  offer.ID,
		UserID:   userID,
		Amount:   amount,
		Discount: discount,
		UsedAt:   now,
	}

	if err := s.db.Create(offerUse).Error; err != nil {
		return nil, err
	}

	// Update offer usage count
	offer.CurrentUses++
	if err := s.db.Save(&offer).Error; err != nil {
		return nil, err
	}

	return offerUse, nil
}

func (s *FinancialService) GetOffers(filters OfferFilters) ([]models.Offer, error) {
	var offers []models.Offer
	query := s.db.Preload("Building").Preload("Company").Preload("Uses")

	if filters.BuildingID != nil {
		query = query.Where("building_id = ?", *filters.BuildingID)
	}
	if filters.CompanyID != nil {
		query = query.Where("company_id = ?", *filters.CompanyID)
	}
	if filters.ActiveOnly {
		now := time.Now()
		query = query.Where("is_active = ? AND start_date <= ? AND end_date >= ?", true, now, now)
	}

	err := query.Order("created_at DESC").Find(&offers).Error
	return offers, err
}

// Financial Reporting
func (s *FinancialService) GenerateFinancialReport(reportType string, periodStart, periodEnd time.Time, generatedBy uuid.UUID) (*models.FinancialReport, error) {
	var totalSales, totalRentals, totalRevenue, totalCommissions, totalTaxes, totalFees float64

	// Calculate sales revenue
	var salesTransactions []models.SaleTransaction
	s.db.Where("status = ? AND created_at BETWEEN ? AND ?", models.TransactionStatusCompleted, periodStart, periodEnd).Find(&salesTransactions)
	
	for _, transaction := range salesTransactions {
		totalSales += transaction.Price
		totalCommissions += transaction.Commission
		totalTaxes += transaction.TaxAmount
		totalFees += transaction.Fees
	}

	// Calculate rental revenue
	var leasePayments []models.LeasePayment
	s.db.Where("status = ? AND payment_date BETWEEN ? AND ?", models.TransactionStatusCompleted, periodStart, periodEnd).Find(&leasePayments)
	
	for _, payment := range leasePayments {
		totalRentals += payment.Amount
	}

	totalRevenue = totalSales + totalRentals

	// Create detailed report data
	reportData := map[string]interface{}{
		"sales_count":      len(salesTransactions),
		"payments_count":   len(leasePayments),
		"avg_sale_price":   totalSales / float64(len(salesTransactions)),
		"avg_rental":       totalRentals / float64(len(leasePayments)),
	}

	report := &models.FinancialReport{
		ReportType:       reportType,
		PeriodStart:      periodStart,
		PeriodEnd:        periodEnd,
		TotalSales:       totalSales,
		TotalRentals:     totalRentals,
		TotalRevenue:     totalRevenue,
		TotalCommissions: totalCommissions,
		TotalTaxes:       totalTaxes,
		TotalFees:        totalFees,
		ReportData:       utils.ToJSON(reportData),
		GeneratedAt:      time.Now(),
		GeneratedBy:      generatedBy,
	}

	if err := s.db.Create(report).Error; err != nil {
		return nil, err
	}

	return report, nil
}

// Helper functions
func (s *FinancialService) generatePromoCode() string {
	// Generate a random 8-character promo code
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 8)
	for i := range code {
		code[i] = chars[utils.RandomInt(0, len(chars)-1)]
	}
	return string(code)
}

// Input types
type CreateSaleTransactionInput struct {
	BuildingID    uuid.UUID                    `json:"building_id"`
	BuyerID       uuid.UUID                    `json:"buyer_id"`
	SellerID      uuid.UUID                    `json:"seller_id"`
	AgentID       *uuid.UUID                   `json:"agent_id"`
	Price         float64                      `json:"price"`
	PaymentMethod models.PaymentMethod         `json:"payment_method"`
	Commission    float64                      `json:"commission"`
	TaxAmount     float64                      `json:"tax_amount"`
	Fees          float64                      `json:"fees"`
	Notes         string                       `json:"notes"`
}

type CreateLeaseContractInput struct {
	PropertyID       uuid.UUID                    `json:"property_id"`
	TenantID         uuid.UUID                    `json:"tenant_id"`
	LandlordID       uuid.UUID                    `json:"landlord_id"`
	AgentID          *uuid.UUID                   `json:"agent_id"`
	DurationMonths   int                          `json:"duration_months"`
	StartDate        time.Time                    `json:"start_date"`
	MonthlyRent      float64                      `json:"monthly_rent"`
	DepositAmount    float64                      `json:"deposit_amount"`
	PaymentFrequency models.PaymentFrequency     `json:"payment_frequency"`
	UtilitiesIncluded bool                        `json:"utilities_included"`
	PetAllowed       bool                         `json:"pet_allowed"`
	Furnished        bool                         `json:"furnished"`
	LateFeeAmount    float64                      `json:"late_fee_amount"`
	GracePeriodDays  int                          `json:"grace_period_days"`
	Notes            string                       `json:"notes"`
}

type CreateLeasePaymentInput struct {
	LeaseID       uuid.UUID                    `json:"lease_id"`
	Amount        float64                      `json:"amount"`
	PaymentDate   time.Time                    `json:"payment_date"`
	DueDate       time.Time                    `json:"due_date"`
	PaymentMethod models.PaymentMethod         `json:"payment_method"`
	Notes         string                       `json:"notes"`
}

type CreateOfferInput struct {
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	DiscountPercent int        `json:"discount_percent"`
	DiscountAmount  float64    `json:"discount_amount"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	BuildingID      *uuid.UUID `json:"building_id"`
	CompanyID       *uuid.UUID `json:"company_id"`
	ImageURL        string     `json:"image_url"`
	TermsConditions string     `json:"terms_conditions"`
	MaxUses         int        `json:"max_uses"`
	MinAmount       float64    `json:"min_amount"`
	MaxAmount       float64    `json:"max_amount"`
	Code            string     `json:"code"`
}

type SaleTransactionFilters struct {
	BuildingID *uuid.UUID                    `json:"building_id"`
	BuyerID    *uuid.UUID                    `json:"buyer_id"`
	SellerID   *uuid.UUID                    `json:"seller_id"`
	AgentID    *uuid.UUID                    `json:"agent_id"`
	Status     models.TransactionStatus      `json:"status"`
	StartDate  *time.Time                    `json:"start_date"`
	EndDate    *time.Time                    `json:"end_date"`
}

type LeaseContractFilters struct {
	PropertyID *uuid.UUID           `json:"property_id"`
	TenantID   *uuid.UUID           `json:"tenant_id"`
	LandlordID *uuid.UUID           `json:"landlord_id"`
	Status     models.LeaseStatus   `json:"status"`
	ActiveOnly bool                 `json:"active_only"`
}

type OfferFilters struct {
	BuildingID *uuid.UUID `json:"building_id"`
	CompanyID  *uuid.UUID `json:"company_id"`
	ActiveOnly bool       `json:"active_only"`
} 