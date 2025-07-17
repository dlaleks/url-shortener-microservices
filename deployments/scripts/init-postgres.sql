-- Create databases for microservices
CREATE DATABASE url_service;
CREATE DATABASE user_service;

-- Create user for services (optional, for better security)
CREATE USER url_service_user WITH PASSWORD 'url_service_pass';
CREATE USER user_service_user WITH PASSWORD 'user_service_pass';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE url_service TO url_service_user;
GRANT ALL PRIVILEGES ON DATABASE user_service TO user_service_user;

-- Connect to url_service database
\c url_service;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Connect to user_service database
\c user_service;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";