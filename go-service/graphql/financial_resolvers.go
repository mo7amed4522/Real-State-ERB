package graphql

import (
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
func (r *FinancialResolvers) CreateSaleTransaction(ctx graphql.ResolveContext) (interface{}, error) {
	var input services.CreateSaleTransactionInput
	if err := graphql.GetFieldContext(ctx).Args["input"].(map[string]interface{}); err != nil {
		return nil, err
	}

	// Parse input
	if buildingID, ok := input["buildingId"].(string); ok {
		if id, err := uuid.Parse(buildingID); err == nil {
			input.BuildingID = id
		}
	}
	if buyerID, ok := input["buyerId"].(string); ok {
		if id, err := uuid.Parse(buyerID); err == nil {
			input.BuyerID = id
		}
	}
	if sellerID, ok := input["sellerId"].(string); ok {
		if id, err := uuid.Parse(sellerID); err == nil {
			input.SellerID = id
		}
	}
	if agentID, ok := input["agentId"].(string); ok && agentID != "" {
		if id, err := uuid.Parse(agentID); err == nil {
			input.AgentID = &id
		}
	}
	if price, ok := input["price"].(float64); ok {
		input.Price = price
	}
	if paymentMethod, ok := input["paymentMethod"].(string); ok {
		input.PaymentMethod = models.PaymentMethod(paymentMethod)
	}
	if commission, ok := input["commission"].(float64); ok {
		input.Commission = commission
	}
	if taxAmount, ok := input["taxAmount"].(float64); ok {
		input.TaxAmount = taxAmount
	}
	if fees, ok := input["fees"].(float64); ok {
		input.Fees = fees
	}
	if notes, ok := input["notes"].(string); ok {
		input.Notes = notes
	}

	transaction, err := r.financialService.CreateSaleTransaction(input)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *FinancialResolvers) UpdateSaleTransactionStatus(ctx graphql.ResolveContext) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args
	
	transactionID, err := uuid.Parse(args["id"].(string))
	if err != nil {
		return nil, err
	}

	status := models.TransactionStatus(args["status"].(string))

	transaction, err := r.financialService.UpdateSaleTransactionStatus(transactionID, status)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *FinancialResolvers) GetSaleTransactions(ctx graphql.ResolveContext) (interface{}, error) {
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
func (r *FinancialResolvers) CreateLeaseContract(ctx graphql.ResolveContext) (interface{}, error) {
	var input services.CreateLeaseContractInput
	if err := graphql.GetFieldContext(ctx).Args["input"].(map[string]interface{}); err != nil {
		return nil, err
	}

	// Parse input
	if propertyID, ok := input["propertyId"].(string); ok {
		if id, err := uuid.Parse(propertyID); err == nil {
			input.PropertyID = id
		}
	}
	if tenantID, ok := input["tenantId"].(string); ok {
		if id, err := uuid.Parse(tenantID); err == nil {
			input.TenantID = id
		}
	}
	if landlordID, ok := input["landlordId"].(string); ok {
		if id, err := uuid.Parse(landlordID); err == nil {
			input.LandlordID = id
		}
	}
	if agentID, ok := input["agentId"].(string); ok && agentID != "" {
		if id, err := uuid.Parse(agentID); err == nil {
			input.AgentID = &id
		}
	}
	if durationMonths, ok := input["durationMonths"].(int); ok {
		input.DurationMonths = durationMonths
	}
	if startDate, ok := input["startDate"].(string); ok {
		if date, err := time.Parse(time.RFC3339, startDate); err == nil {
			input.StartDate = date
		}
	}
	if monthlyRent, ok := input["monthlyRent"].(float64); ok {
		input.MonthlyRent = monthlyRent
	}
	if depositAmount, ok := input["depositAmount"].(float64); ok {
		input.DepositAmount = depositAmount
	}
	if paymentFrequency, ok := input["paymentFrequency"].(string); ok {
		input.PaymentFrequency = models.PaymentFrequency(paymentFrequency)
	}
	if utilitiesIncluded, ok := input["utilitiesIncluded"].(bool); ok {
		input.UtilitiesIncluded = utilitiesIncluded
	}
	if petAllowed, ok := input["petAllowed"].(bool); ok {
		input.PetAllowed = petAllowed
	}
	if furnished, ok := input["furnished"].(bool); ok {
		input.Furnished = furnished
	}
	if lateFeeAmount, ok := input["lateFeeAmount"].(float64); ok {
		input.LateFeeAmount = lateFeeAmount
	}
	if gracePeriodDays, ok := input["gracePeriodDays"].(int); ok {
		input.GracePeriodDays = gracePeriodDays
	}
	if notes, ok := input["notes"].(string); ok {
		input.Notes = notes
	}

	lease, err := r.financialService.CreateLeaseContract(input)
	if err != nil {
		return nil, err
	}

	return lease, nil
}

func (r *FinancialResolvers) UpdateLeaseStatus(ctx graphql.ResolveContext) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args
	
	leaseID, err := uuid.Parse(args["id"].(string))
	if err != nil {
		return nil, err
	}

	status := models.LeaseStatus(args["status"].(string))

	lease, err := r.financialService.UpdateLeaseStatus(leaseID, status)
	if err != nil {
		return nil, err
	}

	return lease, nil
}

