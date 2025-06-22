# Financial System Guide

## Overview

The financial system in the Go service provides comprehensive tracking of property sales, lease agreements, and all financial transactions between buyers, sellers, tenants, and agents. It includes promotional offers, payment tracking, and financial reporting capabilities.

## Core Entities

### 1. SaleTransaction
Tracks property sale transactions with detailed financial information.

**Key Features:**
- Complete transaction history
- Multiple payment methods support
- Commission and fee tracking
- Document management
- Status tracking (PENDING, COMPLETED, CANCELLED, FAILED, REFUNDED)

**Fields:**
- `id`: Unique transaction ID (UUID)
- `buildingId`: Reference to Building sold
- `buyerId`: Reference to User (buyer)
- `sellerId`: Reference to User (seller)
- `agentId`: Optional sales agent reference
- `price`: Final sale price
- `paymentMethod`: CASH, BANK_TRANSFER, MORTGAGE, CREDIT_CARD, CHECK, CRYPTO
- `status`: Transaction status
- `commission`: Agent commission amount
- `taxAmount`: Tax amount
- `fees`: Additional fees
- `totalAmount`: Total transaction amount
- `notes`: Transaction notes
- `createdAt`, `updatedAt`, `completedAt`: Timestamps

### 2. LeaseContract
Manages property lease agreements with comprehensive terms and conditions.

**Key Features:**
- Flexible lease terms
- Payment frequency options
- Late fee management
- Document storage
- Status tracking (ACTIVE, TERMINATED, PENDING, EXPIRED)

**Fields:**
- `id`: Unique lease ID (UUID)
- `propertyId`: Reference to property/unit
- `tenantId`: Reference to tenant (User)
- `landlordId`: Reference to landlord (User)
- `agentId`: Optional agent reference
- `durationMonths`: Lease duration in months
- `startDate`, `endDate`: Lease period
- `monthlyRent`: Monthly rent amount
- `depositAmount`: Security deposit
- `paymentFrequency`: MONTHLY, QUARTERLY, YEARLY, WEEKLY
- `status`: Lease status
- `contractFileUrl`, `contractFileKey`: Encrypted contract document
- `utilitiesIncluded`: Whether utilities are included
- `petAllowed`: Pet policy
- `furnished`: Furnishing status
- `lateFeeAmount`: Late payment fee
- `gracePeriodDays`: Payment grace period
- `notes`: Additional notes

### 3. LeasePayment
Tracks individual rent payments for lease contracts.

**Key Features:**
- Payment history tracking
- Late fee calculation
- Multiple payment methods
- Due date management

**Fields:**
- `id`: Unique payment ID (UUID)
- `leaseId`: Reference to lease contract
- `amount`: Payment amount
- `paymentDate`: When payment was made
- `dueDate`: When payment was due
- `status`: Payment status
- `paymentMethod`: Payment method used
- `lateFee`: Calculated late fee
- `notes`: Payment notes

### 4. Offer
Manages promotional offers for properties or companies.

**Key Features:**
- Percentage or fixed amount discounts
- Usage limits and tracking
- Time-based validity
- Building or company-specific offers
- Promo code generation

**Fields:**
- `id`: Unique offer ID (UUID)
- `title`: Offer title
- `description`: Full description
- `discountPercent`: Percentage discount
- `discountAmount`: Fixed amount discount
- `startDate`, `endDate`: Validity period
- `buildingId`: Optional building-specific offer
- `companyId`: Optional company-specific offer
- `imageUrl`: Offer banner image
- `termsConditions`: Offer terms
- `isActive`: Offer status
- `maxUses`: Usage limit (0 = unlimited)
- `currentUses`: Current usage count
- `minAmount`, `maxAmount`: Amount limits
- `code`: Promo code

### 5. OfferUse
Tracks when offers are used by users.

**Fields:**
- `id`: Unique usage ID (UUID)
- `offerId`: Reference to offer
- `userId`: User who used the offer
- `amount`: Transaction amount
- `discount`: Applied discount
- `usedAt`: Usage timestamp

### 6. FinancialReport
Generates comprehensive financial reports.

**Key Features:**
- Sales and rental revenue tracking
- Commission and fee summaries
- Period-based reporting
- Detailed breakdown data

**Fields:**
- `id`: Unique report ID (UUID)
- `reportType`: Report type (SALES, RENTALS, REVENUE)
- `periodStart`, `periodEnd`: Reporting period
- `totalSales`: Total sales revenue
- `totalRentals`: Total rental revenue
- `totalRevenue`: Combined revenue
- `totalCommissions`: Total commissions
- `totalTaxes`: Total taxes
- `totalFees`: Total fees
- `reportData`: Detailed breakdown (JSON)
- `generatedAt`: Report generation timestamp
- `generatedBy`: User who generated the report

## GraphQL API

### Queries

