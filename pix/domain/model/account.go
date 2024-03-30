package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Account struct {
	Base      `valid:"required"`
	OwnerName string `json:"owner_name" valid:"notnull"`
	Number    int    `json:"number" valid:"notnull"`
	Bank      *Bank  `valid:"-"`
}

func (account *Account) validate() error {
	_, err := govalidator.ValidateStruct(account)

	return err
}

func NewAccount(bank *Bank, ownerName string, number int) (*Account, error) {
	account := Account{
		OwnerName: ownerName,
		Number:    number,
		Bank:      bank,
	}

	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()

	if err := account.validate(); err != nil {
		return nil, err
	}

	return &account, nil
}
