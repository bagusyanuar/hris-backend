CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies (id),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(150) NOT NULL,
    parent_id UUID REFERENCES departments (id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_departments_company_code ON departments (company_id, code) WHERE deleted_at IS NULL;
CREATE INDEX idx_departments_company_id ON departments (company_id);
CREATE INDEX idx_departments_parent_id ON departments (parent_id);

CREATE TABLE job_titles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies (id),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    grade_level INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_job_titles_company_code ON job_titles (company_id, code) WHERE deleted_at IS NULL;
CREATE INDEX idx_job_titles_company_id ON job_titles (company_id);

CREATE TABLE job_positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies (id),
    department_id UUID NOT NULL REFERENCES departments (id),
    job_title_id UUID NOT NULL REFERENCES job_titles (id),
    name VARCHAR(150) NOT NULL,
    reports_to_id UUID REFERENCES job_positions (id),
    headcount_quota INT NOT NULL DEFAULT 1 CHECK (headcount_quota >= 1),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_job_positions_company_id ON job_positions (company_id);
CREATE INDEX idx_job_positions_department_id ON job_positions (department_id);
CREATE INDEX idx_job_positions_job_title_id ON job_positions (job_title_id);
CREATE INDEX idx_job_positions_reports_to_id ON job_positions (reports_to_id);