#### Get Sale Transactions
```graphql
query GetSaleTransactions(
  $buildingId: ID
  $buyerId: ID
  $sellerId: ID
  $agentId: ID
  $status: String
  $startDate: String
  $endDate: String
) {
  getSaleTransactions(
    buildingId: $buildingId
    buyerId: $buyerId
    sellerId: $sellerId
    agentId: $agentId
    status: $status
    startDate: $startDate
    endDate: $endDate
  ) {
    id
    buildingId
    buyerId
    sellerId
    agentId
    price
    paymentMethod
    status
    commission
    taxAmount
    fees
    totalAmount
    notes
    createdAt
    updatedAt
    completedAt
  }
}
```

#### Get Lease Contracts
```graphql
query GetLeaseContracts(
  $propertyId: ID
  $tenantId: ID
  $landlordId: ID
  $status: String
  $activeOnly: Boolean
) {
  getLeaseContracts(
    propertyId: $propertyId
    tenantId: $tenantId
    landlordId: $landlordId
    status: $status
    activeOnly: $activeOnly
  ) {
    id
    propertyId
    tenantId
    landlordId
    agentId
    durationMonths
    startDate
    endDate
    monthlyRent
    depositAmount
    paymentFrequency
    status
    utilitiesIncluded
    petAllowed
    furnished
    lateFeeAmount
    gracePeriodDays
    notes
    createdAt
    updatedAt
    signedAt
    terminatedAt
  }
}
```

#### Get Offers
```graphql
query GetOffers(
  $buildingId: ID
  $companyId: ID
  $activeOnly: Boolean
) {
  getOffers(
    buildingId: $buildingId
    companyId: $companyId
    activeOnly: $activeOnly
  ) {
    id
    title
    description
    discountPercent
    discountAmount
    startDate
    endDate
    buildingId
    companyId
    imageUrl
    termsConditions
    isActive
    maxUses
    currentUses
    minAmount
    maxAmount
    code
    createdAt
    updatedAt
  }
}
```

### Mutations

#### Create Sale Transaction
```graphql
mutation CreateSaleTransaction($input: CreateSaleTransactionInput!) {
  createSaleTransaction(input: $input) {
    id
    buildingId
    buyerId
    sellerId
    agentId
    price
    paymentMethod
    status
    commission
    taxAmount
    fees
    totalAmount
    notes
    createdAt
    updatedAt
  }
}
```

**Input Example:**
```json
{
  "buildingId": "123e4567-e89b-12d3-a456-426614174000",
  "buyerId": "123e4567-e89b-12d3-a456-426614174001",
  "sellerId": "123e4567-e89b-12d3-a456-426614174002",
  "agentId": "123e4567-e89b-12d3-a456-426614174003",
  "price": 500000.00,
  "paymentMethod": "BANK_TRANSFER",
  "commission": 15000.00,
  "taxAmount": 25000.00,
  "fees": 5000.00,
  "notes": "Property sale transaction"
}
```

#### Create Lease Contract
```graphql
mutation CreateLeaseContract($input: CreateLeaseContractInput!) {
  createLeaseContract(input: $input) {
    id
    propertyId
    tenantId
    landlordId
    agentId
    durationMonths
    startDate
    endDate
    monthlyRent
    depositAmount
    paymentFrequency
    status
    utilitiesIncluded
    petAllowed
    furnished
    lateFeeAmount
    gracePeriodDays
    notes
    createdAt
    updatedAt
  }
}
```

**Input Example:**
```json
{
  "propertyId": "123e4567-e89b-12d3-a456-426614174000",
  "tenantId": "123e4567-e89b-12d3-a456-426614174001",
  "landlordId": "123e4567-e89b-12d3-a456-426614174002",
  "agentId": "123e4567-e89b-12d3-a456-426614174003",
  "durationMonths": 12,
  "startDate": "2024-01-01T00:00:00Z",
  "monthlyRent": 2500.00,
  "depositAmount": 5000.00,
  "paymentFrequency": "MONTHLY",
  "utilitiesIncluded": true,
  "petAllowed": false,
  "furnished": true,
  "lateFeeAmount": 100.00,
  "gracePeriodDays": 5,
  "notes": "Standard lease agreement"
}
```

#### Create Lease Payment
```graphql
mutation CreateLeasePayment($input: CreateLeasePaymentInput!) {
  createLeasePayment(input: $input) {
    id
    leaseId
    amount
    paymentDate
    dueDate
    status
    paymentMethod
    lateFee
    notes
    createdAt
    updatedAt
  }
}
```

**Input Example:**
```json
{
  "leaseId": "123e4567-e89b-12d3-a456-426614174000",
  "amount": 2500.00,
  "paymentDate": "2024-01-01T00:00:00Z",
  "dueDate": "2024-01-01T00:00:00Z",
  "paymentMethod": "BANK_TRANSFER",
  "notes": "January rent payment"
}
```

#### Create Offer
```graphql
mutation CreateOffer($input: CreateOfferInput!) {
  createOffer(input: $input) {
    id
    title
    description
    discountPercent
    discountAmount
    startDate
    endDate
    buildingId
    companyId
    imageUrl
    termsConditions
    isActive
    maxUses
    currentUses
    minAmount
    maxAmount
    code
    createdAt
    updatedAt
  }
}
```

