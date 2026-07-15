# Phase 3: Master Data (Customers & Services) — Implementation Plan

> **For agentic workers:** Execute inline. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Full CRUD for Customers and Services with search, pagination, and consistent UI following the existing design system.

**Architecture:** Handler → Service → Repository pattern. Templates extend main layout with sidebar/navbar. HTMX for search and pagination interactions.

**UI Design System (konsisten):**
- Cards: white bg, rounded-xl, shadow-sm, border-slate-200
- Tables: thead bg-slate-50, text-sm, sticky header
- Buttons: btn-primary (indigo-600), btn-danger (red-600), btn-ghost (slate)
- Badges: IN=emerald, OUT=red, active=emerald, inactive=slate
- Forms: label above input, rounded-lg, border-slate-300, focus:ring-indigo-500
- Pagination: numbered pages with prev/next

---

### Task 1: Customer Repository + Service + Handler

**Files:**
- Create `internal/repositories/customer_repository.go`
- Create `internal/services/customer_service.go`
- Create `internal/handlers/customer_handler.go`
- Modify `internal/routes/routes.go` (wire routes)

- [ ] **Step 1: Create CustomerRepository**

```go
// Package with FindAll (paginated+search), FindByID, Create, Update, Delete
```

- [ ] **Step 2: Create CustomerService** (pass-through + validation)

- [ ] **Step 3: Create CustomerHandler** (Index, New, Create, Show, Edit, Update, Delete)

- [ ] **Step 4: Wire routes**
  - GET /customers → CustomerHandler.Index
  - GET /customers/new → CustomerHandler.New
  - POST /customers → CustomerHandler.Create
  - GET /customers/:id → CustomerHandler.Show
  - GET /customers/:id/edit → CustomerHandler.Edit
  - POST /customers/:id → CustomerHandler.Update
  - POST /customers/:id/delete → CustomerHandler.Delete

- [ ] **Step 5: Verify compile**

```bash
go build ./... && echo "BUILD OK"
```

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat: add customer CRUD (repository, service, handler)"
```

---

### Task 2: Customer Templates

**Files:**
- Create `templates/customers/index.html`
- Create `templates/customers/form.html`
- Create `templates/customers/show.html`
- Create `templates/partials/pagination.html`

- [ ] **Step 1: Create index template** — table with search bar, action buttons, pagination

- [ ] **Step 2: Create form template** — create/edit form, shared for both

- [ ] **Step 3: Create show template** — detail view

- [ ] **Step 4: Create pagination partial** — reusable component

- [ ] **Step 5: Verify compile + commit**

---

### Task 3: Service Repository + Service + Handler

**Files:**
- Create `internal/repositories/service_repository.go`
- Create `internal/services/service_service.go`
- Create `internal/handlers/service_handler.go`
- Modify `internal/routes/routes.go`

- [ ] **Step 1: Create ServiceRepository** (paginated+search by name)

- [ ] **Step 2: Create ServiceService** (pass-through + validation)

- [ ] **Step 3: Create ServiceHandler** (Index, New, Create, Edit, Update, Delete)

- [ ] **Step 4: Wire routes**

- [ ] **Step 5: Verify compile + commit**

---

### Task 4: Service Templates

**Files:**
- Create `templates/services/index.html`
- Create `templates/services/form.html`

- [ ] **Step 1: Create index template** — table with search, status toggle, actions

- [ ] **Step 2: Create form template** — name, price per kg, estimated hours, description, status

- [ ] **Step 3: Verify compile**

```bash
go build ./... && echo "BUILD OK"
```

- [ ] **Step 4: Test with server**

```bash
# Start server, test customer CRUD flow
go run ./cmd/server/
```

- [ ] **Step 5: Commit**

---

### Task 5: Final Test + Push

- [ ] **Step 1: Full verification**

```bash
go build ./... && go vet ./...
```

- [ ] **Step 2: Push to GitHub**

```bash
git push origin master
```
