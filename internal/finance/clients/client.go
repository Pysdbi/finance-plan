package clients

import (
	"finance/internal/finance/balance"
	"finance/internal/finance/expense"
	"finance/internal/finance/income"
)

// ========= Аккаунт =========

type Client struct {
	Name     string
	Incomes  []income.Income
	Expenses []expense.Expense
}

func (a *Client) BalanceIncomes() []balance.TransactionIncome {
	var incomes []balance.TransactionIncome
	for _, i := range a.Incomes {
		incomes = append(incomes, balance.TransactionIncome{
			Amount:  i.Amount,
			FinDate: i.FinDate,
		})
	}
	return incomes
}

func (a *Client) BalanceExpenses() []balance.TransactionExpense {
	var expenses []balance.TransactionExpense
	for _, i := range a.Expenses {
		expenses = append(expenses, balance.TransactionExpense{
			Amount:  i.Amount,
			FinDate: i.FinDate,
		})
	}
	return expenses
}