func (r *FinancialResolvers) GetLeaseContracts(ctx graphql.ResolveContext) (interface{}, error) {
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
func (r *FinancialResolvers) CreateLeasePayment(ctx graphql.ResolveContext) (interface{}, error) {
	var input services.CreateLeasePaymentInput
	if err := graphql.GetFieldContext(ctx).Args["input"].(map[string]interface{}); err != nil {
		return nil, err
	}

	// Parse input
	if leaseID, ok := input["leaseId"].(string); ok {
		if id, err := uuid.Parse(leaseID); err == nil {
			input.LeaseID = id
		}
	}
	if amount, ok := input["amount"].(float64); ok {
		input.Amount = amount
	}
	if paymentDate, ok := input["paymentDate"].(string); ok {
		if date, err := time.Parse(time.RFC3339, paymentDate); err == nil {
			input.PaymentDate = date
		}
	}
	if dueDate, ok := input["dueDate"].(string); ok {
		if date, err := time.Parse(time.RFC3339, dueDate); err == nil {
			input.DueDate = date
		}
	}
	if paymentMethod, ok := input["paymentMethod"].(string); ok {
		input.PaymentMethod = models.PaymentMethod(paymentMethod)
	}
	if notes, ok := input["notes"].(string); ok {
		input.Notes = notes
	}

	payment, err := r.financialService.CreateLeasePayment(input)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *FinancialResolvers) UpdatePaymentStatus(ctx graphql.ResolveContext) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args
	
	paymentID, err := uuid.Parse(args["id"].(string))
	if err != nil {
		return nil, err
	}

	status := models.TransactionStatus(args["status"].(string))

	payment, err := r.financialService.UpdatePaymentStatus(paymentID, status)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// Offer Resolvers
func (r *FinancialResolvers) CreateOffer(ctx graphql.ResolveContext) (interface{}, error) {
	var input services.CreateOfferInput
	if err := graphql.GetFieldContext(ctx).Args["input"].(map[string]interface{}); err != nil {
		return nil, err
	}

	// Parse input
	if title, ok := input["title"].(string); ok {
		input.Title = title
	}
	if description, ok := input["description"].(string); ok {
		input.Description = description
	}
	if discountPercent, ok := input["discountPercent"].(int); ok {
		input.DiscountPercent = discountPercent
	}
	if discountAmount, ok := input["discountAmount"].(float64); ok {
		input.DiscountAmount = discountAmount
	}
	if startDate, ok := input["startDate"].(string); ok {
		if date, err := time.Parse(time.RFC3339, startDate); err == nil {
			input.StartDate = date
		}
	}
	if endDate, ok := input["endDate"].(string); ok {
		if date, err := time.Parse(time.RFC3339, endDate); err == nil {
			input.EndDate = date
		}
	}
	if buildingID, ok := input["buildingId"].(string); ok && buildingID != "" {
		if id, err := uuid.Parse(buildingID); err == nil {
			input.BuildingID = &id
		}
	}
	if companyID, ok := input["companyId"].(string); ok && companyID != "" {
		if id, err := uuid.Parse(companyID); err == nil {
			input.CompanyID = &id
		}
	}
	if imageURL, ok := input["imageUrl"].(string); ok {
		input.ImageURL = imageURL
	}
	if termsConditions, ok := input["termsConditions"].(string); ok {
		input.TermsConditions = termsConditions
	}
	if maxUses, ok := input["maxUses"].(int); ok {
		input.MaxUses = maxUses
	}
	if minAmount, ok := input["minAmount"].(float64); ok {
		input.MinAmount = minAmount
	}
	if maxAmount, ok := input["maxAmount"].(float64); ok {
		input.MaxAmount = maxAmount
	}
	if code, ok := input["code"].(string); ok {
		input.Code = code
	}

	offer, err := r.financialService.CreateOffer(input)
	if err != nil {
		return nil, err
	}

	return offer, nil
}

func (r *FinancialResolvers) UseOffer(ctx graphql.ResolveContext) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args
	
	offerCode := args["code"].(string)
	userID, err := uuid.Parse(args["userId"].(string))
	if err != nil {
		return nil, err
	}
	amount := args["amount"].(float64)

	offerUse, err := r.financialService.UseOffer(offerCode, userID, amount)
	if err != nil {
		return nil, err
	}

	return offerUse, nil
}

func (r *FinancialResolvers) GetOffers(ctx graphql.ResolveContext) (interface{}, error) {
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
func (r *FinancialResolvers) GenerateFinancialReport(ctx graphql.ResolveContext) (interface{}, error) {
	args := graphql.GetFieldContext(ctx).Args
	
	reportType := args["reportType"].(string)
	periodStart, err := time.Parse(time.RFC3339, args["periodStart"].(string))
	if err != nil {
		return nil, err
	}
	periodEnd, err := time.Parse(time.RFC3339, args["periodEnd"].(string))
	if err != nil {
		return nil, err
	}
	generatedBy, err := uuid.Parse(args["generatedBy"].(string))
	if err != nil {
		return nil, err
	}

	report, err := r.financialService.GenerateFinancialReport(reportType, periodStart, periodEnd, generatedBy)
	if err != nil {
		return nil, err
	}

	return report, nil
} 