-- Insert a new company into the company table
-- name: CreateCompany :one 
INSERT INTO company (name, description, employees, registered, type)
VALUES ($1, $2, $3, $4, $5) RETURNING *;


-- Update an existing company's information
-- name: UpdateCompany :one 
UPDATE  company
SET name = $1, description = $2 , employees = $3, type = $4
WHERE id = $5 RETURNING *;

-- name: ListCompanies :many
SELECT * FROM company;


-- Delete a company from the company table
-- name: DeleteCompany :exec
DELETE FROM company
WHERE id = $1;


-- Retrieve information about a specific company by its ID
-- name: GetCompany :one
SELECT id, name, description, employees, registered, type
FROM company
WHERE id = $1;


-- Insert a new user into the users table
-- name: InsertUser :one
INSERT INTO users (name)
VALUES ($1) RETURNING *;


-- Retrieve information about a specific user by their name
-- name: GetUserByName :one
SELECT name
FROM users
WHERE name = $1;
