CREATE table IF NOT EXISTS rooms
(
    id       serial primary key,
    title    varchar NOT NULL,
    capacity int     NOT NULL
);


CREATE TABLE IF NOT EXISTS courses
(
    id              serial primary key,
    title           varchar                                 NOT NULL,
    duration_lesson int                                     NOT NULL,
    course_duration int                                     NOT NULL,
    price           double precision check ( price > 5000 ) NOT NULL,
    description     text
);


CREATE TABLE IF NOT EXISTS groups
(
    id          bigserial PRIMARY KEY,
    name        varchar                                                NOT NULL,
    course_id   int                                                    NOT NULL,
    teacher_id  uuid                                                   NOT NULL,
    room_id     int references rooms (id),
    date_type   varchar check (date_type in ('JUFT', 'TOQ', 'BOSHQA')) NOT NULL,
    days        TEXT[]                                                 NOT NULL,
    start_time  varchar                                                NOT NULL,
    start_date  date                                                   NOT NULL,
    end_date    date                                                   NOT NULL,
    is_archived boolean   DEFAULT FALSE                                NOT NULL,
    created_at  timestamp DEFAULT NOW(),
    CONSTRAINT valid_days CHECK (array_length(days, 1) > 0 AND days <@
                                                               ARRAY ['DUSHANBA', 'SESHANBA', 'CHORSHANBA', 'PAYSHANBA', 'JUMA', 'SHANBA', 'YAKSHANBA'])
);

CREATE TABLE IF NOT EXISTS attendance
(
    group_id    bigint references groups (id),
    student_id  uuid NOT NULL,
    teacher_id  uuid NOT NULL,
    attend_date date NOT NULL,
    status      int  NOT NULL,
    created_at  timestamp DEFAULT NOW(),
    PRIMARY KEY (group_id, student_id, attend_date)
);

CREATE TABLE IF NOT EXISTS students
(
    id                 uuid PRIMARY KEY,
    name               varchar NOT NULL,
    phone              varchar NOT NULL,
    date_of_birth      date,
    balance            double precision                                      DEFAULT 0,
    condition          varchar CHECK ( condition IN ('ACTIVE', 'ARCHIVED') ) DEFAULT 'ACTIVE',
    additional_contact varchar,
    address            varchar,
    telegram_username  varchar,
    passport_id        varchar,
    gender             boolean,
    created_at         timestamp                                             DEFAULT now()
);

CREATE TABLE IF NOT EXISTS student_note
(
    id         uuid primary key,
    student_id uuid references students (id) NOT NULL,
    comment    text                          NOT NULL,
    created_at timestamp DEFAULT NOW(),
    created_by uuid
);

CREATE TABLE IF NOT EXISTS group_history
(
    id          uuid primary key,
    group_id    bigint references groups (id) NOT NULL,
    description text                          NOT NULL,
    created_at  timestamp DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS student_history
(
    id          uuid primary key,
    student_id  uuid references students (id) NOT NULL,
    description text                          NOT NULL,
    created_at  timestamp DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS transfer_lesson
(
    id            uuid PRIMARY KEY,
    group_id      bigint references groups (id) NOT NULL,
    real_date     date                          NOT NULL,
    transfer_date date                          NOT NULL
);

CREATE TABLE IF NOT EXISTS group_students
(
    id                 uuid PRIMARY KEY,
    group_id           bigint references groups (id) NOT NULL,
    student_id         uuid                          NOT NULL,
    condition          varchar check ( condition in ('FREEZE', 'ACTIVE', 'DELETE')) DEFAULT 'FREEZE',
    last_specific_date date                          NOT NULL                       DEFAULT NOW(),
    created_at         timestamp                                                    DEFAULT NOW(),
    created_by         uuid                          NOT NULL,
    UNIQUE (group_id, student_id)
);

CREATE TABLE IF NOT EXISTS group_student_condition_history
(
    id                uuid primary key,
    group_student_id  uuid references group_students (id)                                  NOT NULL,
    student_id        uuid references students (id)                                        NOT NULL,
    group_id          bigint references groups (id)                                        NOT NULL,
    old_condition     varchar check ( old_condition in ('FREEZE', 'ACTIVE', 'DELETE'))     NOT NULL,
    current_condition varchar check ( current_condition in ('FREEZE', 'ACTIVE', 'DELETE')) NOT NULL,
    specific_date     date                                                                 NOT NULL DEFAULT NOW(),
    return_the_money  boolean                                                              NOT NULL DEFAULT FALSE,
    created_at        timestamp                                                                     DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_attendance_group_date ON attendance (group_id, attend_date);
CREATE INDEX IF NOT EXISTS idx_group_students_group ON group_students (group_id);