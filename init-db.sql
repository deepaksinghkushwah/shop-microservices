-- Initialize databases for microservices

-- Create auth database
CREATE DATABASE IF NOT EXISTS auth;

-- Create catalog database
CREATE DATABASE IF NOT EXISTS catalog;

-- Create order database
CREATE DATABASE IF NOT EXISTS order;

-- Grant privileges to postgres user
GRANT ALL PRIVILEGES ON DATABASE auth TO postgres;
GRANT ALL PRIVILEGES ON DATABASE catalog TO postgres;
GRANT ALL PRIVILEGES ON DATABASE order TO postgres;
