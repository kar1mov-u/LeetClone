
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'level') THEN
        CREATE TYPE level AS ENUM('easy', 'medium', 'hard');
    END IF;
END$$;

CREATE TABLE problems(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    difficulty level NOT NULL
);

CREATE TABLE example(
    id UUID PRIMARY KEY NOT NULL, 
    problem_id BIGINT NOT NULL, 
    input TEXT NOT NULL, 
    output TEXT NOT NULL, 
    explanation TEXT,
    CONSTRAINT fk_problem FOREIGN KEY(problem_id) REFERENCES problems(id)
);

CREATE TABLE testcase(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id BIGINT NOT NULL, 
    input TEXT NOT NULL, 
    output TEXT NOT NULL,
    CONSTRAINT fk_problem FOREIGN KEY(problem_id) REFERENCES problems(id)
);

