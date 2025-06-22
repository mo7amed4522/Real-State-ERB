package graphql

import (
	"context"
	"fmt"
	"my-property/go-service/models"
	"my-property/go-service/services"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

// Financial Resolvers
type FinancialResolvers struct {
	financialService *services.FinancialService
}

func NewFinancialResolvers(financialService *services.FinancialService) *FinancialResolvers {
	return &FinancialResolvers{
		financialService: financialService,
	}
}

// Sale Transaction Resolvers
func (r *FinancialResolvers) CreateSaleTransaction(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args
	inputMap, ok := args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid 'input' argument")
	}

	var input services.CreateSaleTransactionInput

	// Parse input from map
	if buildingID, ok := inputMap["buildingId"].(string); ok {
		if id, err := uuid.Parse(buildingID); err == nil {
			input.BuildingID = id
		}
	}
	if buyerID, ok := inputMap["buyerId"].(string); ok {
		if id, err := uuid.Parse(buyerID); err == nil {
			input.BuyerID = id
		}
	}
	if sellerID, ok := inputMap["sellerId"].(string); ok {
		if id, err := uuid.Parse(sellerID); err == nil {
			input.SellerID = id
		}
	}
	if agentID, ok := inputMap["agentId"].(string); ok && agentID != "" {
		if id, err := uuid.Parse(agentID); err == nil {
			input.AgentID = &id
		}
	}
	if price, ok := inputMap["price"].(float64); ok {
		input.Price = price
	}
	if paymentMethod, ok := inputMap["paymentMethod"].(string); ok {
		input.PaymentMethod = models.PaymentMethod(paymentMethod)
	}
	if commission, ok := inputMap["commission"].(float64); ok {
		input.Commission = commission
	}
	if taxAmount, ok := inputMap["taxAmount"].(float64); ok {
		input.TaxAmount = taxAmount
	}
	if fees, ok := inputMap["fees"].(float64); ok {
		input.Fees = fees
	}
	if notes, ok := inputMap["notes"].(string); ok {
		input.Notes = notes
	}

	transaction, err := r.financialService.CreateSaleTransaction(input)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *FinancialResolvers) UpdateSaleTransactionStatus(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args

	transactionIDStr, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'id' argument")
	}
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'id' argument: %w", err)
	}

	statusStr, ok := args["status"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'status' argument")
	}
	status := models.TransactionStatus(statusStr)

	transaction, err := r.financialService.UpdateSaleTransactionStatus(transactionID, status)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *FinancialResolvers) GetSaleTransactions(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args

	var filters services.SaleTransactionFilters

	if buildingID, ok := args["buildingId"].(string); ok && buildingID != "" {
		if id, err := uuid.Parse(buildingID); err == nil {
			filters.BuildingID = &id
		}
	}
	if buyerID, ok := args["buyerId"].(string); ok && buyerID != "" {
		if id, err := uuid.Parse(buyerID); err == nil {
			filters.BuyerID = &id
		}
	}
	if sellerID, ok := args["sellerId"].(string); ok && sellerID != "" {
		if id, err := uuid.Parse(sellerID); err == nil {
			filters.SellerID = &id
		}
	}
	if agentID, ok := args["agentId"].(string); ok && agentID != "" {
		if id, err := uuid.Parse(agentID); err == nil {
			filters.AgentID = &id
		}
	}
	if status, ok := args["status"].(string); ok && status != "" {
		filters.Status = models.TransactionStatus(status)
	}
	if startDate, ok := args["startDate"].(string); ok && startDate != "" {
		if date, err := time.Parse(time.RFC3339, startDate); err == nil {
			filters.StartDate = &date
		}
	}
	if endDate, ok := args["endDate"].(string); ok && endDate != "" {
		if date, err := time.Parse(time.RFC3339, endDate); err == nil {
			filters.EndDate = &date
		}
	}

	transactions, err := r.financialService.GetSaleTransactions(filters)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// Lease Contract Resolvers
func (r *FinancialResolvers) CreateLeaseContract(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args
	inputMap, ok := args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid 'input' argument")
	}

	var input services.CreateLeaseContractInput

	if propertyID, ok := inputMap["propertyId"].(string); ok {
		if id, err := uuid.Parse(propertyID); err == nil {
			input.PropertyID = id
		}
	}
	if tenantID, ok := inputMap["tenantId"].(string); ok {
		if id, err := uuid.Parse(tenantID); err == nil {
			input.TenantID = id
		}
	}
	if landlordID, ok := inputMap["landlordId"].(string); ok {
		if id, err := uuid.Parse(landlordID); err == nil {
			input.LandlordID = id
		}
	}
	if agentID, ok := inputMap["agentId"].(string); ok && agentID != "" {
		if id, err := uuid.Parse(agentID); err == nil {
			input.AgentID = &id
		}
	}
	if durationMonths, ok := inputMap["durationMonths"].(float64); ok { // JSON numbers are float64
		input.DurationMonths = int(durationMonths)
	}
	if startDate, ok := inputMap["startDate"].(string); ok {
		if date, err := time.Parse(time.RFC3339, startDate); err == nil {
			input.StartDate = date
		}
	}
	if monthlyRent, ok := inputMap["monthlyRent"].(float64); ok {
		input.MonthlyRent = monthlyRent
	}
	if depositAmount, ok := inputMap["depositAmount"].(float64); ok {
		input.DepositAmount = depositAmount
	}
	if paymentFrequency, ok := inputMap["paymentFrequency"].(string); ok {
		input.PaymentFrequency = models.PaymentFrequency(paymentFrequency)
	}
	if utilitiesIncluded, ok := inputMap["utilitiesIncluded"].(bool); ok {
		input.UtilitiesIncluded = utilitiesIncluded
	}
	if petAllowed, ok := inputMap["petAllowed"].(bool); ok {
		input.PetAllowed = petAllowed
	}
	if furnished, ok := inputMap["furnished"].(bool); ok {
		input.Furnished = furnished
	}
	if lateFeeAmount, ok := inputMap["lateFeeAmount"].(float64); ok {
		input.LateFeeAmount = lateFeeAmount
	}
	if gracePeriodDays, ok := inputMap["gracePeriodDays"].(float64); ok { // JSON numbers are float64
		input.GracePeriodDays = int(gracePeriodDays)
	}
	if notes, ok := inputMap["notes"].(string); ok {
		input.Notes = notes
	}

	lease, err := r.financialService.CreateLeaseContract(input)
	if err != nil {
		return nil, err
	}

	return lease, nil
}

