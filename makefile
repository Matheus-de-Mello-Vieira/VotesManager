unit_test:
	( cd repositories ; go test ./... )

load_test:
	( cd repositories ; go run k6/test_load.go )

up_depedencies:
	docker compose up -d postgres
	docker compose up -d kafka

