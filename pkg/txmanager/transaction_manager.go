package txmanager

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoTransactionInContext = errors.New("no transaction found in context")
)

// define a private type for context key to avoid collisions.
type contextKey string

// TxContextKey is the key used to store the transaction (*gorm.DB) in a context.Context.
// Repositories can retrieve the transaction from context to ensure operations are executed
// within the same database transaction.
const txContextKey contextKey = "do-transaction"

// TransactionManager defines the interface for managing database transactions.
type TransactionManager interface {
	// Do executes the given job function within a database transaction.
	// The transaction is automatically committed if the job returns nil,
	// or rolled back if the job returns an error or panics.
	//
	// The transaction (*gorm.DB) is stored in the context with key TxContextKey,
	// so that repositories can retrieve and use it to perform operations within the same transaction.
	Do(ctx context.Context, job func(ctx context.Context) error) error
}

// txManager is a concrete implementation of TransactionManager using GORM.
type txManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new TransactionManager with the given GORM DB connection.
func NewTransactionManager(db *gorm.DB) *txManager {
	return &txManager{db: db}
}

// Do executes the job function within a transaction, injecting the transaction into the context.
func (txm *txManager) Do(ctx context.Context, job func(ctx context.Context) error) error {
	return txm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Inject the transaction into a new context so repositories can access it.
		txnContext := context.WithValue(ctx, txContextKey, tx)

		// Execute the job. If it returns an error, GORM rolls back automatically.
		return job(txnContext)
	})
}

// GetTx retrieves the current GORM transaction (*gorm.DB) from the given context.
// It returns nil if no transaction is found in the context.
//
// This function is intended to be used inside repository methods so that all
// database operations are executed within the same transaction managed by TransactionManager.
//
// Example usage in a repository:
//
//	func (r *repo) CustomerRepository(ctx context.Context, customerID string) (*Customer, error) {
//	    tx, err := GetTx(ctx)
//	    if err != nil {
//	        return nil, err
//	    }
//	    // use tx to query database
//	}
func GetTx(ctx context.Context) (*gorm.DB, error) {
	tx, ok := ctx.Value(txContextKey).(*gorm.DB)
	if !ok {
		return nil, ErrNoTransactionInContext
	}

	return tx, nil
}
