-- Create the replicator user with replication privileges
CREATE USER replicator WITH REPLICATION ENCRYPTED PASSWORD 'replicator_password';

-- Create a physical replication slot
SELECT pg_create_physical_replication_slot('replication_slot');

-- Grant SELECT on all tables to the replicator user
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO replicator;

-- Grant SELECT on all sequences to the replicator user
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON SEQUENCES TO replicator;

-- Grant SELECT on all functions to the replicator user
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT EXECUTE ON FUNCTIONS TO replicator;
