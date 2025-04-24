$env:GOOSE_DRIVER = "postgres"
$env:GOOSE_DBSTRING = "postgres://postgres:postgres@localhost:5432/notesdb?sslmode=disable"
goose -dir internal/db/migrations up
