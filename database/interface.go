package database

import "context"

// TransactionExecutor defines the interface for executing database transactions.
type TransactionExecutor interface {
	ExecuteInTransaction(ctx context.Context, fn func(context.Context) error) error
}