func (r *FinancialResolvers) UpdateLeaseStatus(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args

	leaseIDStr, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'id' argument")
	}
	leaseID, err := uuid.Parse(leaseIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'id' argument: %w", err)
	}

	statusStr, ok := args["status"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'status' argument")
	}
	status := models.LeaseStatus(statusStr)

	lease, err := r.financialService.UpdateLeaseStatus(leaseID, status)
	if err != nil {
		return nil, err
	}

	return lease, nil
}

func (r *FinancialResolvers) GetLeaseContracts(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args

	var filters services.LeaseContractFilters

	if propertyID, ok := args["propertyId"].(string); ok && propertyID != "" {
		if id, err := uuid.Parse(propertyID); err == nil {
			filters.PropertyID = &id
		}
	}
	if tenantID, ok := args["tenantId"].(string); ok && tenantID != "" {
		if id, err := uuid.Parse(tenantID); err == nil {
			filters.TenantID = &id
		}
	}
	if landlordID, ok := args["landlordId"].(string); ok && landlordID != "" {
		if id, err := uuid.Parse(landlordID); err == nil {
			filters.LandlordID = &id
		}
	}
	if status, ok := args["status"].(string); ok && status != "" {
		filters.Status = models.LeaseStatus(status)
	}
	if activeOnly, ok := args["activeOnly"].(bool); ok {
		filters.ActiveOnly = activeOnly
	}

	leases, err := r.financialService.GetLeaseContracts(filters)
	if err != nil {
		return nil, err
	}

	return leases, nil
}

// Lease Payment Resolvers
func (r *FinancialResolvers) CreateLeasePayment(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args
	inputMap, ok := args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid 'input' argument")
	}

	var input services.CreateLeasePaymentInput

	if leaseID, ok := inputMap["leaseId"].(string); ok {
		if id, err := uuid.Parse(leaseID); err == nil {
			input.LeaseID = id
		}
	}
	if amount, ok := inputMap["amount"].(float64); ok {
		input.Amount = amount
	}
	if paymentDate, ok := inputMap["paymentDate"].(string); ok {
		if date, err := time.Parse(time.RFC3339, paymentDate); err == nil {
			input.PaymentDate = date
		}
	}
	if dueDate, ok := inputMap["dueDate"].(string); ok {
		if date, err := time.Parse(time.RFC3339, dueDate); err == nil {
			input.DueDate = date
		}
	}
	if paymentMethod, ok := inputMap["paymentMethod"].(string); ok {
		input.PaymentMethod = models.PaymentMethod(paymentMethod)
	}
	if notes, ok := inputMap["notes"].(string); ok {
		input.Notes = notes
	}

	payment, err := r.financialService.CreateLeasePayment(input)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *FinancialResolvers) UpdatePaymentStatus(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args

	paymentIDStr, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'id' argument")
	}
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'id' argument: %w", err)
	}

	statusStr, ok := args["status"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'status' argument")
	}
	status := models.TransactionStatus(statusStr)

	payment, err := r.financialService.UpdatePaymentStatus(paymentID, status)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// Offer Resolvers
