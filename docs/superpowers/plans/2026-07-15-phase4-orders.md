# Phase 4: Orders (Pesanan Laundry) — Implementation Plan

> **For agentic workers:** Execute inline. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Full order management with auto-generated order number, customer & service selection, automatic total calculation, and status tracking history.

**Architecture:** Order → OrderDetail (many-to-many with services) → OrderTracking (status history). Order number auto-generated with format `INV/YYYYMMDD/XXXX`.

**UI Design System (konsisten):**
- Cards, tables, forms, badges same as Customers/Services
- Status badges: menunggu=slate, dicuci=blue, dikeringkan=amber, disetrika=purple, siap_diambil=indigo, sudah_diambil=emerald, dibatalkan=red
- Dropdown select for customer and service

---

### Task 1: Order Models + Migration Updates

**Files:**
- Modify `internal/models/order.go` (add EstimatedDoneDate, customer/service relations)
- Create `internal/models/order_detail.go` (junction table for order → services)
- Create `internal/models/order_tracking.go` (status history)

- [ ] **Step 1: Update Order model** with proper fields and relations

```go
type Order struct {
    ID               uint
    OrderNumber      string    `gorm:"size:50;uniqueIndex;not null"`
    CustomerID       uint      `gorm:"not null;index"`
    Customer         *Customer `gorm:"foreignKey:CustomerID"`
    UserID           uint      `gorm:"not null;index"`
    User             *User     `gorm:"foreignKey:UserID"`
    WeightKg         float64   `gorm:"type:decimal(10,2)"`
    Discount         float64   `gorm:"type:decimal(12,2);default:0"`
    ExtraCost        float64   `gorm:"type:decimal(12,2);default:0"`
    Total            float64   `gorm:"type:decimal(12,2)"`
    EntryDate        time.Time
    EstimatedDoneDate time.Time
    Status           string    `gorm:"size:20;default:menunggu;index"`
    Notes            string    `gorm:"type:text"`
    // Relations
    Details          []OrderDetail    `gorm:"foreignKey:OrderID"`
    Tracking         []OrderTracking  `gorm:"foreignKey:OrderID"`
    // Timestamps
    CreatedAt        time.Time
    UpdatedAt        time.Time
    DeletedAt        gorm.DeletedAt `gorm:"index"`
}
```

- [ ] **Step 2: Create OrderDetail model**

```go
type OrderDetail struct {
    ID        uint
    OrderID   uint     `gorm:"not null;index"`
    ServiceID uint     `gorm:"not null"`
    Service   *Service `gorm:"foreignKey:ServiceID"`
    PricePerKg float64 `gorm:"type:decimal(12,2)"`
}
```

- [ ] **Step 3: Create OrderTracking model**

```go
type OrderTracking struct {
    ID        uint
    OrderID   uint   `gorm:"not null;index"`
    Status    string `gorm:"size:20;not null"`
    Note      string `gorm:"type:text"`
    CreatedBy uint
    CreatedAt time.Time
}
```

- [ ] **Step 4: Add to AutoMigrate in database.go**

- [ ] **Step 5: Verify compile + commit**

---

### Task 2: Order Repository

**Files:**
- Create `internal/repositories/order_repository.go`

Features:
- GenerateOrderNumber — format `INV/YYYYMMDD/XXXX` (auto-increment per day)
- FindAll — paginated, filterable by status/date/customer
- FindByID — with Preload for Customer, Details, Tracking
- Create — with transaction (order + details + initial tracking)
- UpdateStatus — add tracking entry
- Delete — soft delete

- [ ] **Step 1: Write repository** with all methods

- [ ] **Step 2: Verify compile + commit**

---

### Task 3: Order Service

**Files:**
- Create `internal/services/order_service.go`

Business logic:
- CreateOrder: validate customer, calculate total from service prices + weight, generate order number, initial status "menunggu"
- UpdateStatus: validate valid transitions, add tracking entry
- CalculateTotal: sum of (service_price * weight) - discount + extra_cost

Status transitions:
```
menunggu → dicuci → dikeringkan → disetrika → siap_diambil → sudah_diambil
                                                              ↳ dibatalkan (any → dibatalkan)
```

- [ ] **Step 1: Write service** with all business rules

- [ ] **Step 2: Verify compile + commit**

---

### Task 4: Order Handler + Routes

**Files:**
- Create `internal/handlers/order_handler.go`
- Modify `internal/routes/routes.go`

Handler methods:
- Index — list orders with filters (status, date range, search customer)
- New — form with customer + service dropdowns
- Create — process order creation
- Show — detail with tracking history
- Edit — form pre-filled
- Update — process update
- UpdateStatus — POST to change status + add tracking
- Delete

Routes:
```
GET    /orders            → OrderHandler.Index
GET    /orders/new        → OrderHandler.New
POST   /orders            → OrderHandler.Create
GET    /orders/:id        → OrderHandler.Show
GET    /orders/:id/edit   → OrderHandler.Edit
POST   /orders/:id        → OrderHandler.Update
POST   /orders/:id/status → OrderHandler.UpdateStatus
POST   /orders/:id/delete → OrderHandler.Delete
```

- [ ] **Step 1: Write handler**

- [ ] **Step 2: Wire routes**

- [ ] **Step 3: Verify compile + commit**

---

### Task 5: Order Templates

**Files:**
- Create `templates/orders/index.html`
- Create `templates/orders/form.html`
- Create `templates/orders/show.html`

**index.html:** Table with columns: Order Number, Customer, Service, Weight, Total, Status, Date, Actions
- Filters: status dropdown, date range, customer search
- Status badges with colors

**form.html:** Create/Edit form
- Customer select2-like dropdown (searchable)
- Multiple service selection with weight input per service
- Discount and extra cost fields
- Auto total calculation preview (JS)

**show.html:** Detail view
- Order info card
- Customer info card
- Service list with prices
- Payment status
- Status tracking timeline
- Action buttons (Update Status, Payment, Delete)

- [ ] **Step 1: Create index template**

- [ ] **Step 2: Create form template** with Alpine.js for dynamic calculations

- [ ] **Step 3: Create show template** with tracking timeline

- [ ] **Step 4: Verify compile**

```bash
go build ./... && echo "BUILD OK"
```

- [ ] **Step 5: Test with server**

```bash
# Start server and test order CRUD
go run ./cmd/server/
```

- [ ] **Step 6: Commit**

---

### Task 6: Final Test + Push

- [ ] **Step 1: Full verification**

```bash
go build ./... && go vet ./...
```

- [ ] **Step 2: Push to GitHub**

```bash
git push origin master
```
