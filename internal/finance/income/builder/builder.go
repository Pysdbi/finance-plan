package income

import (
	"finance/internal/finance/findate"
	"finance/internal/finance/income"
)

type incomeBuilder struct {
	finDate     findate.FinDate
	amount      int64
	name        string
	description string
}

type IIncomeBuilder interface {
	Name(v string) IIncomeBuilder
	Description(v string) IIncomeBuilder
	// Amount сумма поступления
	Amount(v int64) IIncomeBuilder

	FinDate(v findate.FinDate) IIncomeBuilder

	Build() *income.Income
}

func NewIncomeBuilder() IIncomeBuilder {
	return &incomeBuilder{}
}

func (b *incomeBuilder) Name(v string) IIncomeBuilder {
	b.name = v
	return b
}

func (b *incomeBuilder) Description(v string) IIncomeBuilder {
	b.description = v
	return b
}

func (b *incomeBuilder) Amount(v int64) IIncomeBuilder {
	b.amount = v
	return b
}

func (b *incomeBuilder) FinDate(v findate.FinDate) IIncomeBuilder {
	b.finDate = v
	return b
}

func (b *incomeBuilder) Build() *income.Income {
	return &income.Income{
		Amount:      b.amount,
		FinDate:     b.finDate,
		Name:        b.name,
		Description: b.description,
	}
}
