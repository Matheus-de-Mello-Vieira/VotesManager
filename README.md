# BBB Voting

## Usage
For the first time:
```bash
cp .env.example .env
```

To run the depedencies:
```bash
docker compose up -d postgres
```

To run the components:
```bash
go run voters-frontend/main.go
```