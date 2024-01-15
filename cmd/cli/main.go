package main

import (
	"finance/internal/database"
	balance "finance/internal/finance/balance/builder"
	"finance/internal/finance/clients"
	expense "finance/internal/finance/expense/builder"
	"finance/internal/finance/findate"
	findateb "finance/internal/finance/findate/builder"
	income "finance/internal/finance/income/builder"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"time"
)

func main() {

	db, err := database.NewDatabase()
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:  "Finance Plan",
		Usage: "Финансовое приложение для анализа, прогнозов, аналитики личных финансов",
		Commands: []cli.Command{
			{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "Управление клиентами",
				Subcommands: []cli.Command{
					{
						Name:  "list",
						Usage: "Получить всех клиентов",
						Action: func(c *cli.Context) error {
							clns := make([]database.Client, 0)
							err := db.Find(&clns).Error
							if err != nil {
								panic(err)
							}
							fmt.Println("Список всех клиентов:")

							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"ID", "Имя"})
							for _, cln := range clns {
								table.Append([]string{
									strconv.Itoa(cln.ID),
									cln.Name,
								})
							}
							table.Render()

							return nil
						},
					},
					{
						Name:  "new",
						Usage: "Создать нового клиента",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "name", Usage: "Имя клиента", Required: true},
						},
						Action: func(c *cli.Context) error {
							clientName := c.String("name")

							cln := &database.Client{
								Name: clientName,
							}
							err := db.Create(cln).Error
							if err != nil {
								panic(err)
							}

							fmt.Printf("Создан клиент c id %d\n", cln.ID)
							return nil
						},
					},
					{
						Name:  "rm",
						Usage: "Удалить клиента по id",
						Action: func(c *cli.Context) error {
							if len(c.Args()) == 0 {
								fmt.Println("Необходимо указать id клиента")
								return nil
							}
							idStr := c.Args().First()
							id, err := strconv.Atoi(idStr)
							if err != nil {
								fmt.Println("ID должен быть числом")
								return err
							}

							err = db.Delete(&database.Client{ID: id}).Error
							if err != nil {
								panic(err)
							}

							fmt.Printf("Клиент с id %d удален\n", id)
							return nil
						},
					},
				},
			},
			{
				Name:    "income",
				Aliases: []string{"i"},
				Usage:   "Управление доходами",
				Subcommands: []cli.Command{
					{
						Name:  "list",
						Usage: "Показывает доходы клиента",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "clientId", Usage: "ID клиента"},
						},
						Action: func(c *cli.Context) error {
							clientId := c.Int("clientId")

							incomes := make([]database.Income, 0)
							err := db.Where(&database.Income{ClientID: clientId}).Joins("FinDate").Find(&incomes).Error
							if err != nil {
								panic(err)
							}

							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"ID", "Название", "Описание", "Сумма", "Тип", "Период (d:m:y)"})
							for _, i := range incomes {
								d := i.FinDate
								var Type string
								var periodParams string

								switch d.Type {
								case findate.Once:
									Type = "Разовый"
								case findate.Periodic:
									Type = "Период"
									periodParams = fmt.Sprintf("%d:%d:%d", d.PeriodDays, d.PeriodMonths, d.PeriodYears)
								}

								table.Append([]string{
									strconv.Itoa(i.ID),
									i.Name,
									i.Description,
									strconv.FormatInt(i.Amount, 10),
									Type,
									periodParams,
								})
							}
							table.Render()

							return nil
						},
					},
					{
						Name:  "new",
						Usage: "Создать новый доход",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "clientId", Usage: "ID клиента", Required: true},

							&cli.StringFlag{Name: "name", Usage: "Название", Required: true},
							&cli.StringFlag{Name: "desc", Usage: "Описание"},
							&cli.Int64Flag{Name: "amount", Usage: "Сумма", Required: true},
							// Разовый доход
							&cli.StringFlag{Name: "dateAt", Usage: "Дата "},
							// Периодический доход
							&cli.BoolFlag{Name: "period", Usage: "Указать условия периодичности"},
							&cli.StringFlag{Name: "startAt", Usage: "Дата начала периода"},
							&cli.StringFlag{Name: "endAt", Usage: "Дата конеца периода"},
							&cli.StringFlag{Name: "periodAt", Usage: "Параметры периодичности. Формат: \"days:months:years\". Пример (раз в месяц): 0:1:0"},
						},
						Action: func(c *cli.Context) error {
							clientId := c.Int("clientId")
							Name := c.String("name")
							Description := c.String("desc")
							Amount := c.Int64("amount")

							startAtArgName := "dateAt"
							Type := findate.Once
							if c.Bool("period") {
								startAtArgName = "startAt"
								Type = findate.Periodic
							}

							var StartDay, StartMonth, StartYear int
							startAtStr := c.String(startAtArgName)
							if startAtStr == "" {
								date := time.Now()
								StartDay, StartMonth, StartYear = date.Day(), int(date.Month()), date.Year()
							} else if !findate.IsDate(startAtStr) {
								msg := fmt.Sprintf("%s должен быть формата dd.mm.yyyy", startAtArgName)
								panic(msg)
							} else {
								StartDay, StartMonth, StartYear = findate.ParseDate(startAtStr)
							}

							endAtStr := c.String("endAt")
							if endAtStr != "" && !findate.IsDate(endAtStr) {
								panic("endAt должен быть формата dd.mm.yyyy")
							}
							EndDay, EndMonth, EndYear := findate.ParseDatePtr(endAtStr)

							periodAtStr := c.String("periodAt")
							if periodAtStr != "" && !findate.IsPeriod(periodAtStr) {
								panic("periodAt должен быть формата \"%d:%d:%d\"")
							}
							PeriodDays, PeriodMonths, PeriodYears := findate.ParsePeriod(periodAtStr)
							if PeriodDays == 0 && PeriodMonths == 0 && PeriodYears == 0 {
								PeriodDays = 1
							}

							finDate := &database.FinDate{
								StartYear:    StartYear,
								StartMonth:   StartMonth,
								StartDay:     StartDay,
								EndYear:      EndYear,
								EndMonth:     EndMonth,
								EndDay:       EndDay,
								PeriodYears:  PeriodYears,
								PeriodMonths: PeriodMonths,
								PeriodDays:   PeriodDays,
								Type:         Type,
							}
							err := db.Create(finDate).Error
							if err != nil {
								panic(err)
							}

							cln := &database.Income{
								ClientID:    clientId,
								FinDateID:   finDate.ID,
								Amount:      Amount,
								Name:        Name,
								Description: Description,
							}
							err = db.Create(cln).Error
							if err != nil {
								panic(err)
							}

							fmt.Printf("Создан доход успешно id=%d\n", cln.ID)
							return nil
						},
					},
					{
						Name:  "rm",
						Usage: "Удалить доход по id",
						Action: func(c *cli.Context) error {
							if len(c.Args()) == 0 {
								fmt.Println("Необходимо указать id клиента")
								return nil
							}
							idStr := c.Args().First()
							id, err := strconv.Atoi(idStr)
							if err != nil {
								fmt.Println("ID должен быть числом")
								return err
							}

							err = db.Delete(&database.Income{ID: id}).Error
							if err != nil {
								panic(err)
							}

							fmt.Printf("Доход с id %d удален\n", id)
							return nil
						},
					},
				},
			},
			{
				Name:    "expense",
				Aliases: []string{"e"},
				Usage:   "Управление расходами",
				Subcommands: []cli.Command{
					{
						Name:  "list",
						Usage: "Показывает расходы клиента",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "clientId", Usage: "ID клиента"},
						},
						Action: func(c *cli.Context) error {
							clientId := c.Int("clientId")

							expenses := make([]database.Expense, 0)
							err := db.Where(&database.Expense{ClientID: clientId}).Joins("FinDate").Find(&expenses).Error
							if err != nil {
								panic(err)
							}

							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"ID", "Название", "Описание", "Сумма", "Тип", "Период (d:m:y)"})
							for _, i := range expenses {
								d := i.FinDate
								var Type string
								var periodParams string

								switch d.Type {
								case findate.Once:
									Type = "Разовый"
								case findate.Periodic:
									Type = "Период"
									periodParams = fmt.Sprintf("%d:%d:%d", d.PeriodDays, d.PeriodMonths, d.PeriodYears)
								}

								table.Append([]string{
									strconv.Itoa(i.ID),
									i.Name,
									i.Description,
									strconv.FormatInt(i.Amount, 10),
									Type,
									periodParams,
								})
							}
							table.Render()

							return nil
						},
					},
					{
						Name:  "new",
						Usage: "Создать новый расход",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "clientId", Usage: "ID клиента", Required: true},

							&cli.StringFlag{Name: "name", Usage: "Название", Required: true},
							&cli.StringFlag{Name: "desc", Usage: "Описание"},
							&cli.Int64Flag{Name: "amount", Usage: "Сумма", Required: true},
							// Разовый доход
							&cli.StringFlag{Name: "dateAt", Usage: "Дата "},
							// Периодический доход
							&cli.BoolFlag{Name: "period", Usage: "Указать условия периодичности"},
							&cli.StringFlag{Name: "startAt", Usage: "Дата начала периода"},
							&cli.StringFlag{Name: "endAt", Usage: "Дата конеца периода"},
							&cli.StringFlag{Name: "periodAt", Usage: "Параметры периодичности. Формат: \"days:months:years\". Пример (раз в месяц): 0:1:0"},
						},
						Action: func(c *cli.Context) error {
							clientId := c.Int("clientId")
							Name := c.String("name")
							Description := c.String("desc")
							Amount := c.Int64("amount")

							startAtArgName := "dateAt"
							Type := findate.Once
							if c.Bool("period") {
								startAtArgName = "startAt"
								Type = findate.Periodic
							}

							var StartDay, StartMonth, StartYear int
							startAtStr := c.String(startAtArgName)
							if startAtStr == "" {
								date := time.Now()
								StartDay, StartMonth, StartYear = date.Day(), int(date.Month()), date.Year()
							} else if !findate.IsDate(startAtStr) {
								msg := fmt.Sprintf("%s должен быть формата dd.mm.yyyy", startAtArgName)
								panic(msg)
							} else {
								StartDay, StartMonth, StartYear = findate.ParseDate(startAtStr)
							}

							endAtStr := c.String("endAt")
							if endAtStr != "" && !findate.IsDate(endAtStr) {
								panic("endAt должен быть формата dd.mm.yyyy")
							}
							EndDay, EndMonth, EndYear := findate.ParseDatePtr(endAtStr)

							periodAtStr := c.String("periodAt")
							if periodAtStr != "" && !findate.IsPeriod(periodAtStr) {
								panic("periodAt должен быть формата \"%d:%d:%d\"")
							}
							PeriodDays, PeriodMonths, PeriodYears := findate.ParsePeriod(periodAtStr)
							if PeriodDays == 0 && PeriodMonths == 0 && PeriodYears == 0 {
								PeriodDays = 1
							}

							finDate := &database.FinDate{
								StartYear:    StartYear,
								StartMonth:   StartMonth,
								StartDay:     StartDay,
								EndYear:      EndYear,
								EndMonth:     EndMonth,
								EndDay:       EndDay,
								PeriodYears:  PeriodYears,
								PeriodMonths: PeriodMonths,
								PeriodDays:   PeriodDays,
								Type:         Type,
							}
							err := db.Create(finDate).Error
							if err != nil {
								panic(err)
							}

							cln := &database.Expense{
								ClientID:    clientId,
								FinDateID:   finDate.ID,
								Amount:      Amount,
								Name:        Name,
								Description: Description,
							}
							err = db.Create(cln).Error
							if err != nil {
								panic(err)
							}

							fmt.Printf("Создан доход успешно id=%d\n", cln.ID)
							return nil
						},
					},
					{
						Name:  "rm",
						Usage: "Удалить расход по id",
						Action: func(c *cli.Context) error {
							if len(c.Args()) == 0 {
								fmt.Println("Необходимо указать id клиента")
								return nil
							}
							idStr := c.Args().First()
							id, err := strconv.Atoi(idStr)
							if err != nil {
								fmt.Println("ID должен быть числом")
								return err
							}

							err = db.Delete(&database.Expense{ID: id}).Error
							if err != nil {
								panic(err)
							}

							fmt.Printf("Расход с id %d удален\n", id)
							return nil
						},
					},
				},
			},

			{
				Name:    "balance",
				Aliases: []string{"c"},
				Usage:   "Баланс",
				Subcommands: []cli.Command{
					{
						Name:  "range",
						Usage: "",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "clientId", Usage: "ID клиента", Required: true},
							&cli.StringFlag{Name: "dateFrom", Usage: "Дата начала периода", Required: true},
							&cli.StringFlag{Name: "dateTo", Usage: "Дата конца периода", Required: true},
						},
						Action: func(c *cli.Context) error {
							clientId := c.Int("clientId")

							dateFrom := c.String("dateFrom")
							if !findate.IsDate(dateFrom) {
								panic("dateFrom должен быть формата dd.mm.yyyy")
							}

							dateTo := c.String("dateTo")
							if dateTo != "" && !findate.IsDate(dateTo) {
								panic("dateTo должен быть формата dd.mm.yyyy")
							}

							Client := &database.Client{ID: clientId}
							err = db.
								Preload("Incomes.FinDate").
								Preload("Expenses.FinDate").
								First(Client).Error
							if err != nil {
								fmt.Printf("Not found client with id = %d", clientId)
								return nil
							}
							fmt.Printf("%+v\n", Client)

							NewClientB := clients.NewClientBuilder
							NewIncomeB := income.NewIncomeBuilder
							NewExpenseB := expense.NewExpenseBuilder
							NewFinDateB := findateb.NewFinDateBuilder

							clientBuilder := NewClientB().Name(Client.Name)

							FromDbFindate := func(fd *database.FinDate) *findate.FinDate {
								fdb := NewFinDateB()
								fdb = fdb.StartAt(fd.StartDay, fd.StartMonth, fd.StartYear)

								switch fd.Type {
								case findate.Once:
									fdb = fdb.Once()
								case findate.Periodic:
									fdb = fdb.Periodic().SetPeriod(fd.PeriodDays, fd.PeriodMonths, fd.PeriodYears)
									if fd.EndDay != nil || fd.EndMonth != nil || fd.EndYear != nil {
										fdb = fdb.EndAt(*fd.EndDay, *fd.EndMonth, *fd.EndYear)
									}
								}
								return fdb.Build()
							}

							// ------ Incomes ------
							for _, i := range Client.Incomes {
								builder := NewIncomeB().
									Name(i.Name).Description(i.Description).
									Amount(i.Amount).FinDate(*FromDbFindate(&i.FinDate)).
									Build()
								clientBuilder.AddIncome(builder)
							}
							// ------ Expenses ------
							for _, i := range Client.Expenses {
								builder := NewExpenseB().
									Name(i.Name).Description(i.Description).
									Amount(i.Amount).FinDate(*FromDbFindate(&i.FinDate)).
									Build()
								clientBuilder.AddExpense(builder)
							}

							client := clientBuilder.Build()

							blc := balance.NewBalanceB(0).
								DateFrom(dateFrom).
								DateTo(dateTo).
								Expenses(client.BalanceExpenses()).
								Incomes(client.BalanceIncomes()).
								Build().
								Compute()

							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"Дата", "Доход", "Расход", "Баланс"})
							for _, i := range blc.RangeDays() {
								table.Append([]string{
									i.Date,
									strconv.FormatInt(i.Income, 10),
									strconv.FormatInt(i.Expense, 10),
									strconv.FormatInt(i.Amount, 10),
								})
							}
							table.Render()

							blc.ShowInfo()

							return nil
						},
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println("Ошибка при выполнении приложения:", err)
	}

}
