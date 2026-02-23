# FeedbackBot Regression Test Plan

## 1. Auth Flows

### 1.1 Register
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 1 | Register with valid name, email, password | POST /auth/register | 200, returns access_token + refresh_token + user info |
| 2 | Register with duplicate email | POST /auth/register | 409, "email already registered" |
| 3 | Register with missing fields | POST /auth/register | 400 |
| 4 | Auto-creates tenant from email domain | POST /auth/register | UserTenant record created with correct tenant |

### 1.2 Login
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 5 | Login with valid credentials | POST /auth/login | 200, returns tokens + user info |
| 6 | Login with wrong password | POST /auth/login | 401, "invalid email or password" |
| 7 | Login with non-existent email | POST /auth/login | 401, "invalid email or password" |
| 8 | Login with missing fields | POST /auth/login | 400 |

### 1.3 Refresh Token
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 9 | Refresh with valid refresh_token | POST /auth/refresh | 200, new token pair |
| 10 | Refresh with expired token | POST /auth/refresh | 401 |
| 11 | Refresh with access token (wrong type) | POST /auth/refresh | 401, "invalid token type" |
| 12 | Refresh with invalid string | POST /auth/refresh | 401 |

### 1.4 Me
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 13 | Get current user with valid token | GET /auth/me | 200, user_id, email, name, role, tenant_id |
| 14 | Get current user without token | GET /auth/me | 401 |
| 15 | Get current user with expired token | GET /auth/me | 401 |

## 2. Tenant CRUD

### 2.1 Create Tenant
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 16 | Create tenant with valid data + auth | POST /tenants | 201, tenant object |
| 17 | Create tenant without auth | POST /tenants | 401 |
| 18 | Create tenant with duplicate slug | POST /tenants | 500 (unique constraint) |
| 19 | Create tenant with missing fields | POST /tenants | 400 |
| 20 | UserTenant association created | POST /tenants | UserTenant record exists |

### 2.2 Get Tenant
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 21 | Get existing tenant by ID | GET /tenants/:id | 200, tenant object |
| 22 | Get non-existent tenant | GET /tenants/:id | 404 |

## 3. Bot CRUD + Token Verification

### 3.1 List Bots
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 23 | List bots for tenant | GET /bots | 200, array of bots |
| 24 | List bots without auth | GET /bots | 401 |
| 25 | Bots scoped to tenant | GET /bots | Only tenant's bots returned |

### 3.2 Create Bot
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 26 | Create bot with valid Telegram token | POST /bots | 201, bot with username + verified=true |
| 27 | Create bot with invalid token | POST /bots | 400, "Invalid bot token" |
| 28 | Create bot without auth | POST /bots | 401 |
| 29 | Create bot with missing token field | POST /bots | 400 |

### 3.3 Get Bot
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 30 | Get existing bot | GET /bots/:id | 200, bot object (no token exposed) |
| 31 | Get bot from different tenant | GET /bots/:id | 404 |

### 3.4 Delete Bot
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 32 | Delete own bot | DELETE /bots/:id | 200 |
| 33 | Delete bot from different tenant | DELETE /bots/:id | 404 |

## 4. Group Management

### 4.1 List Groups
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 34 | List groups for tenant | GET /groups | 200, array |
| 35 | Groups scoped to tenant | GET /groups | Only tenant's groups |

### 4.2 Get Group
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 36 | Get existing group | GET /groups/:id | 200, group object |
| 37 | Get group from different tenant | GET /groups/:id | 404 |

### 4.3 Update Group
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 38 | Toggle is_active | PATCH /groups/:id | 200, updated group |
| 39 | Update group from different tenant | PATCH /groups/:id | 404 |

### 4.4 Update Group Config
| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 40 | Set post_to_group=true | PATCH /groups/:id/config | 200, updated config |
| 41 | Set forum_topic_id | PATCH /groups/:id/config | 200, updated config |
| 42 | Update config for non-existent group | PATCH /groups/:id/config | 404 |

## 5. Feedback Submission (Telegram Bot DM)

