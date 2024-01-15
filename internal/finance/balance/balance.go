package balance

import (
	helpers "finance/internal"
	"finance/internal/finance/findate"
	"fmt"
	"log"
	"time"
)

type Balance struct {
	// Amount Чистый доход за указанный промежуток времени (Amount = Income - Expense)
	Amount int64
	// Income Доход за указанный промежуток времени
	Income int64
	// Expense Расход за указанный промежуток времени
	Expense int64

	// From Дата начала выборки
	DateFrom time.Time
	// To Дата конца выборки
	DateTo time.Time

	Incomes  []TransactionIncome
	Expenses []TransactionExpense

	rangeDays  []balanceDayResult
	isComputed bool
}

type balanceDayResult struct {
	findate.DateRangeItem
	Income            int64
	Expense           int64
	DiffIncomeExpense int64
	Amount            int64
}

func (b *Balance) Compute() *Balance {
	for _, d := range findate.DefineRange(b.DateFrom, b.DateTo) {
		var totalIncome, totalExpense int64
		for _, item := range b.Incomes {
			if item.FinDate.InDate(d.Time) {
				totalIncome += item.Amount
			}
		}
		for _, item := range b.Expenses {
			if item.FinDate.InDate(d.Time) {
				totalExpense += item.Amount
			}
		}
		b.Income += totalIncome
		b.Expense += totalExpense
		b.Amount += totalIncome - totalExpense

		b.rangeDays = append(b.rangeDays, *&balanceDayResult{
			DateRangeItem:     d,
			Income:            totalIncome,
			Expense:           totalExpense,
			DiffIncomeExpense: totalIncome - totalExpense,
			Amount:            b.Amount,
		})
	}
	b.isComputed = true
	return b
}

func (b *Balance) RangeDays() []balanceDayResult {
	return b.rangeDays
}

func (b *Balance) ShowDetailsTable() {
	if !b.isComputed {
		log.Printf("Выборка не просчитана. Выполните метод Balance.Compute")
		return
	}

	headers := []string{"Дата", "Доход", "Расход", "Разница Д-Р", "Баланс"}
	tableData := make([][]string, 0)

	for _, i := range b.RangeDays() {
		tableData = append(tableData, []string{
			i.Date,
			helpers.FormatNumber(i.Income),
			helpers.FormatNumber(i.Expense),
			helpers.FormatNumber(i.DiffIncomeExpense),
			helpers.FormatNumber(i.Amount),
		})
	}

	helpers.PrintTable(headers, tableData)
}

func (b *Balance) ShowInfo() {
	if !b.isComputed {
		log.Printf("Выборка не просчитана. Выполните метод Balance.Compute")
		return
	}
	fmt.Printf(
		"==================\nЗа весь период:\n  Период: %s - %s\n  Доход: %v\n  Расход: %v\n  Баланс: %v\n==================\n",
		b.DateFrom.Format("02.01.2006"),
		b.DateTo.Format("02.01.2006"),
		helpers.FormatNumber(b.Income),
		helpers.FormatNumber(b.Expense),
		helpers.FormatNumber(b.Amount),
	)
}
