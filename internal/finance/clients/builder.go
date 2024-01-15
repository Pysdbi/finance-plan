package clients

import (
	"finance/internal/finance/expense"
	"finance/internal/finance/income"
)

type clientBuilder struct {
	name string

	incomes  []income.Income
	expenses []expense.Expense
}

func NewClientBuilder() IClientBuilder {
	return &clientBuilder{}
}

type IClientBuilder interface {
	Name(v string) IClientBuilder

	AddIncome(income *income.Income) IClientBuilder

	AddExpense(expense *expense.Expense) IClientBuilder

	Build() *Client
}

func (b *clientBuilder) Name(v string) IClientBuilder {
	b.name = v
	return b
}

func (b *clientBuilder) AddIncome(income *income.Income) IClientBuilder {
	b.incomes = append(b.incomes, *income)
	return b
}

func (b *clientBuilder) AddExpense(expense *expense.Expense) IClientBuilder {
	b.expenses = append(b.expenses, *expense)
	return b
}

func (b *clientBuilder) Build() *Client {
	return &Client{
		Name:     b.name,
		Incomes:  b.incomes,
		Expenses: b.expenses,
	}
}
