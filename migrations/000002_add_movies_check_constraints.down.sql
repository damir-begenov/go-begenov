ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_runtime_check;
ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_year_check;
ALTER TABLE movies DROP CONSTRAINT IF EXISTS genres_length_check;
migrate -source file://Users/damirbegenov/Downloads/greenlight-2/migrations -database "postgres://localhost:5432/postgres:0000/postgres?sslmode=disable" up

