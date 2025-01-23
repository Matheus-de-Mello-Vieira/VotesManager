# BBB Voting

## Usage
To run everything:
```bash
make setup
```

* With you brownser, so go to:
* **voters-frontend** http://localhost:8080/
* **prodution-frontend**: http://localhost:8081/votes/

### Swagger:
* http://localhost:8080/swagger
* http://localhost:8081/swagger

### Grafana
* Open your browser and navigate to http://localhost:3000.
* Default login credentials
  * Username: `admin`
  * Password: `admin` (you will be prompted to change it on first login).
* There will be a dashboard called Metrics with the CPU and RAM usage for each selected container

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