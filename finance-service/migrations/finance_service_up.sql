CREATE TABLE IF NOT EXISTS student_discount
(
    student_id uuid             NOT NULL,
    discount   double precision NOT NULL,
    group_id   bigint           NOT NULL,
    comment    varchar          NOT NULL,
    created_at timestamp DEFAULT NOW(),
    PRIMARY KEY (student_id, group_id)
);
CREATE TABLE IF NOT EXISTS category
(
    id          serial primary key,
    name        varchar NOT NULL,
    description varchar not null
);

CREATE TABLE IF NOT EXISTS expense
(
    id             uuid primary key,
    title          varchar                                                       NOT NULL,
    user_id        uuid,
    category_id    int,
    expense_type   varchar check ( expense_type in ('USER', 'CATEGORY')),
    sum            double precision                                              NOT NULL,
    given_date     date                                                          NOT NULL,
    created_at     timestamp default NOW(),
    created_by     uuid                                                          NOT NULL,
    payment_method varchar check ( payment_method in ('CASH', 'CLICK', 'PAYME')) NOT NULL
);

CREATE TABLE IF NOT EXISTS student_finance
(
    id           uuid primary key,
    student_id   uuid             NOT NULL,
    payment_type varchar          NOT NULL,
    finance_type varchar check ( finance_type in ('KIRIM', 'CHIQIM')),
    amount       double precision NOT NULL,
    given_date   date             NOT NULL,
    comment      varchar          NOT NULL,
    created_at   timestamp DEFAULT NOW(),
    created_by   uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS teacher_finance
(
    id         uuid primary key,
    group_id   bigint           NOT NULL,
    student_id uuid             NOT NULL,
    sum        DOUBLE PRECISION NOT NULL,
    created_at timestamp DEFAULT NOW()
)