test:
	( cd repositories ; go test ./... )
	
up_depedencies:
	docker compose up -d postgres
	docker compose up -d kafka