# Laundry Enhancement Features Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add 5 missing features to the Laundry Management System: Invoice Printing, Payment on Order Page, Dashboard Charts Fix, Search All Pages, Stock Movement Log.

**Architecture:** Each feature is independent and can be developed in parallel. Features touch the model/repository/service/handler/template stack consistently.

**Tech Stack:** Go 1.25, Fiber v3, GORM, PostgreSQL, Alpine.js, Chart.js, HTML/CSS (Tailwind)

## Global Constraints

- All UI must match existing design system (slate-900 sidebar, rounded-xl cards, indigo-600 buttons, status badges with colored dots)
- Indonesian language for all UI text
- Use `render(c, ...)` helper for template rendering (already wired in all handlers)
- Use Fiber v3 `session.FromContext(c)` for session access
- Follow existing Clean Architecture pattern: models → repositories → services → handlers → templates
- All new models must be registered in `internal/config/database.go` AutoMigrate
- All new routes in `internal/routes/routes.go`
- Module path: `github.com/afdhalpower/golanglaundry`

---

### Task 1: Nota/Invoice Printing (Cetak Nota)

**Files:**
- Create: `templates/orders/print.html`
- Modify: `internal/handlers/order_handler.go` (add Print handler)
- Modify: `internal/routes/routes.go` (add print route)

**Interfaces:**
- Consumes: `OrderService.GetByID(id)`, `models.Order` with `Customer`, `Details`, `Details.Service`, `User`, `Payment` preloaded
- Produces: Print handler at `GET /orders/:id/print`

- [ ] **Step 1: Create print template**

Create `templates/orders/print.html` — a clean, printable invoice layout:

