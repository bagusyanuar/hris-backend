# User Stories & System Flows: Employee Module

## 1. User Stories

**US-01: Create New Employee**
- **As an** HR Admin
- **I want to** add a new employee to the system with their personal and bank details
- **So that** they are officially registered and can receive payroll.
- **Acceptance Criteria:**
  - System must validate that KTP is 16 digits.
  - System must ensure at least 1 primary bank account is provided.
  - System must rollback entirely if any sub-data (bank, personal) fails to save.

**US-02: Offboard Employee**
- **As an** HR Admin
- **I want to** offboard a resigning employee
- **So that** they no longer have access to the system but their historical data remains intact.
- **Acceptance Criteria:**
  - Status becomes `INACTIVE`.
  - `resign_date` is populated.
  - No records are physically deleted from the database (Soft Delete).
  - Triggers an event to the Auth Module to revoke login access.

---

## 2. Sequence Diagrams

### 2.1. Create Employee Flow (Happy Path)

This diagram illustrates the layered architecture communication when creating a new employee.

```mermaid
sequenceDiagram
    actor Admin as HR Admin
    participant FE as Frontend UI
    participant Handler as HTTP Handler
    participant Svc as Employee Service
    participant Auth as Auth Module (Internal API)
    participant DB as Postgres Database

    Admin->>FE: Submits Employee Wizard Form
    FE->>Handler: POST /api/v1/employees
    Handler->>Svc: Pass DTO (Mapped)
    
    rect rgb(240, 248, 255)
        Note over Svc: Domain Validation
        Svc->>Svc: Validate KTP format
        Svc->>Svc: Ensure >= 1 Primary Bank
    end
    
    Svc->>Auth: Request create user account
    Auth-->>Svc: Return user_id
    
    rect rgb(255, 240, 245)
        Note over Svc, DB: Database Transaction
        Svc->>DB: Begin TX
        DB->>DB: Insert employees
        DB->>DB: Insert employee_personal_data
        DB->>DB: Insert employee_banks
        Svc->>DB: Commit TX
    end
    
    Svc-->>Handler: Return Entity
    Handler-->>FE: 201 Created Response
    FE-->>Admin: Show Success Notification
```
