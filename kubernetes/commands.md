# Setup 

If you are using minikube, execute:

```bash
minikube start --cpus='4' --driver=docker
eval $(minikube docker-env) && \
    docker build -t voters-frontend:latest --build-arg MAIN_PATH="voters-frontend/main.go" . && \
    docker build -t prodution-frontend:latest --build-arg MAIN_PATH="prodution-frontend/main.go" . && \
    docker build -t votes-register:latest --build-arg MAIN_PATH="votes-register/main.go" .
```

Now, execute:

```bash
kubectl create configmap postgres-init-sql --from-file=ddl/script.sql

kubectl apply -f kubernetes/zookeeper.yaml
kubectl apply -f kubernetes/kafka.yaml
kubectl apply -f kubernetes/postgresql.yaml
kubectl apply -f kubernetes/redis.yaml

postgresPod=$(kubectl get pods -l app=postgres --no-headers -o custom-columns=":metadata.name")
kafkaPod=$(kubectl get pods -l app=kafka-broker --no-headers -o custom-columns=":metadata.name")
redisPod=$(kubectl get pods -l app=redis --no-headers -o custom-columns=":metadata.name")

echo "wainting depedences"
kubectl wait --for=condition=Ready pod/$postgresPod --timeout="200s"
kubectl wait --for=condition=Ready pod/$kafkaPod --timeout="200s"
kubectl wait --for=condition=Ready pod/$redisPod --timeout="200s"

kubectl exec -it $postgresPod -- psql -U postgres -d postgres -f ddl/script.sql
kubectl exec -it $kafkaPod -- kafka-topics --bootstrap-server localhost:9092 --create --topic votes

kubectl apply -f kubernetes/voters-frontend.yaml
kubectl apply -f kubernetes/prodution-frontend.yaml
kubectl apply -f kubernetes/votes-register.yaml

echo "wainting system"
kubectl wait --for=condition=Ready pod/$(kubectl get pods -l app=voters-frontend --no-headers -o custom-columns=":metadata.name")
kubectl wait --for=condition=Ready pod/$(kubectl get pods -l app=prodution-frontend --no-headers -o custom-columns=":metadata.name")
kubectl wait --for=condition=Ready pod/$(kubectl get pods -l app=votes-register --no-headers -o custom-columns=":metadata.name")
```

Get acess (minikube):

```bash
minikube service voters-frontend prodution-frontend
```

# Utilities

## Deploy new code

### voters-frontend

```bash
eval $(minikube docker-env) && \
docker build -t voters-frontend:latest --build-arg MAIN_PATH="voters-frontend/main.go" . && \
kubectl rollout restart deployment voters-frontend
```

### prodution-frontend

```bash
eval $(minikube docker-env) && \
docker build -t prodution-frontend:latest --build-arg MAIN_PATH="voters-frontend/main.go" . && \
kubectl rollout restart deployment prodution-frontend
```

### voters-frontend

```bash
eval $(minikube docker-env) && \
docker build -t voters-frontend:latest --build-arg MAIN_PATH="prodution-frontend/main.go" . && \
kubectl rollout restart deployment voters-frontend
```

## Get logs

### voters-frontend

```bash
kubectl logs $(kubectl get pods -l app=voters-frontend --no-headers -o custom-columns=":metadata.name")
```

### prodution-frontend

```bash
kubectl logs $(kubectl get pods -l app=prodution-frontend --no-headers -o custom-columns=":metadata.name")
```

### votes-register

```bash
kubectl logs $(kubectl get pods -l app=votes-register --no-headers -o custom-columns=":metadata.name")
```

## Get number of votes on database

```bash
kubectl exec -it $postgresPod -- psql -U postgres -d postgres -c 'select count(*) from votes;'
```

# load test

obs.: this test can fail due to docker network limitations

run:

```bash
eval $(minikube docker-env) &&
docker build -t load-test:latest --build-arg MAIN_PATH="k6/test_load.go" k6/. && \
kubectl apply -f kubernetes/load_test.yaml
```

get logs:

```bash
kubectl logs $(kubectl get pods -l app=load-test --no-headers -o custom-columns=":metadata.name")
```

delete

```bash
kubectl delete -f kubernetes/load_test.yaml
```