```html
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Nota - {{.order.OrderNumber}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        @media print { @page { margin: 1.5cm; } body { -webkit-print-color-adjust: exact; print-color-adjust: exact; } }
    </style>
</head>
<body class="bg-white">
    <div class="max-w-lg mx-auto p-6">
        <!-- Header -->
        <div class="text-center border-b-2 border-slate-200 pb-4 mb-4">
            <h1 class="text-xl font-bold text-slate-900">{{.laundryName}}</h1>
            <p class="text-xs text-slate-500">{{.laundryAddress}}</p>
            <p class="text-xs text-slate-500">Telp: {{.laundryPhone}}</p>
            <h2 class="text-lg font-bold mt-2 text-indigo-600">NOTA LAUNDRY</h2>
        </div>

        <!-- Order Info -->
        <div class="text-sm mb-4">
            <div class="flex justify-between"><span class="text-slate-500">No. Nota:</span><span class="font-mono font-medium">{{.order.OrderNumber}}</span></div>
            <div class="flex justify-between"><span class="text-slate-500">Pelanggan:</span><span class="font-medium">{{.order.Customer.Name}}</span></div>
            <div class="flex justify-between"><span class="text-slate-500">Tanggal:</span><span>{{.order.EntryDate.Format "02 Jan 2006 15:04"}}</span></div>
            <div class="flex justify-between"><span class="text-slate-500">Estimasi Selesai:</span><span>{{.order.EstimatedDoneDate.Format "02 Jan 2006 15:04"}}</span></div>
            <div class="flex justify-between"><span class="text-slate-500">Status:</span>
                <span class="font-medium {{if eq .order.Status "sudah_diambil"}}text-emerald-600{{else}}text-amber-600{{end}}">
                    {{if eq .order.Status "menunggu"}}Menunggu{{else if eq .order.Status "dicuci"}}Dicuci{{else if eq .order.Status "dikeringkan"}}Dikeringkan{{else if eq .order.Status "disetrika"}}Disetrika{{else if eq .order.Status "siap_diambil"}}Siap Diambil{{else if eq .order.Status "sudah_diambil"}}Sudah Diambil{{else if eq .order.Status "dibatalkan"}}Dibatalkan{{end}}
                </span>
            </div>
        </div>

        <!-- Service Details -->
        <table class="w-full text-sm mb-4">
            <thead>
                <tr class="bg-slate-100">
                    <th class="text-left px-3 py-2 font-semibold text-slate-700">Layanan</th>
                    <th class="text-center px-3 py-2 font-semibold text-slate-700">Berat</th>
                    <th class="text-right px-3 py-2 font-semibold text-slate-700">Harga</th>
                    <th class="text-right px-3 py-2 font-semibold text-slate-700">Subtotal</th>
                </tr>
            </thead>
            <tbody>
                {{range .order.Details}}
                <tr class="border-b border-slate-100">
                    <td class="px-3 py-2">{{if .Service}}{{.Service.Name}}{{else}}Layanan #{{.ServiceID}}{{end}}</td>
                    <td class="px-3 py-2 text-center">{{$.order.WeightKg}} kg</td>
                    <td class="px-3 py-2 text-right">Rp {{printf "%.0f" .PricePerKg}}</td>
                    <td class="px-3 py-2 text-right font-medium">Rp {{printf "%.0f" (mul $.order.WeightKg .PricePerKg)}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>

        <!-- Total -->
        <div class="border-t-2 border-slate-300 pt-3 text-sm space-y-1">
            <div class="flex justify-between"><span>Subtotal</span><span>Rp {{printf "%.0f" (mul .order.WeightKg (index .order.Details 0).PricePerKg)}}</span></div>
            {{if gt .order.Discount 0.0}}<div class="flex justify-between text-red-600"><span>Diskon</span><span>-Rp {{printf "%.0f" .order.Discount}}</span></div>{{end}}
            {{if gt .order.ExtraCost 0.0}}<div class="flex justify-between text-amber-600"><span>Biaya Tambahan</span><span>+Rp {{printf "%.0f" .order.ExtraCost}}</span></div>{{end}}
            <div class="flex justify-between text-base font-bold text-indigo-600 pt-2 border-t border-slate-200">
                <span>TOTAL</span>
                <span>Rp {{printf "%.0f" .order.Total}}</span>
            </div>
        </div>

        <!-- Payment -->
        {{if .payment}}
        <div class="mt-4 pt-4 border-t border-slate-200 text-sm">
            <div class="flex justify-between"><span class="text-slate-500">Pembayaran:</span><span class="font-medium text-emerald-600">LUNAS</span></div>
            <div class="flex justify-between"><span class="text-slate-500">Metode:</span><span>{{if eq .payment.Method "tunai"}}Tunai{{else if eq .payment.Method "qris"}}QRIS{{else if eq .payment.Method "transfer"}}Transfer{{else}}{{.payment.Method}}{{end}}</span></div>
            <div class="flex justify-between"><span class="text-slate-500">Dibayar:</span><span>Rp {{printf "%.0f" .payment.Amount}}</span></div>
        </div>
        {{end}}

        <!-- Footer -->
        <div class="mt-6 pt-4 border-t border-slate-200 text-center text-xs text-slate-400">
            <p>Terima kasih telah menggunakan jasa kami</p>
            <p>{{.laundryName}} | {{.laundryPhone}}</p>
        </div>
    </div>
    <script>window.print();</script>
</body>
</html>
```

- [ ] **Step 2: Add Print handler in order_handler.go**

Add method `Print` to `OrderHandler`:

```go
func (h *OrderHandler) Print(c fiber.Ctx) error {
    id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
    order, err := h.orderService.GetByID(uint(id))
    if err != nil {
        return c.Status(fiber.StatusNotFound).SendString("Pesanan tidak ditemukan")
    }

    // Get settings
    settingSvc := h.settingService
    laundryName, _ := settingSvc.Get("laundry_name")
    if laundryName == "" {
        laundryName = "Laundry Management"
    }
    laundryAddress, _ := settingSvc.Get("laundry_address")
    laundryPhone, _ := settingSvc.Get("laundry_phone")

    return c.Render("orders/print", fiber.Map{
        "order":         order,
        "laundryName":   laundryName,
        "laundryAddress": laundryAddress,
        "laundryPhone":  laundryPhone,
    })
}
```

Note: Need to add `settingService *services.SettingService` to `OrderHandler` struct and constructor. Add `settingSvc` parameter to `NewOrderHandler`.

- [ ] **Step 3: Update OrderHandler constructor**

Change `NewOrderHandler` to accept `settingService`, store it in the struct, and add `SettingService` to `OrderHandler` struct.

