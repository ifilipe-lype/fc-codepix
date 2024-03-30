package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type PixKeyRepositoryInterface interface {
	RegisterKey(pixkey *PixKey) (*PixKey, error)
	FindKeyByKind(key, kind string) (*PixKey, error)
	AddBank(bank *Bank) error
	AddAccount(account *Account) error
	FindAccount(id string) (*Account, error)
}

type PixKey struct {
	Base      `valid:"required"`
	Kind      string    `json:"kind" valid:"notnull"`
	Key       string    `json:"key" valid:"notnull"`
	Status    string    `json:"status" valid:"notnull"`
	AccountID string    `json:"account_id" valid:"notnull"`
	Account   *Account  `valid:"-"`
	PixKeys   []*PixKey `valid:"-"`
}

func (pixKey *PixKey) validate() error {
	_, err := govalidator.ValidateStruct(pixKey)

	if pixKey.Kind != "email" && pixKey.Kind != "cpf" {
		return errors.New("invalid pixkey kind")
	}

	if pixKey.Status != "active" && pixKey.Status != "inactive" {
		return errors.New("invalid pixkey status")
	}

	return err
}

func NewPixKey(account *Account, key, kind string) (*PixKey, error) {
	pixKey := PixKey{
		Key:       key,
		Kind:      kind,
		Account:   account,
		AccountID: account.ID,
		Status:    "active",
	}

	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()

	if err := pixKey.validate(); err != nil {
		return nil, err
	}

	return &pixKey, nil
}
