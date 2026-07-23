CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL,
    legal_name VARCHAR(150) NOT NULL,
    npwp VARCHAR(25),
    bpjs_no VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_companies_npwp ON companies (npwp) WHERE npwp IS NOT NULL AND deleted_at IS NULL;

CREATE TABLE branches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies (id),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(150) NOT NULL,
    city VARCHAR(100),
    is_main BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_branches_company_code ON branches (company_id, code) WHERE deleted_at IS NULL;
CREATE INDEX idx_branches_company_id ON branches (company_id);
CREATE UNIQUE INDEX idx_branches_company_main ON branches (company_id) WHERE is_main = true AND deleted_at IS NULL;
