package main

import (
	balance "finance/internal/finance/balance/builder"
	"finance/internal/finance/clients"
	expense "finance/internal/finance/expense/builder"
	findateb "finance/internal/finance/findate/builder"
	income "finance/internal/finance/income/builder"
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()

	NewClientB := clients.NewClientBuilder
	NewIncomeB := income.NewIncomeBuilder
	NewExpenseB := expense.NewExpenseBuilder
	NewFinDateB := findateb.NewFinDateBuilder

	client := NewClientB().
		Name("Дмитрий").
		// ------ Incomes ------
		AddIncome(NewIncomeB().
			Name("Зарплата").Amount(100000).
			FinDate(*NewFinDateB().StartAt(10, 1, 2024).EveryMonth().Build()).Build(),
		).
		AddIncome(NewIncomeB().
			Name("Дивиденты").Amount(50000).
			FinDate(*NewFinDateB().StartAt(1, 1, 2024).EveryMonth().Build()).Build(),
		).
		AddIncome(NewIncomeB().
			Name("Проценты").Amount(5000).
			FinDate(*NewFinDateB().StartAt(4, 1, 2024).EveryDay().Build()).Build(),
		).
		// ------ Expenses ------
		AddExpense(NewExpenseB().
			Name("Ипотека").Amount(60000).
			FinDate(*NewFinDateB().StartAt(5, 1, 2024).Every14Days().Build()).Build(),
		).
		Build()

	blc := balance.NewBalanceB(0).
		DateFrom("01.01.2024").
		DateTo("31.01.2025").
		Expenses(client.BalanceExpenses()).
		Incomes(client.BalanceIncomes()).
		Build().
		Compute()

	defer fmt.Printf("Время выполнения: %.3f ms\n", float64(time.Since(startTime).Nanoseconds())/1000000.0)

	blc.ShowInfo()
	//blc.ShowDetailsTable()
}
