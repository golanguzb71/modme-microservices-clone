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

CREATE TABLE IF NOT EXISTS transfer_lesson
(
    id            uuid PRIMARY KEY,
    group_id      bigint references groups (id) NOT NULL,
    real_date     date                          NOT NULL,
    transfer_date date                          NOT NULL
);

CREATE TABLE IF NOT EXISTS group_students
(
    id         uuid PRIMARY KEY,
    group_id   bigint references groups (id) NOT NULL,
    student_id uuid                         NOT NULL,
    condition  varchar check ( condition in ('FREEZE', 'ACTIVE', 'DELETE')) DEFAULT 'FREEZE',
    created_at timestamp                                                   DEFAULT NOW(),
    created_by uuid                         NOT NULL
);