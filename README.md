# BBB Voting

## Usage
To run everything:
```bash
docker compose up -d postgres
docker compose up -d kafka

sleep 10
docker compose up
```

* With you brownser, so go to:
* **voters-frontend** http://localhost:8080/
* **prodution-frontend**: http://localhost:8081/votes/