**Input Example:**
```json
{
  "title": "Summer Sale Discount",
  "description": "Get 10% off on all properties this summer",
  "discountPercent": 10,
  "startDate": "2024-06-01T00:00:00Z",
  "endDate": "2024-08-31T23:59:59Z",
  "companyId": "123e4567-e89b-12d3-a456-426614174000",
  "imageUrl": "https://example.com/summer-sale.jpg",
  "termsConditions": "Valid for new contracts only",
  "maxUses": 100,
  "minAmount": 1000.00,
  "maxAmount": 100000.00
}
```

#### Use Offer
```graphql
mutation UseOffer($code: String!, $userId: ID!, $amount: Float!) {
  useOffer(code: $code, userId: $userId, amount: $amount) {
    id
    offerId
    userId
    amount
    discount
    usedAt
  }
}
```

#### Update Transaction Status
```graphql
mutation UpdateSaleTransactionStatus($id: ID!, $status: String!) {
  updateSaleTransactionStatus(id: $id, status: $status) {
    id
    status
    updatedAt
    completedAt
  }
}
```

#### Update Lease Status
```graphql
mutation UpdateLeaseStatus($id: ID!, $status: String!) {
  updateLeaseStatus(id: $id, status: $status) {
    id
    status
    updatedAt
    signedAt
    terminatedAt
  }
}
```

#### Generate Financial Report
```graphql
mutation GenerateFinancialReport(
  $reportType: String!
  $periodStart: String!
  $periodEnd: String!
  $generatedBy: ID!
) {
  generateFinancialReport(
    reportType: $reportType
    periodStart: $periodStart
    periodEnd: $periodEnd
    generatedBy: $generatedBy
  ) {
    id
    reportType
    periodStart
    periodEnd
    totalSales
    totalRentals
    totalRevenue
    totalCommissions
    totalTaxes
    totalFees
    reportData
    generatedAt
    generatedBy
  }
}
```

## Business Logic

### Payment Methods
- **CASH**: Cash payment
- **BANK_TRANSFER**: Bank transfer
- **MORTGAGE**: Mortgage payment
- **CREDIT_CARD**: Credit card payment
- **CHECK**: Check payment
- **CRYPTO**: Cryptocurrency payment

### Transaction Statuses
- **PENDING**: Transaction initiated but not completed
- **COMPLETED**: Transaction successfully completed
- **CANCELLED**: Transaction cancelled
- **FAILED**: Transaction failed
- **REFUNDED**: Transaction refunded

### Lease Statuses
- **PENDING**: Lease contract pending approval
- **ACTIVE**: Active lease contract
- **TERMINATED**: Lease contract terminated
- **EXPIRED**: Lease contract expired

### Payment Frequencies
- **MONTHLY**: Monthly payments
- **QUARTERLY**: Quarterly payments
- **YEARLY**: Yearly payments
- **WEEKLY**: Weekly payments

## Features

### 1. Automatic Calculations
- Total amount calculation (price + commission + tax + fees)
- Late fee calculation based on grace period
- End date calculation based on duration
- Discount calculation (percentage or fixed amount)

### 2. Document Management
- Encrypted file storage for contracts and documents
- File key management for secure access
- Support for multiple document types

### 3. Offer System
- Automatic promo code generation
- Usage tracking and limits
- Time-based validity
- Amount-based restrictions

### 4. Financial Reporting
- Comprehensive revenue tracking
- Period-based reporting
- Detailed breakdown data
- Commission and fee summaries

### 5. Status Management
- Complete transaction lifecycle tracking
- Lease status management
- Payment status tracking
- Automatic timestamp updates

## Security Features

### 1. Encryption
- File encryption for sensitive documents
- Secure key management
- Encrypted storage for contract files

### 2. Access Control
- User-based access to financial data
- Role-based permissions
- Audit trail for all transactions

### 3. Data Validation
- Input validation for all financial data
- Business rule enforcement
- Error handling and logging

## Integration Points

### 1. User Management
- Integration with user system for buyers, sellers, tenants, landlords
- User role management
- User-based filtering and reporting

### 2. Property Management
- Integration with building/property system
- Property-specific offers
- Property-based transaction tracking

### 3. Company Management
- Company-specific offers
- Company-based reporting
- Multi-tenant support

### 4. Notification System
- Payment due notifications
- Transaction status updates
- Offer expiration alerts

## Best Practices

### 1. Data Management
- Always validate financial data before processing
- Use appropriate decimal precision for monetary values
- Implement proper error handling for failed transactions

### 2. Security
- Encrypt sensitive financial documents
- Implement proper access controls
- Maintain audit trails for all transactions

### 3. Performance
- Use database indexes for frequently queried fields
- Implement pagination for large result sets
- Cache frequently accessed data

### 4. Monitoring
- Monitor transaction volumes and patterns
- Track offer usage and effectiveness
- Monitor system performance and errors

## Error Handling

The system includes comprehensive error handling for:
- Invalid input data
- Database connection issues
- File upload failures
- Business rule violations
- Authentication and authorization errors

All errors are logged with appropriate context and returned to the client with meaningful error messages. 