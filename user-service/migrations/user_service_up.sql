CREATE TABLE IF NOT EXISTS users
(
    id           uuid primary key,
    full_name    varchar                                                          NOT NULL,
    phone_number varchar UNIQUE                                                   NOT NULL,
    password     varchar                                                          NOT NULL,
    role         varchar check ( role in ('CEO', 'TEACHER', 'ADMIN', 'EMPLOYEE')) NOT NULL,
    birth_date   date                                                             NOT NULL,
    gender       boolean                                                          NOT NULL DEFAULT TRUE,
    is_deleted   boolean                                                          NOT NULL DEFAULT FALSE,
    created_at   timestamp                                                                 DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS users_history
(
    id            uuid primary key,
    user_id       uuid references users (id),
    updated_field varchar   NOT NULL,
    old_value     varchar   NOT NULL,
    current_value varchar   NOT NULL,
    created_at    timestamp NOT NULL DEFAULT NOW()
);