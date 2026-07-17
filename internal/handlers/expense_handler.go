package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/services"
	"github.com/afdhalpower/golanglaundry/internal/helpers"
)

type ExpenseHandler struct {
	expenseService  *services.ExpenseService
	categoryService *services.ExpenseCategoryService
}

func NewExpenseHandler(expenseService *services.ExpenseService, categoryService *services.ExpenseCategoryService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService:  expenseService,
		categoryService: categoryService,
	}
}

func (h *ExpenseHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	categoryID := c.Query("category_id", "")
	search := c.Query("search", "")

	expenses, total, err := h.expenseService.GetAll(page, limit, categoryID, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat pengeluaran")
	}

	categories, _ := h.categoryService.GetAll()

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return render(c, "expenses/index", fiber.Map{
		"title":      "Pengeluaran",
		"expenses":   expenses,
		"categories": categories,
		"categoryID": categoryID,
		"search":     search,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
	}, "layouts/main")
}

func (h *ExpenseHandler) New(c fiber.Ctx) error {
	categories, _ := h.categoryService.GetAll()
	return render(c, "expenses/form", fiber.Map{
		"title":      "Tambah Pengeluaran",
		"expense":    nil,
		"categories": categories,
	}, "layouts/main")
}

func (h *ExpenseHandler) Create(c fiber.Ctx) error {
	categoryID, _ := strconv.ParseUint(c.FormValue("expense_category_id"), 10, 32)
	amount, _ := strconv.ParseFloat(c.FormValue("amount"), 64)
	dateStr := c.FormValue("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = time.Now()
	}

	expense := &models.Expense{
		ExpenseCategoryID: uint(categoryID),
		Amount:            amount,
		Description:       c.FormValue("description"),
		Date:              date,
		CreatedBy:         helpers.LogAndGetUserID(c),
	}

	if err := h.expenseService.Create(expense); err != nil {
		categories, _ := h.categoryService.GetAll()
		return render(c, "expenses/form", fiber.Map{
			"title":      "Tambah Pengeluaran",
			"expense":    expense,
			"categories": categories,
			"error":      "Gagal menyimpan pengeluaran",
		}, "layouts/main")
	}

	return c.Redirect().To("/expenses")
}

func (h *ExpenseHandler) Edit(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	expense, err := h.expenseService.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pengeluaran tidak ditemukan")
	}

	categories, _ := h.categoryService.GetAll()
	return render(c, "expenses/form", fiber.Map{
		"title":      "Edit Pengeluaran",
		"expense":    expense,
		"categories": categories,
	}, "layouts/main")
}

func (h *ExpenseHandler) Update(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	expense, err := h.expenseService.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pengeluaran tidak ditemukan")
	}

	categoryID, _ := strconv.ParseUint(c.FormValue("expense_category_id"), 10, 32)
	amount, _ := strconv.ParseFloat(c.FormValue("amount"), 64)
	dateStr := c.FormValue("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = expense.Date
	}

	expense.ExpenseCategoryID = uint(categoryID)
	expense.Amount = amount
	expense.Description = c.FormValue("description")
	expense.Date = date

	if err := h.expenseService.Update(expense); err != nil {
		categories, _ := h.categoryService.GetAll()
		return render(c, "expenses/form", fiber.Map{
			"title":      "Edit Pengeluaran",
			"expense":    expense,
			"categories": categories,
			"error":      "Gagal memperbarui pengeluaran",
		}, "layouts/main")
	}

	return c.Redirect().To("/expenses")
}

func (h *ExpenseHandler) Delete(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	h.expenseService.Delete(uint(id))
	return c.Redirect().To("/expenses")
}
