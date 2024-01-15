package expense

import (
	"finance/internal/finance/expense"
	"finance/internal/finance/findate"
)

type expenseBuilder struct {
	finDate     findate.FinDate
	amount      int64
	name        string
	description string
}

type IExpenseBuilder interface {
	Name(v string) IExpenseBuilder
	Description(v string) IExpenseBuilder
	// Amount сумма поступления
	Amount(v int64) IExpenseBuilder

	FinDate(v findate.FinDate) IExpenseBuilder

	Build() *expense.Expense
}

func NewExpenseBuilder() IExpenseBuilder {
	return &expenseBuilder{}
}

func (b *expenseBuilder) Amount(v int64) IExpenseBuilder {
	b.amount = v
	return b
}

func (b *expenseBuilder) Name(v string) IExpenseBuilder {
	b.name = v
	return b
}

func (b *expenseBuilder) Description(v string) IExpenseBuilder {
	b.description = v
	return b
}

func (b *expenseBuilder) FinDate(v findate.FinDate) IExpenseBuilder {
	b.finDate = v
	return b
}

func (b *expenseBuilder) Build() *expense.Expense {
	return &expense.Expense{
		Amount:      b.amount,
		FinDate:     b.finDate,
		Name:        b.name,
		Description: b.description,
	}
}
