module myapp

go 1.23.4

replace github.com/fouched/celeritas => ../celeritas

require github.com/fouched/celeritas v0.0.0-00010101000000-000000000000

require (
	github.com/go-chi/chi/v5 v5.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)
