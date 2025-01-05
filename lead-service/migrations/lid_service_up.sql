CREATE TABLE IF NOT EXISTS lead_section
(
    id         serial PRIMARY KEY,
    title      varchar NOT NULL UNIQUE,
    created_at timestamp DEFAULT NOW(),
    company_id int
);

CREATE TABLE IF NOT EXISTS expect_section
(
    id         serial UNIQUE,
    title      varchar NOT NULL UNIQUE,
    created_at timestamp DEFAULT NOW(),
    company_id int
);

CREATE TABLE IF NOT EXISTS set_section
(
    id         serial PRIMARY KEY,
    title      varchar                                               NOT NULL,
    course_id  int                                                   NOT NULL,
    teacher_id uuid                                                  NOT NULL,
    date_type  varchar check (date_type in ('JUFT', 'TOQ', 'OTHER')) NOT NULL,
    days       TEXT[]                                                NOT NULL,
    start_time varchar                                               NOT NULL,
    created_at timestamp DEFAULT NOW(),
    CONSTRAINT valid_days CHECK (array_length(days, 1) > 0 AND days <@
                                                               ARRAY ['DUSHANBA', 'SESHANBA', 'CHORSHANBA', 'PAYSHANBA', 'JUMA', 'SHANBA', 'YAKSHANBA']),
    company_id int
);

CREATE TABLE IF NOT EXISTS lead_user
(
    id           serial PRIMARY KEY,
    phone_number varchar NOT NULL,
    full_name    varchar NOT NULL,
    lead_id      int REFERENCES lead_section (id),
    expect_id    int REFERENCES expect_section (id),
    set_id       int REFERENCES set_section (id) ON DELETE CASCADE,
    comment      varchar,
    created_at   timestamp DEFAULT now(),
    company_id   int
);

CREATE TABLE IF NOT EXISTS lead_source_reports
(
    id         uuid primary key,
    lead_count int,
    source     varchar,
    created_at timestamp DEFAULT NOW(),
    company_id int
);

CREATE TABLE IF NOT EXISTS lead_conversion_reports
(
    id              uuid primary key,
    lead_count      int,
    conversion_date varchar,
    created_at      timestamp DEFAULT NOW(),
    company_id      int
);