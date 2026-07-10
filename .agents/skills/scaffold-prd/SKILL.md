---
name: scaffold-prd
description: Guide for scaffolding Product Requirements Documents (PRD) in the HRIS project to ensure consistency between Backend and Frontend teams.
---

# Scaffolding PRD (Product Requirements Document)

When the user asks to create or brainstorm a requirement document, PRD, or specifications for a new module/feature, you MUST follow this format and convention.

## 1. Directory & Naming Convention
- **Location:** All PRD documents MUST be saved in the `docs/requirement/` directory.
- **Naming Convention:** Use lowercase with hyphens (e.g., `employee.md`, `attendance-tracking.md`).
- **Do not** use `docs/technical/` for PRDs. `docs/technical/` is strictly for implementation plans and system architecture documents.

## 2. Mandatory PRD Structure (The 6 Pillars)
An Enterprise-Grade PRD must act as a single source of truth for Business, QA, and Engineering. Every PRD MUST contain the following 6 core sections:

### 2.1. Tujuan & Dampak (The "Why")
Explain *why* this module is being built. What business problem does it solve?
*Example: "Mempercepat proses input data dari 10 menit menjadi 2 menit."*

### 2.2. Scope & Out-of-Scope (Batasan Tegas)
Clearly define what is being built and, critically, what is **NOT** being built right now to prevent feature creep.

### 2.3. User Roles & Permissions
Define who will use this feature and their access levels (e.g., Superadmin, HR Manager, Regular Employee). Detail what each role can Read, Write, or Approve.

### 2.4. Kriteria Penerimaan (Acceptance Criteria)
The strict definition of "Done" to prevent debates between QA and Engineering. Use the **Given-When-Then** format to ensure boundaries are black-and-white and easily convertible into unit tests.

### 2.5. Technical & Architectural Constraints
Define the engineering rules for this module.
- **Backend:** DDD isolation rules, strict typing, or soft-delete mandates.
- **Frontend:** Form structure (e.g., Wizard/Multi-step), UI constraints, data masking, or client-side validations.

### 2.6. Dependencies (Ketergantungan)
List what this module relies on to function properly:
- Internal dependencies (e.g., Module X depends on Module Y's API).
- External integrations (e.g., 3rd party APIs, SSO, Payment Gateways).

## 3. Data Schema & Business Rules
After the 6 core pillars, provide the logical breakdown of the data models (Entities):
- **Header:** `## [Entity Name]`
- **Aturan Bisnis:** Unique constraints, relations.
- **Sample Data:** A Markdown table illustrating what the data looks like. Columns should reflect the actual fields (e.g., `id`, `name`, `status`, `created_at`). This helps Frontend developers mock the UI.
