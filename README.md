# BBB Voting

## Usage
To run everything:
```bash
make setup
```

To refresh the rough totals (shown after someone voted), run:
```bash
make refresh_rough_totals
```

* With you brownser, so go to:
* **voters-frontend** http://localhost:8080/
* **prodution-frontend**: http://localhost:8081/votes/

### Test
Unit tests:
```bash
make unit_test
```

Load test:
```bash
make load_test
```