- [ ] **Step 4: Wire route in routes.go**

Add: `orders.Get("/:id/print", orderHandler.Print)`

- [ ] **Step 5: Add print button to order show page**

In `templates/orders/show.html`, add a print button next to the "Update Status" card:

```html
<a href="/orders/{{.order.ID}}/print" target="_blank"
   class="block w-full bg-indigo-600 hover:bg-indigo-700 text-white py-2 rounded-lg text-sm font-medium transition text-center mb-3">
   <span class="flex items-center justify-center gap-2">
       <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"></path></svg>
       Cetak Nota
   </span>
</a>
```

- [ ] **Step 6: Build + Test**

Run `go build ./...`, start server, login, open order detail, click "Cetak Nota" button, verify print page appears with correct data.

---

### Task 2: Payment on Order Page

**Files:**
- Modify: `templates/orders/show.html` (add payment card)
- Modify: `internal/handlers/order_handler.go` (add Payment handler)
- Modify: `internal/handlers/order_handler.go` (preload Payment in Show)
- Modify: `internal/repositories/order_repository.go` (add FindPaymentByOrderID)
- Modify: `internal/services/order_service.go` (add GetPayment)
- Modify: `internal/routes/routes.go` (add payment route on orders)

**Interfaces:**
- Consumes: `OrderRepository.FindByID(id)` (already preloads Customer, Details, Service, Tracking)
- Consumes: `PaymentService.Create(orderID, amount, method, note)`
- Produces: Payment directly from order detail page

- [ ] **Step 1: Add FindPaymentByOrderID to OrderRepository**

In `order_repository.go`:
```go
func (r *OrderRepository) FindPaymentByOrderID(orderID uint) (*models.Payment, error) {
    var payment models.Payment
    err := r.db.Where("order_id = ?", orderID).First(&payment).Error
    if err != nil {
        return nil, err
    }
    return &payment, nil
}
```

- [ ] **Step 2: Add GetPayment to OrderService**

In `order_service.go`:
```go
func (s *OrderService) GetPayment(orderID uint) (*models.Payment, error) {
    return s.orderRepo.FindPaymentByOrderID(orderID)
}
```

- [ ] **Step 3: Add Pay handler to OrderHandler**

Add `Pay` method that calls `paymentService.Create()`. Need to add `paymentService *services.PaymentService` to `OrderHandler`.

```go
func (h *OrderHandler) Pay(c fiber.Ctx) error {
    id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
    userID := c.Locals("user_id").(uint)

    payment, err := h.paymentService.CreateOrUpdate(uint(id), userID, c.FormValue("method"), c.FormValue("note"))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }
    _ = payment
    return c.Redirect().To("/orders/" + c.Params("id"))
}
```

- [ ] **Step 4: Update OrderHandler struct + constructor**

Add `paymentService` and `settingService` fields, update constructor.

- [ ] **Step 5: Add payment card to order show template**

In `templates/orders/show.html`, add after the "Update Status" card (or replace it when status is "sudah_diambil"):

```html
<!-- Payment Card -->
{{if .payment}}
<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-sm">
    <h3 class="text-sm font-semibold text-slate-900 mb-4">Pembayaran</h3>
    <div class="space-y-2 text-sm">
        <div class="flex justify-between">
            <span class="text-slate-500">Status</span>
            <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-emerald-50 text-emerald-700">
                <span class="w-1.5 h-1.5 bg-emerald-500 rounded-full"></span>LUNAS
            </span>
        </div>
        <div class="flex justify-between">
            <span class="text-slate-500">Jumlah</span>
            <span class="font-semibold">Rp {{printf "%.0f" .payment.Amount}}</span>
        </div>
        <div class="flex justify-between">
            <span class="text-slate-500">Metode</span>
            <span>{{if eq .payment.Method "tunai"}}Tunai{{else if eq .payment.Method "qris"}}QRIS{{else if eq .payment.Method "transfer"}}Transfer{{else}}{{.payment.Method}}{{end}}</span>
        </div>
    </div>
</div>
{{else if and (ne .order.Status "dibatalkan") (ne .order.Status "sudah_diambil")}}
<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-sm">
    <h3 class="text-sm font-semibold text-slate-900 mb-4">Pembayaran</h3>
    <form action="/orders/{{.order.ID}}/pay" method="POST" class="space-y-3">
        <div>
            <p class="text-sm text-slate-500 mb-2">Total: <span class="font-bold text-indigo-600">Rp {{printf "%.0f" .order.Total}}</span></p>
            <label class="block text-xs font-medium text-slate-500 mb-1">Metode Pembayaran</label>
            <select name="method" required class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:ring-2 focus:ring-indigo-500 outline-none">
                <option value="tunai">Tunai</option>
                <option value="qris">QRIS</option>
                <option value="transfer">Transfer</option>
            </select>
        </div>
        <input type="text" name="note" placeholder="Catatan (opsional)" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:ring-2 focus:ring-indigo-500 outline-none">
        <button type="submit" class="w-full bg-emerald-600 hover:bg-emerald-700 text-white py-2 rounded-lg text-sm font-medium transition">
            <span class="flex items-center justify-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                Konfirmasi Pembayaran
            </span>
        </button>
    </form>
</div>
{{end}}
```

