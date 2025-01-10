# BBB Voting

## Usage
To run everything:
```bash
docker compose up -d postgres
docker compose up -d kafka
```

At first time:
```bash
docker exec --workdir /bin/ -it kafka ./kafka-topics --bootstrap-server localhost:9092 --create --topic test-topic
```
You need to execute the DDL of database too, in `ddl/script.sql`

```
sleep 10
docker compose up
```

* With you brownser, so go to:
* **voters-frontend** http://localhost:8080/
* **prodution-frontend**: http://localhost:8081/votes/