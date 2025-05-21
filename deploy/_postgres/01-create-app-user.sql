CREATE EXTENSION IF NOT EXISTS pg_cron;

DO
$$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'return_zero_app') THEN
      CREATE USER return_zero_app WITH PASSWORD 'strong_password_here';
   END IF;
END
$$;

GRANT CONNECT ON DATABASE returnzero TO return_zero_app;

GRANT USAGE ON SCHEMA public TO return_zero_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO return_zero_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO return_zero_app;
GRANT USAGE ON SCHEMA cron TO return_zero_app;
GRANT EXECUTE ON FUNCTION cron.schedule(text, text, text) TO return_zero_app;

ALTER DEFAULT PRIVILEGES IN SCHEMA public 
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO return_zero_app;
    
ALTER DEFAULT PRIVILEGES IN SCHEMA public 
    GRANT USAGE, SELECT ON SEQUENCES TO return_zero_app;
