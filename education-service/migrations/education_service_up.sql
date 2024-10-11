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
    id            uuid primary key,
    name          varchar not null,
    phone         varchar not null,
    date_of_birth date,
    balance       double precision                                     DEFAULT 0,
    condition     varchar check ( condition in ('ACTIVE', 'ARCHIVED')) DEFAULT 'ACTIVE',
    gender        boolean,
    created_at    timestamp                                            DEFAULT now()
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
    id         uuid PRIMARY KEY,
    group_id   bigint references groups (id) NOT NULL,
    student_id uuid                          NOT NULL,
    condition  varchar check ( condition in ('FREEZE', 'ACTIVE', 'DELETE')) DEFAULT 'FREEZE',
    created_at timestamp                                                    DEFAULT NOW(),
    created_by uuid                          NOT NULL
);

CREATE TABLE IF NOT EXISTS group_student_condition_history
(
    id               uuid primary key,
    group_student_id uuid references group_students (id)                          NOT NULL,
    student_id       uuid references students (id)                                NOT NULL,
    group_id         bigint references groups (id)                                  NOT NULL,
    condition        varchar check ( condition in ('FREEZE', 'ACTIVE', 'DELETE')) NOT NULL,
    created_at       timestamp DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION log_group_update()
    RETURNS TRIGGER AS
$$
DECLARE
    description TEXT := ' ';
BEGIN
    IF NEW.name IS DISTINCT FROM OLD.name THEN
        description := description || 'name: ' || OLD.name || ' => ' || NEW.name || '; ';
    END IF;
    IF NEW.course_id IS DISTINCT FROM OLD.course_id THEN
        description := description || 'course_id: ' || OLD.course_id || ' => ' || NEW.course_id || '; ';
    END IF;
    IF NEW.teacher_id IS DISTINCT FROM OLD.teacher_id THEN
        description := description || 'teacher_id: ' || OLD.teacher_id || ' => ' || NEW.teacher_id || '; ';
    END IF;
    IF NEW.room_id IS DISTINCT FROM OLD.room_id THEN
        description := description || 'room_id: ' || COALESCE(OLD.room_id::text, 'NULL') || ' => ' ||
                       COALESCE(NEW.room_id::text, 'NULL') || '; ';
    END IF;
    IF NEW.date_type IS DISTINCT FROM OLD.date_type THEN
        description := description || 'date_type: ' || OLD.date_type || ' => ' || NEW.date_type || '; ';
    END IF;
    IF NEW.days IS DISTINCT FROM OLD.days THEN
        description := description || 'days: ' || array_to_string(OLD.days, ', ') || ' => ' ||
                       array_to_string(NEW.days, ', ') || '; ';
    END IF;
    IF NEW.start_time IS DISTINCT FROM OLD.start_time THEN
        description := description || 'start_time: ' || OLD.start_time || ' => ' || NEW.start_time || '; ';
    END IF;
    IF NEW.start_date IS DISTINCT FROM OLD.start_date THEN
        description := description || 'start_date: ' || COALESCE(OLD.start_date::text, 'NULL') || ' => ' ||
                       COALESCE(NEW.start_date::text, 'NULL') || '; ';
    END IF;
    IF NEW.end_date IS DISTINCT FROM OLD.end_date THEN
        description := description || 'end_date: ' || COALESCE(OLD.end_date::text, 'NULL') || ' => ' ||
                       COALESCE(NEW.end_date::text, 'NULL') || '; ';
    END IF;
    IF NEW.is_archived IS DISTINCT FROM OLD.is_archived THEN
        description := description || 'is_archived: ' || OLD.is_archived || ' => ' || NEW.is_archived || '; ';
    END IF;

    INSERT INTO group_history (id, group_id, description, created_at)
    VALUES (gen_random_uuid(), NEW.id, description, NOW());

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION log_student_update()
    RETURNS TRIGGER AS
$$
DECLARE
    description TEXT := ' ';
BEGIN
    IF NEW.name IS DISTINCT FROM OLD.name THEN
        description := description || 'name: ' || OLD.name || ' => ' || NEW.name || '; ';
    END IF;
    IF NEW.phone IS DISTINCT FROM OLD.phone THEN
        description := description || 'phone: ' || OLD.phone || ' => ' || NEW.phone || '; ';
    END IF;
    IF NEW.date_of_birth IS DISTINCT FROM OLD.date_of_birth THEN
        description := description || 'date_of_birth: ' || COALESCE(OLD.date_of_birth::text, 'NULL') || ' => ' ||
                       COALESCE(NEW.date_of_birth::text, 'NULL') || '; ';
    END IF;
    IF NEW.gender IS DISTINCT FROM OLD.gender THEN
        description := description || 'gender: ' || COALESCE(OLD.gender::text, 'NULL') || ' => ' ||
                       COALESCE(NEW.gender::text, 'NULL') || '; ';
    END IF;
    IF NEW.condition IS DISTINCT FROM OLD.condition THEN
        description := description || 'condition: ' || OLD.condition || ' => ' || NEW.condition || '; ';
    END IF;

    INSERT INTO student_history (id, student_id, description, created_at)
    VALUES (gen_random_uuid(), NEW.id, description, NOW());

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION log_group_student_condition_update()
    RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.condition IS DISTINCT FROM OLD.condition THEN
        INSERT INTO group_student_condition_history (id, group_student_id, condition, student_id, group_id, created_at)
        VALUES (gen_random_uuid(), NEW.id, NEW.condition, NEW.student_id, NEW.group_id, NOW());
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trigger_group_update
    AFTER UPDATE
    ON groups
    FOR EACH ROW
EXECUTE FUNCTION log_group_update();

CREATE OR REPLACE TRIGGER trigger_student_update
    AFTER UPDATE
    ON students
    FOR EACH ROW
EXECUTE FUNCTION log_student_update();

CREATE OR REPLACE TRIGGER trigger_group_student_condition_update
    AFTER UPDATE
    ON group_students
    FOR EACH ROW
    WHEN (OLD.condition IS DISTINCT FROM NEW.condition)
EXECUTE FUNCTION log_group_student_condition_update();

CREATE INDEX IF NOT EXISTS idx_attendance_group_date ON attendance (group_id, attend_date);
CREATE INDEX IF NOT EXISTS idx_group_students_group ON group_students (group_id);