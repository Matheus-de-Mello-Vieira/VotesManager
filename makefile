setup:
	$(MAKE) setup_depedencies
	sleep 10

	docker exec --workdir /bin/ -it kafka ./kafka-topics --bootstrap-server localhost:9092 --create --topic votes
	docker exec postgres /bin/psql -h 127.0.0.1 -p 5432 -U postgres -d postgres -f ddl/script.sql

	$(MAKE) setup_main

tear_down:
	docker compose down

setup_depedencies:
	docker compose up -d postgres kafka redis

setup_main:
	docker compose up prodution-frontend voters-frontend voters-register --build

unit_test:
	( cd repositories ; go test ./... )

load_test:
	( docker compose up k6 --build )

GOPATH = $(shell cd repositories ; go env GOPATH)
swagger:
	cd repositories ; \
		$(GOPATH)/bin/swag init --generalInfo prodution-frontend/main.go --output prodution-frontend/docs --exclude voters-frontend; \
		$(GOPATH)/bin/swag init --generalInfo voters-frontend/main.go --output voters-frontend/docs --exclude prodution-frontend
