# FeedbackBot User Stories

## Overview

User stories for FeedbackBot — a multi-tenant SaaS Telegram bot for anonymous feedback collection. Stories are organized by feature area and prioritized by impact.

---

## US-1: Dashboard with Feedback Stats

**Priority:** High
**As a** tenant admin,
**I want** a dashboard showing feedback statistics at a glance,
**so that** I can quickly understand feedback volume and trends without scrolling through individual messages.

### Acceptance Criteria
- [ ] Dashboard is the default page after login (`/` route)
- [ ] Shows total feedback count for the current tenant
- [ ] Shows feedback count per group (bar or list)
- [ ] Shows feedback received today / this week / this month
- [ ] Shows ratio of admin-only vs public feedback
- [ ] Shows number of active groups and bots
- [ ] Data refreshes when navigating back to dashboard

### API Requirements
- `GET /dashboard/stats` — returns aggregated counts:
  ```json
  {
    "total_feedbacks": 142,
    "feedbacks_today": 5,
    "feedbacks_this_week": 23,
    "feedbacks_this_month": 67,
    "admin_only_count": 18,
    "public_count": 124,
    "active_groups": 3,
    "active_bots": 1,
    "per_group": [
      { "group_id": 1, "group_title": "Support", "count": 89 },
      { "group_id": 2, "group_title": "Bugs", "count": 53 }
    ]
  }
  ```

### Implementation Notes
- Backend: New `svc_dashboard` service with a single stats endpoint
- Frontend: New `DashboardPage` component with Mantine `StatsGroup`, `SimpleGrid`, and a bar chart (Mantine Charts or recharts)
- Scoped to tenant via `TenantMiddleware`

---

## US-2: Bot List Page

**Priority:** High (already partially implemented)
**As a** tenant admin,
**I want** to see all my registered Telegram bots in a list,
**so that** I can manage, verify, or remove them.

### Acceptance Criteria
- [x] Bot list page at `/settings/bot` shows all bots for the tenant
- [x] Each bot card shows username and verification status
- [x] Delete button removes the bot
- [ ] Bot card shows the number of groups the bot is in
- [ ] Bot card shows when the bot was added (created_at)
- [ ] Empty state message when no bots exist
- [x] Loading spinner while fetching

### API Requirements
- `GET /bots` — already implemented, returns array of bots
- Consider enriching response with `group_count` field

### Implementation Notes
- Most of this is already done via `BotConfigPage` + `useGetBots`
- Remaining work: add `group_count` to bot response, show created_at

---

## US-3: Export Feedbacks as CSV

**Priority:** High
**As a** tenant admin,
**I want** to export feedback data as a CSV file,
**so that** I can analyze feedback in spreadsheets or share reports with stakeholders.

### Acceptance Criteria
- [ ] "Export CSV" button on the feedbacks page
- [ ] Exports all feedbacks matching current filters (group, date range, admin_only)
- [ ] CSV columns: ID, Group, Message, AdminOnly, PostedToGroup, CreatedAt
- [ ] SenderID is NEVER included in the export (anonymity preserved)
- [ ] File downloads with name like `feedbacks_2026-02-23.csv`
- [ ] Works for up to 10,000 records without timeout

### API Requirements
- `GET /feedbacks/export?group_id=X&date_from=Y&date_to=Z&admin_only=true`
- Returns `Content-Type: text/csv` with `Content-Disposition: attachment`

### Implementation Notes
- Backend: Add export handler in `svc_feedback` that streams CSV
- Frontend: Trigger download via `window.open()` or Axios blob response
- Reuse existing feedback query filters

---

## US-4: Feedback Search

**Priority:** Medium
**As a** tenant admin,
**I want** to search feedback by keyword,
**so that** I can quickly find specific feedback messages.

### Acceptance Criteria
- [ ] Search input on the feedbacks page
- [ ] Searches feedback message text (case-insensitive)
- [ ] Works in combination with existing filters (group, date, admin_only)
- [ ] Results update as user types (debounced 300ms)
- [ ] Shows "No results" state when search yields nothing
- [ ] Search term is highlighted in results

### API Requirements
- Add `search` query parameter to `GET /feedbacks`:
  `GET /feedbacks?group_id=X&search=billing`
- Backend uses `WHERE message ILIKE '%keyword%'`

### Implementation Notes
- Backend: Add `search` param to existing feedback list handler
- Frontend: Add `TextInput` with search icon, debounce with `useDebouncedValue` from Mantine hooks
- Pass search term to `useGetFeedbacks` query params

---

## US-5: User Management

**Priority:** Medium
**As a** tenant admin,
**I want** to invite and manage other users in my tenant,
**so that** my team members can also view and manage feedback.

### Acceptance Criteria
- [ ] User management page at `/settings/users`
- [ ] List all users in the current tenant
- [ ] Invite new user by email (sends invitation or creates account)
- [ ] Assign roles: admin, viewer
- [ ] Admin can remove users from tenant
- [ ] Viewer role can only read feedbacks and groups (no bot/tenant management)

### API Requirements
- `GET /users` — list tenant users
- `POST /users/invite` — invite user by email
- `PATCH /users/:id/role` — update user role
- `DELETE /users/:id` — remove user from tenant

### Implementation Notes
- Backend: New `svc_users` service
- Role-based access control (RBAC) middleware checks `UserTenant.role`
- Add `role` field to `UserTenant` model (default: "admin" for creator)
- Invitation flow: create user + UserTenant if email not registered; add UserTenant if already registered

---

## US-6: Notification Settings

**Priority:** Low
**As a** tenant admin,
**I want** to configure email or in-app notifications for new feedback,
**so that** I don't have to check the dashboard constantly.

### Acceptance Criteria
- [ ] Notification settings page at `/settings/notifications`
- [ ] Toggle: email notifications on/off
- [ ] Toggle: receive notifications for all feedback or admin-only feedback
- [ ] Configure notification frequency: instant, hourly digest, daily digest
- [ ] Notifications scoped per group

### API Requirements
- `GET /notifications/settings` — get current settings
- `PATCH /notifications/settings` — update settings
- Background job for digest emails

### Implementation Notes
- Backend: New `NotificationSetting` model (user_id, group_id, email_enabled, frequency, filter)
- Email sending via SMTP or third-party service (SendGrid, Resend)
- Digest requires a cron job / scheduled task
- Consider starting with just instant email notifications as MVP

---

## Priority Matrix

| Story | Impact | Effort | Priority |
|-------|--------|--------|----------|
| US-1: Dashboard Stats | High | Medium | **P0** |
| US-3: Export CSV | High | Low | **P0** |
| US-4: Feedback Search | Medium | Low | **P1** |
| US-2: Bot List Enhancements | Medium | Low | **P1** |
| US-5: User Management | Medium | High | **P2** |
| US-6: Notification Settings | Low | High | **P3** |

## Recommended Implementation Order

1. **US-3: Export CSV** — Low effort, high value. Admins need reporting immediately.
2. **US-1: Dashboard Stats** — High impact, gives the app a professional landing page.
3. **US-4: Feedback Search** — Quick win, improves daily usability.
4. **US-2: Bot List Enhancements** — Small improvements to existing page.
5. **US-5: User Management** — Important for team adoption but higher effort.
6. **US-6: Notifications** — Nice-to-have, defer until core features are solid.