### 5.1 Private Messages
| # | Test Case | Expected |
|---|-----------|----------|
| 43 | /start command | Welcome message returned |
| 44 | Plain text message (1 group) | Feedback created, confirmation sent |
| 45 | /adminOnly + message | Feedback created with admin_only=true |
| 46 | /adminOnly without text | Error: "Please write your feedback after /adminOnly" |
| 47 | Empty message | Error: "Please send a text message" |
| 48 | Message > 4000 chars | Error: "Your message is too long" |
| 49 | Message with multiple active groups | Inline keyboard shown for group selection |
| 50 | No active groups | Error: "No active groups found" |

### 5.2 Callback Queries (Group Selection)
| # | Test Case | Expected |
|---|-----------|----------|
| 51 | Select group from keyboard | Feedback submitted to selected group |
| 52 | Pending feedback expired (not found in DB) | "Session expired" message |
| 53 | Invalid group ID in callback | Error handled gracefully |

### 5.3 Feedback Posting
| # | Test Case | Expected |
|---|-----------|----------|
| 54 | post_to_group=true, not admin_only | Message posted to group chat |
| 55 | post_to_group=true, admin_only | Message NOT posted to group |
| 56 | post_to_group=false | Message NOT posted to group |
| 57 | Forum group with topic_id set | Message posted to specific topic |

## 6. Feedback Listing with Filters

| # | Test Case | Method | Expected |
|---|-----------|--------|----------|
| 58 | List feedbacks for group | GET /feedbacks?group_id=X | 200, paginated list |
| 59 | Filter admin_only=true | GET /feedbacks?admin_only=true | Only admin_only feedbacks |
| 60 | Filter admin_only=false | GET /feedbacks?admin_only=false | Only public feedbacks |
| 61 | Date range filter | GET /feedbacks?date_from=X&date_to=Y | Filtered by date |
| 62 | Pagination (page + limit) | GET /feedbacks?page=2&limit=10 | Correct offset |
| 63 | SenderID never exposed in response | GET /feedbacks | sender_id not in JSON |
| 64 | Feedbacks scoped to tenant | GET /feedbacks | Only tenant's feedbacks |

## 7. Bot Group Events (my_chat_member)

| # | Test Case | Expected |
|---|-----------|----------|
| 65 | Bot added to group as member | Group created in DB, FeedbackConfig created |
| 66 | Bot added to group as admin | Group created in DB |
| 67 | Bot re-added to existing group | Group reactivated (is_active=true) |
| 68 | Bot removed from group (left) | Group deactivated (is_active=false) |
| 69 | Bot removed from group (kicked) | Group deactivated |
| 70 | Event from private chat | Ignored |

## 8. Multi-Tenant Isolation

| # | Test Case | Expected |
|---|-----------|----------|
| 71 | User A cannot see User B's bots | GET /bots returns only own tenant's bots |
| 72 | User A cannot see User B's groups | GET /groups returns only own tenant's groups |
| 73 | User A cannot see User B's feedbacks | GET /feedbacks returns only own tenant's feedbacks |
| 74 | User A cannot delete User B's bot | DELETE /bots/:id returns 404 |
| 75 | User A cannot update User B's group | PATCH /groups/:id returns 404 |
| 76 | TenantMiddleware rejects missing tenant_id | 400 response |

## 9. JWT Token Lifecycle

| # | Test Case | Expected |
|---|-----------|----------|
| 77 | Access token expires after 24h | 401 after expiry |
| 78 | Refresh token expires after 7d | 401 after expiry |
| 79 | Refresh yields new token pair | Both tokens are new |
| 80 | Invalid JWT signature | 401 |
| 81 | Malformed JWT string | 401 |
| 82 | Missing Authorization header | 401 |
| 83 | Non-Bearer scheme | 401 |

## 10. Graceful Shutdown

| # | Test Case | Expected |
|---|-----------|----------|
| 84 | SIGTERM sent to process | Bot polling stops, logs shutdown message |
| 85 | SIGINT sent to process | Bot polling stops gracefully |

## 11. Pending Feedback Persistence

| # | Test Case | Expected |
|---|-----------|----------|
| 86 | Pending feedback survives restart | DB-backed, not lost |
| 87 | New pending feedback overwrites old for same user | Only latest stored |
| 88 | Pending feedback deleted after use | Cleaned up from DB |
