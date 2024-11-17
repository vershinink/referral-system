CREATE USER tester WITH PASSWORD 'tester';
CREATE DATABASE tests;
ALTER DATABASE tests OWNER TO tester;
GRANT ALL ON DATABASE tests TO tester;
\c tests;

SET ROLE tester;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    passwd BYTEA NOT NULL,
    created BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS codes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    owner INTEGER REFERENCES users (id) NOT NULL,
    created BIGINT NOT NULL,
    expired BIGINT NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS referrals (
    id INTEGER REFERENCES users (id) NOT NULL,
    referrer_id INTEGER REFERENCES users (id) NOT NULL,
    code_id INTEGER REFERENCES codes (id) NOT NULL
);

-- CREATE USER tester WITH PASSWORD 'tester';



--ALTER TABLE public.referrals, public.codes, public.users OWNER TO tester;
--GRANT ALL ON ALL TABLES IN SCHEMA public TO tester;
--ALTER ALL SEQUENCES IN SCHEMA public OWNER TO tester;
--GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO tester;
-- ALTER ALL FUNCTIONS IN SCHEMA public OWNER TO tester;
-- GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO tester;