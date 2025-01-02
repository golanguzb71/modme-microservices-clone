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
    is_discounted  boolean                                                  DEFAULT FALSE,
    discount_owner varchar CHECK ( discount_owner in ('TEACHER', 'CENTER')) DEFAULT 'TEACHER',
    price          float                                                        NOT NULL,
    group_id       bigint references groups (id),
    student_id     uuid                                                         NOT NULL,
    teacher_id     uuid                                                         NOT NULL,
    attend_date    date                                                         NOT NULL,
    status         int                                                          NOT NULL,
    created_at     timestamp                                                DEFAULT NOW(),
    created_by     uuid                                                         NOT NULL,
    creator_role   varchar CHECK ( creator_role in ('ADMIN', 'CEO', 'TEACHER')) NOT NULL,
    PRIMARY KEY (group_id, student_id, attend_date)
);

CREATE TABLE IF NOT EXISTS students
(
    id                 uuid PRIMARY KEY,
    name               varchar NOT NULL,
    phone              varchar NOT NULL,
    date_of_birth      date                                                  default '2000-12-12',
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
    id            uuid primary key,
    group_id      bigint references groups (id) NOT NULL,
    field         varchar                       NOT NULL,
    old_value     varchar                       NOT NULL,
    current_value varchar                       NOT NULL,
    created_at    timestamp DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS student_history
(
    id            uuid primary key,
    student_id    uuid references students (id) NOT NULL,
    field         varchar                       NOT NULL,
    old_value     varchar                       NOT NULL,
    current_value varchar                       NOT NULL,
    created_at    timestamp DEFAULT NOW()
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
    id                  uuid primary key,
    group_student_id    uuid references group_students (id)                                  NOT NULL,
    student_id          uuid references students (id)                                        NOT NULL,
    group_id            bigint references groups (id)                                        NOT NULL,
    old_condition       varchar check ( old_condition in ('FREEZE', 'ACTIVE', 'DELETE'))     NOT NULL,
    current_condition   varchar check ( current_condition in ('FREEZE', 'ACTIVE', 'DELETE')) NOT NULL,
    is_eliminated_trial bool                                                                          DEFAULT FALSE,
    specific_date          date                                                                 NOT NULL DEFAULT NOW(),
    return_the_money    boolean                                                              NOT NULL DEFAULT FALSE,
    created_at          timestamp                                                                     DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_attendance_group_date ON attendance (group_id, attend_date);
CREATE INDEX IF NOT EXISTS idx_group_students_group ON group_students (group_id);
CREATE OR REPLACE FUNCTION log_group_update()
    RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.name IS DISTINCT FROM OLD.name THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'name', COALESCE(OLD.name, ''), COALESCE(NEW.name, ''), NOW());
    END IF;

    IF NEW.course_id IS DISTINCT FROM OLD.course_id THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'course_id', COALESCE(OLD.course_id::text, ''),
                COALESCE(NEW.course_id::text, ''), NOW());
    END IF;

    IF NEW.teacher_id IS DISTINCT FROM OLD.teacher_id THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'teacher_id', COALESCE(OLD.teacher_id::text, ''),
                COALESCE(NEW.teacher_id::text, ''), NOW());
    END IF;

    IF NEW.room_id IS DISTINCT FROM OLD.room_id THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'room_id', COALESCE(OLD.room_id::text, ''), COALESCE(NEW.room_id::text, ''),
                NOW());
    END IF;

    IF NEW.date_type IS DISTINCT FROM OLD.date_type THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'date_type', COALESCE(OLD.date_type, ''), COALESCE(NEW.date_type, ''),
                NOW());
    END IF;

    IF NEW.start_time IS DISTINCT FROM OLD.start_time THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'start_time', COALESCE(OLD.start_time::text, ''),
                COALESCE(NEW.start_time::text, ''), NOW());
    END IF;

    IF NEW.start_date IS DISTINCT FROM OLD.start_date THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'start_date', COALESCE(OLD.start_date::text, ''),
                COALESCE(NEW.start_date::text, ''), NOW());
    END IF;

    IF NEW.end_date IS DISTINCT FROM OLD.end_date THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'end_date', COALESCE(OLD.end_date::text, ''),
                COALESCE(NEW.end_date::text, ''), NOW());
    END IF;

    IF NEW.is_archived IS DISTINCT FROM OLD.is_archived THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'is_archived', COALESCE(OLD.is_archived::text, ''),
                COALESCE(NEW.is_archived::text, ''), NOW());
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_group_update
    AFTER UPDATE
    ON groups
    FOR EACH ROW
EXECUTE FUNCTION log_group_update();


CREATE OR REPLACE FUNCTION log_student_update()
    RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.name IS DISTINCT FROM OLD.name THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'name', COALESCE(OLD.name, ''), COALESCE(NEW.name, ''), NOW());
    END IF;

    IF NEW.phone IS DISTINCT FROM OLD.phone THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'phone', COALESCE(OLD.phone, ''), COALESCE(NEW.phone, ''), NOW());
    END IF;

    IF NEW.date_of_birth IS DISTINCT FROM OLD.date_of_birth THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'date_of_birth', COALESCE(OLD.date_of_birth::text, ''),
                COALESCE(NEW.date_of_birth::text, ''), NOW());
    END IF;

    IF NEW.condition IS DISTINCT FROM OLD.condition THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'condition', COALESCE(OLD.condition, ''), COALESCE(NEW.condition, ''),
                NOW());
    END IF;

    IF NEW.additional_contact IS DISTINCT FROM OLD.additional_contact THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'additional_contact', COALESCE(OLD.additional_contact, ''),
                COALESCE(NEW.additional_contact, ''), NOW());
    END IF;

    IF NEW.address IS DISTINCT FROM OLD.address THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'address', COALESCE(OLD.address, ''), COALESCE(NEW.address, ''), NOW());
    END IF;

    IF NEW.telegram_username IS DISTINCT FROM OLD.telegram_username THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'telegram_username', COALESCE(OLD.telegram_username, ''),
                COALESCE(NEW.telegram_username, ''), NOW());
    END IF;

    IF NEW.passport_id IS DISTINCT FROM OLD.passport_id THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'passport_id', COALESCE(OLD.passport_id, ''), COALESCE(NEW.passport_id, ''),
                NOW());
    END IF;

    IF NEW.gender IS DISTINCT FROM OLD.gender THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'gender', COALESCE(OLD.gender::text, ''), COALESCE(NEW.gender::text, ''),
                NOW());
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_student_update
    AFTER UPDATE
    ON students
    FOR EACH ROW
EXECUTE FUNCTION log_student_update();
