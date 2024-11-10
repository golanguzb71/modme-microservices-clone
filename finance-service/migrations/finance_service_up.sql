CREATE TABLE IF NOT EXISTS student_discount
(
    student_id  uuid             NOT NULL,
    discount    double precision NOT NULL,
    group_id    bigint           NOT NULL,
    comment     varchar          NOT NULL,
    start_at    date             NOT NULL,
    end_at      date             NOT NULL,
    withTeacher boolean          NOT NULL,
    created_at  timestamp DEFAULT NOW(),
    PRIMARY KEY (student_id, group_id)
);

CREATE TABLE student_discount_history
(
    id          uuid PRIMARY KEY,
    student_id  uuid             NOT NULL,
    group_id    bigint           NOT NULL,
    discount    double precision NOT NULL,
    start_at    date             NOT NULL,
    end_at      date             NOT NULL,
    withTeacher boolean          NOT NULL,
    comment     varchar          NOT NULL,
    action      varchar          NOT NULL,
    created_at  timestamp default now()
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

CREATE TABLE IF NOT EXISTS student_payments
(
    id              uuid primary key,
    student_id      uuid                                                                   NOT NULL,
    method          varchar check ( student_payments.method in ('CLICK', 'CASH', 'PAYME')) NOT NULL,
    amount          double precision                                                       NOT NULL,
    given_date      date                                                                   NOT NULL,
    comment         varchar                                                                NOT NULL,
    created_at      timestamp                                                                       DEFAULT NOW(),
    payment_type    varchar check ( payment_type in ('ADD', 'TAKE_OFF', 'REFUND') )        NOT NULL,
    created_by_id   uuid                                                                   NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    created_by_name varchar                                                                NOT NULL,
    group_id        bigint
);

CREATE TABLE IF NOT EXISTS teacher_salary
(
    teacher_id        uuid PRIMARY KEY,
    salary_type       varchar CHECK (salary_type IN ('PERCENT', 'FIXED')) NOT NULL,
    salary_type_count double precision CHECK (
        CASE
            WHEN salary_type = 'PERCENT' THEN salary_type_count BETWEEN 1 AND 100
            ELSE TRUE
            END
        ),
    created_at        timestamp DEFAULT NOW()
);

