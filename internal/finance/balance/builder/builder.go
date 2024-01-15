package balance

import (
	"finance/internal/finance/balance"
	"time"
)

type balanceBuilder struct {
	amount  int64
	income  int64
	expense int64

	dateFrom time.Time
	dateTo   time.Time

	incomes  []balance.TransactionIncome
	expenses []balance.TransactionExpense
}

func NewBalanceB(amount int64) IBalanceBuilder {
	return &balanceBuilder{
		amount: amount,
	}
}

type IBalanceBuilder interface {
	DateFrom(v string) IBalanceBuilder
	DateTo(v string) IBalanceBuilder

	Incomes(v []balance.TransactionIncome) IBalanceBuilder
	Expenses(v []balance.TransactionExpense) IBalanceBuilder

	Build() *balance.Balance
}

const format = "02.01.2006"

// DateFrom Указание даты начала выборки. Строка в формате "02.01.2006"
func (b *balanceBuilder) DateFrom(v string) IBalanceBuilder {
	date, _ := time.Parse(format, v)
	b.dateFrom = date
	return b
}

// DateTo Указание даты конца выборки. Строка в формате "02.01.2006"
func (b *balanceBuilder) DateTo(v string) IBalanceBuilder {
	date, _ := time.Parse(format, v)
	b.dateTo = date
	return b
}

func (b *balanceBuilder) Incomes(v []balance.TransactionIncome) IBalanceBuilder {
	b.incomes = v
	return b
}

func (b *balanceBuilder) Expenses(v []balance.TransactionExpense) IBalanceBuilder {
	b.expenses = v
	return b
}

func (b *balanceBuilder) Build() *balance.Balance {
	return &balance.Balance{
		DateFrom: b.dateFrom,
		DateTo:   b.dateTo,
		Incomes:  b.incomes,
		Expenses: b.expenses,
		Amount:   b.amount,
	}
}
