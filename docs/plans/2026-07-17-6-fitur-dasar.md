# 6 Fitur Dasar Laundry Management — Implementation Plan

> **For agentic workers:** Use subagent-driven-development to implement task-by-task.

**Goal:** Implement 6 fitur dasar yang masih kurang: Quick Actions Status, Quick Stats Bar, Activity Log Dashboard, Overdue Alerts, Modal Delete, Notifikasi.

**Architecture:** Monolith Go Fiber v3 + GORM + Alpine.js. Tiap fitur hanya ubah handler + template + repository (backend sudah siap).

**Tech Stack:** Go 1.23, Fiber v3, GORM, PostgreSQL, Alpine.js, Chart.js

---

## Task 1: Quick Actions + Quick Stats di Order Index (subagent 1)

**Files:**
- Modify: `internal/handlers/order_handler.go` — add `BulkUpdateStatus` handler (POST with JSON), add `StatusCount` to Index render
- Modify: `internal/services/order_service.go` — `QuickUpdateStatus(id, newStatus, userID)` without note requirement
- Modify: `templates/orders/index.html` — add quick stats bar di atas tabel, add dropdown quick action di kolom aksi tiap row
- Modify: `internal/routes/routes.go` — add `POST /orders/:id/quick-status`

## Task 2: Activity Log + Overdue Alerts di Dashboard (subagent 2)

**Files:**
- Modify: `internal/repositories/dashboard_repository.go` — add `GetRecentTracking(limit int)` and `CountOverdueOrders()`
- Modify: `internal/services/dashboard_service.go` — expose recent tracking data, overdue count
- Modify: `internal/handlers/dashboard_handler.go` — pass activity log and overdue data to template
- Modify: `templates/dashboard/index.html` — replace placeholder "Belum ada aktivitas" with real tracking rows, add overdue alert card

## Task 3: Modal Delete + Notifikasi (subagent 3)

**Files:**
- Create: `templates/partials/delete_modal.html` — Alpine.js modal konfirmasi
- Modify: `templates/orders/index.html`, `templates/customers/index.html`, `templates/services/index.html`, `templates/payments/index.html`, `templates/expenses/index.html`, `templates/inventory/index.html`, `templates/users/index.html` — replace confirm() with modal
- Create: `templates/partials/toast.html` — notifikasi sukses/error

---

## Execution Strategy

3 subagents in parallel (independent file sets):
- SG1: Order index enhancements
- SG2: Dashboard enhancements
- SG3: Modal + Notifikasi

Then: build, verify, commit.
