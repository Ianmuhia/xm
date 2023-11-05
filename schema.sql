CREATE SCHEMA IF NOT EXISTS companies;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE company_type AS ENUM (
    'Corporations',
    'NonProfit',
    'Cooperative',
    'SoleProprietorship'
);

CREATE TABLE company (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(15) NOT NULL UNIQUE,
    description VARCHAR(3000),
    employees INT NOT NULL,
    registered BOOLEAN NOT NULL,
    type company_type NOT NULL
);

CREATE TABLE users (
    name VARCHAR(255) NOT NULL
)