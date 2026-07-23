---
name: scaffold-rfc
description: Guide for scaffolding Technical Specifications / Request For Comments (RFC) documents based on PRDs.
---

# Scaffolding RFC / Technical Specification

When a user asks to create an RFC, Technical Spec, or Implementation Plan based on a PRD, you MUST follow this structure. The RFC acts as the absolute engineering blueprint for human developers and AI agents before any code is written.

## 1. Directory & Naming Convention
- **Location:** All Technical documents MUST be saved in a specific domain sub-directory under `docs/technical/`. For example, `docs/technical/employee/`.
- **Files Generated:** You MUST split the documentation into the following files to maintain enterprise standards:
  - **In `docs/technical/<domain>/`:**
    1. `tech-spec.md` (The Core RFC containing Architecture, API Contracts, Database Schema, and Mermaid ERD).
    2. `user-stories.md` (Detailed User Stories, Acceptance Criteria, Edge Cases, and Mermaid Sequence Diagrams for flows).
    3. `decision-log.md` (ADR - Architecture Decision Records to log *why* specific technical decisions were made).
    4. `data-dictionary.md` (Detailed definitions of ENUMs, statuses, and complex database field semantics).
    5. `infrastructure.md` (Specific infra needs like S3 buckets, Kafka topics, or external services).
    6. `test-plan.md` (QA scenarios, boundary tests, and integration test plans).
  - **In `docs/databases/`:**
    7. `<domain>.dbml` (The detailed database schema definitions in DBML format for ERD generation).

## 2. Mandatory Tech Spec Structure (`tech-spec.md`)

Every Technical Spec/RFC MUST contain the following core sections:

### 2.1. Overview & PRD Reference
Briefly state what is being engineered and provide a direct link to the relevant PRD in `docs/PRD/`. This ensures traceability.

### 2.2. System Architecture & Boundaries (DDD)
Explain how this module fits into the Domain-Driven Design (DDD) architecture of the HRIS backend.
- Define the **Aggregate Root**.
- Define the **Value Objects** and child entities.
- Outline the **Folder Structure** that will be generated following the **domain-first** layout: one context folder `internal/[module]/` containing `domain/`, `application/`, `adapter/` (+ `adapter/models/`), and `transport/http/`.

### 2.3. Cross-Domain Dependencies (Bounded Context Integrations)
Examine how this module communicates with other modules.
- **Upstream/Downstream:** Which modules depend on this one, and which does this module depend on?
- **Communication Method:** State explicitly how they communicate (e.g., Direct Application Service injection, Go Channels/Event Bus, or Pub/Sub messaging).
- **Data Consistency:** Does it require Saga pattern / Eventual Consistency, or synchronous transaction?

### 2.4. Detailed Database Schema & Migrations
Translate the PRD's logical data schema into strict physical database schemas.
- Provide Markdown tables or DDL containing: **Field Name, Data Type** (e.g., UUID, VARCHAR(150), DECIMAL(10,2)), **Constraints** (PK, FK, UNIQUE, NOT NULL), and **Indexes** for performance.
- Detail the soft-delete (`deleted_at`) and auditing columns (`created_at`, `updated_at`).

### 2.4. API Contracts (Endpoints)
Define the HTTP/REST interfaces that the Frontend will consume.
For each endpoint, detail:
- **Method & Path:** e.g., `POST /api/v1/employees`
- **Request Payload:** JSON structure with strict validation rules (e.g., required, min length, email format).
- **Response Payload:** Success JSON (200/201) and standardized Error JSON (400/422/404/500).

### 2.5. Implementation Details & Algorithms
Explain the internal technical flow and business logic execution.
- **Sequence / Flow:** How the layers interact (HTTP Handler -> Application Service -> Domain -> Adapter Repo). Use Mermaid.js diagrams if the flow is highly complex.
- **Database Transactions:** Identify operations that require strict ACID Transactions (e.g., saving an Employee, Personal Data, and Bank Account atomically in one commit).
- **Domain Errors:** List the custom Go errors to be defined (e.g., `ErrEmployeeNotFound`, `ErrKtpDuplicate`).

### 2.6. Security, Performance & Technical Constraints
- **Security (Auth/Authz):** Which JWT roles/permissions are required to hit the endpoints? 
- **Performance:** Strategies to avoid N+1 queries when fetching relational data, and mandatory Pagination for list endpoints.
- **Data Masking/Sanitization:** Technical implementation of how sensitive data is masked before being serialized to JSON.