- [ ] **Step 6: Pass payment data to order show template**

Update `Show` handler in `order_handler.go` to fetch payment:
```go
payment, _ := h.orderService.GetPayment(uint(id))
// add "payment": payment to fiber.Map
```

- [ ] **Step 7: Wire route in routes.go**

Add: `orders.Post("/:id/pay", orderHandler.Pay)`

- [ ] **Step 8: Build + Test**

---

### Task 3: Dashboard Charts Fix (Real Data)

**Files:**
- Modify: `internal/repositories/dashboard_repository.go` (add GetDailyRevenue, GetWeeklyOrderCounts)
- Modify: `internal/services/dashboard_service.go` (add chart data)
- Modify: `internal/handlers/dashboard_handler.go` (pass chart data)
- Modify: `templates/dashboard/index.html` (update chart JS with real data)

**Interfaces:**
- Consumes: existing `DashboardService.GetStats()`
- Produces: Daily revenue array (7 days) + Weekly order counts (7 days) for Chart.js

- [ ] **Step 1: Add chart queries to DashboardRepository**

```go
type DailyRevenue struct {
    Date  string
    Total float64
}

type DailyOrderCount struct {
    Date  string
    Count int64
}

func (r *DashboardRepository) GetDailyRevenue(days int) ([]DailyRevenue, error) {
    var results []DailyRevenue
    startDate := time.Now().AddDate(0, 0, -days+1)
    err := r.db.Model(&models.Payment{}).
        Select("DATE(created_at) as date, COALESCE(SUM(amount), 0) as total").
        Where("created_at >= ? AND status = 'lunas'", startDate).
        Group("DATE(created_at)").
        Order("DATE(created_at) ASC").
        Scan(&results).Error
    return results, err
}

func (r *DashboardRepository) GetDailyOrderCounts(days int) ([]DailyOrderCount, error) {
    var results []DailyOrderCount
    startDate := time.Now().AddDate(0, 0, -days+1)
    err := r.db.Model(&models.Order{}).
        Select("DATE(created_at) as date, COUNT(*) as count").
        Where("created_at >= ?", startDate).
        Group("DATE(created_at)").
        Order("DATE(created_at) ASC").
        Scan(&results).Error
    return results, err
}
```

- [ ] **Step 2: Add chart data to DashboardService**

Add a new response struct `DashboardData` with `RevenueChart` and `OrderChart` arrays. Add method `GetChartData()`. Use a helper to fill 7-day arrays (fill missing days with 0).

- [ ] **Step 3: Update DashboardHandler**

Pass `revenueData` and `orderData` as JSON arrays to the template.

- [ ] **Step 4: Update chart JS in dashboard template**

Replace hardcoded `[0,0,0,0,0,0,0]` with `{{.revenueData}}` and `{{.orderData}}`.

- [ ] **Step 5: Build + Test**

---

### Task 4: Search on All Pages

