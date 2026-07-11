# User Stories & System Flows: Employee Module

## 1. User Stories

**US-01: Create New Employee (Progressive)**
- **As an** HR Admin
- **I want to** add a new employee to the system step-by-step (Progressive Save)
- **So that** I don't lose data if I get interrupted, and I can save drafts easily.
- **Acceptance Criteria:**
  - Step 1 (Core) must return an `employee_id` to be used for the next steps.
  - System must validate that KTP is 16 digits at the Personal Data step.
  - System must allow saving as draft, even if bank details or contacts are not yet provided.

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
    participant FE as Frontend UI (Stepper)
    participant Handler as HTTP Handler
    participant Svc as Employee Service
    participant DB as Postgres Database

    Note over Admin, DB: STEP 1: Core Data
    Admin->>FE: Submits Step 1
    FE->>Handler: POST /api/v1/employees
    Handler->>Svc: Pass Core DTO
    Svc->>DB: Insert into employees
    DB-->>Svc: Return employee_id
    Svc-->>Handler: Return Entity
    Handler-->>FE: 201 Created (with ID)

    Note over Admin, DB: STEP 2: Personal Data
    Admin->>FE: Submits Step 2
    FE->>Handler: PUT /api/v1/employees/{id}/personal-data
    Handler->>Svc: Pass Personal DTO
    Svc->>Svc: Validate KTP format
    Svc->>DB: Insert/Update employee_personal_data
    Svc-->>FE: 200 OK Response

    Note over Admin, DB: STEP 3: Contacts & Banks
    Admin->>FE: Submits Step 3 & 4
    FE->>Handler: PUT /.../contacts & POST /.../banks
    Handler->>DB: Insert Data
    Handler-->>FE: 200 OK Response
    FE-->>Admin: Show Final Success Notification
```
