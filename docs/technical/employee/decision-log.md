# Architecture Decision Records (ADR): Employee Module

## ADR-001: Usage of Soft Deletes over Hard Deletes
- **Date:** 2026-07-10
- **Status:** Accepted
- **Context:** Employees resign or are terminated. Deleting their records physically (`DELETE FROM`) causes referential integrity issues in historical Payroll and Attendance data. If an employee is deleted, their past payslips would break.
- **Decision:** Implement `deleted_at` timestamp for soft deleting employees. The application layer will automatically filter out soft-deleted records for active employee queries.

## ADR-002: Separating Personal Data from Core Employee Table
- **Date:** 2026-07-10
- **Status:** Accepted
- **Context:** The `employees` table can become a "God Table" with 50+ columns if we combine employment data (job position, status) with personal data (KTP, blood type, marriage status).
- **Decision:** Split into `employees` (core employment) and `employee_personal_data` (demographics) with a 1-to-1 relationship to keep table widths manageable and improve query performance for purely organizational queries.

## ADR-003: Progressive Save (Multiple Endpoints) for Creation
- **Date:** 2026-07-11
- **Status:** Accepted (Revised from Strict Single Transaction)
- **Context:** Creating an employee involves inserting into multiple tables (`employees`, `employee_personal_data`, `employee_banks`). Initially, this was planned as a single massive transaction. However, this risks data loss on the frontend if a validation error occurs at the very end of a long wizard form.
- **Decision:** Implement a **Progressive Save** approach. The Frontend will call separate endpoints per step (`POST /employees` -> `PUT /employees/{id}/personal-data`, etc.). This allows 'Save as Draft', simplifies file uploads, and isolates validation errors per step.

## ADR-004: UUID v4 as Primary Keys
- **Date:** 2026-07-10
- **Status:** Accepted
- **Context:** Using auto-increment integer IDs exposes business metrics (e.g., how many employees we have) and creates issues if we ever migrate or merge databases.
- **Decision:** Use UUID v4 for all primary keys, generated at the Domain Entity layer prior to hitting the database.