**Files:**
- Modify: `internal/repositories/customer_repository.go` (add search)
- Modify: `internal/repositories/service_repository.go` (add search)
- Modify: `internal/repositories/payment_repository.go` (add search)
- Modify: `internal/repositories/expense_repository.go` (add search)
- Modify: `internal/repositories/inventory_repository.go` (add search)
- Modify: `internal/handlers/customer_handler.go` (add search param)
- Modify: `internal/handlers/service_handler.go` (add search param)
- Modify: `internal/handlers/payment_handler.go` (add search param)
- Modify: `internal/handlers/expense_handler.go` (add search param)
- Modify: `internal/handlers/inventory_handler.go` (add search param)
- Modify: `templates/customers/index.html` (add search bar)
- Modify: `templates/services/index.html` (add search bar)
- Modify: `templates/payments/index.html` (add search bar)
- Modify: `templates/expenses/index.html` (add search bar)
- Modify: `templates/inventory/index.html` (add search bar)

**Interfaces:**
- Consumes: existing `FindAll(page, limit)` pattern
- Produces: `FindAll(page, limit, search)` pattern

- [ ] **Step 1: Update CustomerRepository.FindAll** to accept `search string` and filter by `name ILIKE %search%` or `phone ILIKE %search%`

- [ ] **Step 2: Update CustomerService.GetAll** to accept and pass `search string`

- [ ] **Step 3: Update CustomerHandler.Index** to read `c.Query("search", "")` and pass to service + template

- [ ] **Step 4: Add search bar to templates/customers/index.html** (pattern from orders/index.html — search input + filter button + reset link)

- [ ] **Step 5-8: Repeat for services** (search by `name`)

- [ ] **Step 9-12: Repeat for payments** (search by `order_number` via join)

- [ ] **Step 13-16: Repeat for expenses** (search by `description` or category name)

- [ ] **Step 17-20: Repeat for inventory** (search by `name` or `category`)

- [ ] **Step 21: Build + Test**

---

### Task 5: Stock Movement Log

**Files:**
- Create: `internal/models/stock_movement.go`
- Create: `internal/repositories/stock_movement_repository.go`
- Create: `internal/services/stock_movement_service.go`
- Modify: `internal/config/database.go` (add StockMovement to AutoMigrate)
- Modify: `internal/handlers/inventory_handler.go` (record movements on Create/Update)
- Modify: `internal/services/inventory_service.go` (call stock movement recording)
- Modify: `templates/inventory/index.html` (add stock log button/column)
- Create: `templates/inventory/movements.html`

**Interfaces:**
- Consumes: `InventoryService.Create`, `InventoryService.Update` (record changes)
- Produces: Stock movement history per item

- [ ] **Step 1: Create StockMovement model**

```go
type StockMovement struct {
    ID          uint      `gorm:"primaryKey"`
    InventoryID uint      `gorm:"not null;index"`
    Inventory   *Inventory `gorm:"foreignKey:InventoryID"`
    Type        string    `gorm:"size:20;not null"` // "in" or "out"
    Quantity    int       `gorm:"not null"`
    PreviousStock int     `gorm:"not null"`
    NewStock    int       `gorm:"not null"`
    Note        string    `gorm:"type:text"`
    CreatedBy   uint      `gorm:"not null"`
    User        *User     `gorm:"foreignKey:CreatedBy"`
    CreatedAt   time.Time
}
```

- [ ] **Step 2: Create StockMovementRepository** — `FindByInventoryID`, `Create`

- [ ] **Step 3: Update InventoryService** — after Create and Update, record stock movement

- [ ] **Step 4: Create movements template** — show log of stock changes per item

- [ ] **Step 5: Add "Riwayat Stok" button in inventory list + movement page**

- [ ] **Step 6: Register StockMovement in AutoMigrate**

- [ ] **Step 7: Build + Test**

---

## Execution Order

Tasks 1-5 are **independent** and can be executed in parallel by separate subagents:
- Subagent A → Task 1 (Nota Printing)
- Subagent B → Task 2 (Payment on Order Page)
- Subagent C → Task 3 (Dashboard Charts)
- Subagent D → Task 4 (Search All Pages)
- Subagent E → Task 5 (Stock Movement Log)

The only shared modification is in `routes.go` (multiple agents add different routes) and `order_handler.go` (print + pay both touch it) — these should be coordinated by the controller (me) after subagents complete their tasks.