func (r *FinancialResolvers) CreateOffer(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args
	inputMap, ok := args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid 'input' argument")
	}

	var input services.CreateOfferInput

	if title, ok := inputMap["title"].(string); ok {
		input.Title = title
	}
	if description, ok := inputMap["description"].(string); ok {
		input.Description = description
	}
	if discountPercent, ok := inputMap["discountPercent"].(float64); ok { // JSON numbers are float64
		input.DiscountPercent = int(discountPercent)
	}
	if discountAmount, ok := inputMap["discountAmount"].(float64); ok {
		input.DiscountAmount = discountAmount
	}
	if startDate, ok := inputMap["startDate"].(string); ok {
		if date, err := time.Parse(time.RFC3339, startDate); err == nil {
			input.StartDate = date
		}
	}
	if endDate, ok := inputMap["endDate"].(string); ok {
		if date, err := time.Parse(time.RFC3339, endDate); err == nil {
			input.EndDate = date
		}
	}
	if buildingID, ok := inputMap["buildingId"].(string); ok && buildingID != "" {
		if id, err := uuid.Parse(buildingID); err == nil {
			input.BuildingID = &id
		}
	}
	if companyID, ok := inputMap["companyId"].(string); ok && companyID != "" {
		if id, err := uuid.Parse(companyID); err == nil {
			input.CompanyID = &id
		}
	}
	if imageURL, ok := inputMap["imageUrl"].(string); ok {
		input.ImageURL = imageURL
	}
	if termsConditions, ok := inputMap["termsConditions"].(string); ok {
		input.TermsConditions = termsConditions
	}
	if maxUses, ok := inputMap["maxUses"].(float64); ok { // JSON numbers are float64
		input.MaxUses = int(maxUses)
	}
	if minAmount, ok := inputMap["minAmount"].(float64); ok {
		input.MinAmount = minAmount
	}
	if maxAmount, ok := inputMap["maxAmount"].(float64); ok {
		input.MaxAmount = maxAmount
	}
	if code, ok := inputMap["code"].(string); ok {
		input.Code = code
	}

	offer, err := r.financialService.CreateOffer(input)
	if err != nil {
		return nil, err
	}

	return offer, nil
}

func (r *FinancialResolvers) UseOffer(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args

	offerCode, ok := args["code"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'code' argument")
	}
	userIDStr, ok := args["userId"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'userId' argument")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'userId' argument: %w", err)
	}
	amount, ok := args["amount"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing 'amount' argument")
	}

	offerUse, err := r.financialService.UseOffer(offerCode, userID, amount)
	if err != nil {
		return nil, err
	}

	return offerUse, nil
}

func (r *FinancialResolvers) GetOffers(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args

	var filters services.OfferFilters

	if buildingID, ok := args["buildingId"].(string); ok && buildingID != "" {
		if id, err := uuid.Parse(buildingID); err == nil {
			filters.BuildingID = &id
		}
	}
	if companyID, ok := args["companyId"].(string); ok && companyID != "" {
		if id, err := uuid.Parse(companyID); err == nil {
			filters.CompanyID = &id
		}
	}
	if activeOnly, ok := args["activeOnly"].(bool); ok {
		filters.ActiveOnly = activeOnly
	}

	offers, err := r.financialService.GetOffers(filters)
	if err != nil {
		return nil, err
	}

	return offers, nil
}

// Financial Report Resolvers
func (r *FinancialResolvers) GenerateFinancialReport(ctx context.Context) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args

	reportType, ok := args["reportType"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'reportType' argument")
	}
	periodStartStr, ok := args["periodStart"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'periodStart' argument")
	}
	periodStart, err := time.Parse(time.RFC3339, periodStartStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'periodStart' argument: %w", err)
	}
	periodEndStr, ok := args["periodEnd"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'periodEnd' argument")
	}
	periodEnd, err := time.Parse(time.RFC3339, periodEndStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'periodEnd' argument: %w", err)
	}
	generatedByStr, ok := args["generatedBy"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'generatedBy' argument")
	}
	generatedBy, err := uuid.Parse(generatedByStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'generatedBy' argument: %w", err)
	}

	report, err := r.financialService.GenerateFinancialReport(reportType, periodStart, periodEnd, generatedBy)
	if err != nil {
		return nil, err
	}

	return report, nil
}
