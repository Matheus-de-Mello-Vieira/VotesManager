# BBB Voting

## Usage
To run everything:
```bash
make setup
```

* With you brownser, so go to:
* **voters-frontend** http://localhost:8080/
* **prodution-frontend**: http://localhost:8081/votes/

swagger:
* http://localhost:8080/swagger
* http://localhost:8081/swagger

### Kubenetes

You can also run it on Kubernetes, I have provided a separated markdown file with the instructions: `./kubernetes/commands.md`

### Test
Unit tests:
```bash
make unit_test
```

Load test:
```bash
make load_test
```