-- seed.sql — test data for local development
-- Run with: psql $POSTGRES_ADDR -f migrations/seed.sql

-- Wipe existing seed data (safe to re-run)
DELETE FROM metrics_daily WHERE user_id IN (SELECT id FROM users WHERE api_key = 'test-token');
DELETE FROM job_results  WHERE job_id IN (SELECT id FROM jobs WHERE user_id IN (SELECT id FROM users WHERE api_key = 'test-token'));
DELETE FROM jobs         WHERE user_id IN (SELECT id FROM users WHERE api_key = 'test-token');
DELETE FROM events       WHERE user_id IN (SELECT id FROM users WHERE api_key = 'test-token');
DELETE FROM users        WHERE api_key = 'test-token';

-- Test user
INSERT INTO users (id, api_key)
VALUES ('00000000-0000-0000-0000-000000000001', 'test-token');

-- Events
INSERT INTO events (user_id, type, payload_json) VALUES
    ('00000000-0000-0000-0000-000000000001', 'click',    '{"button": "signup"}'),
    ('00000000-0000-0000-0000-000000000001', 'pageview', '{"path": "/dashboard"}'),
    ('00000000-0000-0000-0000-000000000001', 'click',    '{"button": "upgrade"}'),
    ('00000000-0000-0000-0000-000000000001', 'pageview', '{"path": "/settings"}'),
    ('00000000-0000-0000-0000-000000000001', 'click',    '{"button": "logout"}');

-- Jobs
INSERT INTO jobs (user_id, type, status, attempts, max_attempts) VALUES
    ('00000000-0000-0000-0000-000000000001', 'export',  'pending',   0, 3),
    ('00000000-0000-0000-0000-000000000001', 'report',  'done',      1, 3),
    ('00000000-0000-0000-0000-000000000001', 'cleanup', 'failed',    3, 3),
    ('00000000-0000-0000-0000-000000000001', 'export',  'pending',   0, 3),
    ('00000000-0000-0000-0000-000000000001', 'report',  'done',      2, 3);

-- metrics_daily: 30 days of data ending today
INSERT INTO metrics_daily (user_id, day, events_received, jobs_done, jobs_failed)
SELECT
    '00000000-0000-0000-0000-000000000001',
    CURRENT_DATE - (n || ' days')::INTERVAL,
    (random() * 100)::BIGINT,
    (random() * 20)::BIGINT,
    (random() * 5)::BIGINT
FROM generate_series(0, 29) AS n;
