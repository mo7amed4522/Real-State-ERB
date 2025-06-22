# Secure Wallet System

This wallet system provides a secure, fast, and reliable way to handle financial transactions with the following security features:

## Security Features

### 1. **RedLock Distributed Locking**
- Uses Redis-based distributed locking to prevent race conditions
- Ensures atomic operations across multiple instances
- Prevents double-spending and concurrent transaction conflicts

### 2. **Database Transactions**
- All wallet operations use database transactions for atomicity
- Ensures data consistency even if operations fail
- Rollback capability for failed transactions

### 3. **AES-256-GCM Encryption**
- All sensitive data is encrypted using AES-256-GCM
- Encrypted fields include:
  - Bank details
  - Account names
  - Bank IDs
  - Sender/Receiver details
  - Payment information

### 4. **Stripe Integration**
- Secure payment processing with Stripe
- Webhook handling for payment confirmations
- Support for multiple payment methods

### 5. **Authentication & Authorization**
- JWT-based authentication
- Role-based access control
- User-specific wallet access

## Environment Variables Required

```bash
# Stripe Configuration
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key_here
STRIPE_WEBHOOK_SECRET=whsec_your_stripe_webhook_secret_here

# Redis Configuration
REDIS_URL=redis://redis:6379

# Encryption
ENCRYPTION_SECRET_KEY=your_32_character_encryption_key_here
```

## GraphQL Operations

### Create Wallet
```graphql
mutation CreateWallet($input: CreateWalletInput!) {
  createWallet(input: $input) {
    id
    userId
    balance
    isActive
    createdAt
  }
}
```

### Deposit Funds
```graphql
mutation Deposit($walletId: String!, $input: DepositInput!) {
  deposit(walletId: $walletId, input: $input) {
    id
    amount
    status
    stripePaymentIntentId
    metadata
  }
}
```

### Withdraw Funds
```graphql
mutation Withdraw($walletId: String!, $input: WithdrawInput!) {
  withdraw(walletId: $walletId, input: $input) {
    id
    amount
    status
    stripeTransferId
  }
}
```

### Transfer Between Wallets
```graphql
mutation Transfer($fromWalletId: String!, $input: TransferInput!) {
  transfer(fromWalletId: $fromWalletId, input: $input) {
    id
    amount
    status
    metadata
  }
}
```

### Get Transaction History
```graphql
query GetTransactionHistory($walletId: String!, $limit: Int, $offset: Int) {
  getTransactionHistory(walletId: $walletId, limit: $limit, offset: $offset) {
    id
    amount
    transactionType
    status
    createdAt
    description
  }
}
```

## Security Best Practices

### 1. **Data Encryption**
- All sensitive data is encrypted before storage
- Encryption keys are stored securely in environment variables
- No sensitive data is logged or exposed in error messages

### 2. **Transaction Security**
- All transactions are locked using RedLock
- Database transactions ensure atomicity
- Proper error handling and rollback mechanisms

### 3. **Payment Security**
- Stripe handles PCI compliance
- Webhook signatures are verified
- Payment intents are used for secure payment processing

### 4. **Access Control**
- User authentication required for all operations
- Users can only access their own wallets
- Role-based permissions for admin operations

## Performance Optimizations

### 1. **Database Indexing**
- Indexes on frequently queried fields
- Composite indexes for complex queries
- Proper foreign key relationships

### 2. **Caching**
- Redis for distributed locking
- Session storage in Redis
- Query result caching where appropriate

### 3. **Connection Pooling**
- Database connection pooling
- Redis connection management
- Stripe API connection optimization

## Error Handling

The system includes comprehensive error handling for:
- Insufficient funds
- Network failures
- Stripe API errors
- Database transaction failures
- RedLock acquisition failures
- Encryption/decryption errors

## Monitoring & Logging

- Structured logging for all operations
- Error tracking and alerting
- Performance metrics collection
- Audit trail for all transactions

## Testing

The wallet system should be thoroughly tested for:
- Concurrent transaction handling
- Network failure scenarios
- Stripe webhook processing
- Encryption/decryption accuracy
- Database transaction rollbacks
- RedLock behavior under load 