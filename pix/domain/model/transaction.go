package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionError     = "error"
	TransactionPending   = "pending"
	TransactionConfirmed = "confirmed"
	TransactionCompleted = "completed"
)

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transactions []*Transaction
}

type Transaction struct {
	Base              `valid:"required"`
	Status            string   `json:"status" valid:"notnull"`
	Amount            float64  `json:"amount" valid:"notnull"`
	Description       string   `json:"description" valid:"notnull"`
	CancelDescription string   `json:"cancel_description" valid:"notnull"`
	AccountFrom       *Account `valid:"-"`
	PixKeyTo          *PixKey  `valid:"-"`
}

func (t *Transaction) validate() error {
	_, err := govalidator.ValidateStruct(t)

	if t.Amount <= 0 {
		return errors.New("the amount must be greater than 0")
	}

	if t.Status != TransactionCompleted && t.Status != TransactionConfirmed && t.Status != TransactionPending && t.Status != TransactionError {
		return errors.New("invalid transaction status")
	}

	if t.PixKeyTo.AccountID == t.AccountFrom.ID {
		return errors.New("the source and destination can't be the same")
	}

	return err
}

func NewTransaction(accountFrom *Account, pixKey *PixKey, amount float64, description string) (*Transaction, error) {
	transaction := Transaction{
		Amount:            amount,
		Description:       description,
		CancelDescription: "",
		PixKeyTo:          pixKey,
		AccountFrom:       accountFrom,
		Status:            TransactionPending,
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	if err := transaction.validate(); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()

	return t.validate()
}

func (t *Transaction) Confirm() error {
	t.Status = TransactionConfirmed
	t.UpdatedAt = time.Now()

	return t.validate()
}

func (t *Transaction) Error(description string) error {
	t.Status = TransactionError
	t.CancelDescription = description
	t.UpdatedAt = time.Now()

	return t.validate()
}
