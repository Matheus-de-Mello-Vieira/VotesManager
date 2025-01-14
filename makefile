setup:
	docker compose up -d postgres kafka redis
	sleep 5

	docker exec --workdir /bin/ -it kafka ./kafka-topics --bootstrap-server localhost:9092 --create --topic votes
	docker exec postgres /bin/psql -h 127.0.0.1 -p 5432 -U postgres -d postgres -f ddl/script.sql
	sleep 5

	docker compose up prodution-frontend voters-frontend voters-register --build

refresh_rough_totals:
	docker exec postgres /bin/psql -h 127.0.0.1 -p 5432 -U postgres -d postgres -c "REFRESH MATERIALIZED VIEW rough_totals"

unit_test:
	( cd repositories ; go test ./... )

load_test:
	( cd repositories ; go run k6/test_load.go )

up_depedencies:
	docker compose up -d postgres kafka

GOPATH = $(shell cd repositories ; go env GOPATH)
swagger:
	cd repositories ; \
		$(GOPATH)/bin/swag init --generalInfo prodution-frontend/main.go --output prodution-frontend/docs --exclude voters-frontend; \
		$(GOPATH)/bin/swag init --generalInfo voters-frontend/main.go --output voters-frontend/docs --exclude prodution-frontend