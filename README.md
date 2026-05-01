# Project Challenge: “Team Workspace API”

Build a REST API using **Gin + GORM + JWT** with:

* Authentication
* Two full CRUD resources
* Real associations
* Authorization rules
* Clean architecture

No shortcuts.

---

# Core Domain

## User

* id
* name
* email (unique)
* password (hashed)
* role (optional: `member` / `admin`)

## Workspace

* id
* name
* owner_id (User)
* members (many-to-many with users)

## Task

* id
* title
* description
* status (`todo`, `in_progress`, `done`)
* workspace_id
* assigned_to (User)

---

# Associations (Important)

You must implement:

* A user **owns many workspaces**
* A workspace **has many tasks**
* A workspace **has many members**
* A task **belongs to a workspace**
* A task **is assigned to a user**

That gives you:

* One-to-many
* Many-to-many
* Multiple foreign keys
* Authorization complexity

---

# Authentication Requirements

### Public

* `POST /register`
* `POST /login`

### Protected

Everything else.

---

# Authorization Rules (This Is Where It Gets Interesting)

You must enforce:

1. A user can only see workspaces they:

    * Own OR
    * Are a member of

2. A user can only create tasks inside workspaces they belong to.

3. A user can only update:

    * Tasks assigned to them
    * OR if they are workspace owner

4. Only workspace owner can:

    * Delete workspace
    * Add/remove members

This forces you to:

* Read JWT claims
* Use context values
* Validate ownership properly
* Write non-trivial DB queries

---

# Required CRUD Routes

## Workspaces

* Create workspace
* List my workspaces
* Get workspace by id
* Update workspace
* Delete workspace
* Add member
* Remove member

## Tasks

* Create task
* List workspace tasks
* Get task
* Update task
* Delete task

---

# Technical Requirements

You must implement:

* JWT authentication middleware
* Password hashing
* Route grouping
* Database migrations
* Proper HTTP status codes
* Error handling (no panics)
* JSON validation
* Foreign key constraints

---

# What This Will Force You To Learn

* Middleware chaining
* Claims-based auth
* Association preloading
* Join tables (many-to-many)
* Query optimization
* Struct tags & validation
* Clean folder structure
* Separation of concerns

---

# Bonus Challenges

* Pagination on task listing
* Filtering by status
* Soft deletes
* Refresh tokens
* Dockerize the project
* Write unit tests for handlers
* Implement repository pattern
* Add rate limiting middleware
