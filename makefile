setup:
	docker compose up -d postgres
	docker compose up -d kafka
	sleep 5

	docker exec --workdir /bin/ -it kafka ./kafka-topics --bootstrap-server localhost:9092 --create --topic votes
	docker exec postgres /bin/psql -h 127.0.0.1 -p 5432 -U postgres -d postgres -f ddl/script.sql
	sleep 5

	docker compose up

unit_test:
	( cd repositories ; go test ./... )

load_test:
	( cd repositories ; go run k6/test_load.go )

up_depedencies:
	docker compose up -d postgres
	docker compose up -d kafka

