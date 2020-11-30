#### SCP POC Rest Service

This is a REST server that uses in memory data structures to server
Clusters
Services
Factores
for the Clarity POC UI




## Running the example

To run this example, from the root of this project:

```sh

go run ./rest-svr/*.go
```
### Cluster Curl Tests

```

// LOGIN 
curl -i --header "Content-Type: application/json" \
  --request POST \
  --data '{ "email": "developer@vmware.com", "password": "VMware1!"}' \
  http://localhost:8080/api/session

// CURRENT USER


curl --header "Content-Type: application/json" http://localhost:8080/api/currentuser

// GET ALL

curl --header "Content-Type: application/json" http://localhost:8080/api/cluster

// GET ONE

curl --header "Content-Type: application/json" http://localhost:8080/api/cluster/2

// CREATE

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{ "name":"New Cluster", "url":"https://192.168.44.10", "token": "", "cert": "", "certauth": "", "connected": "true"}' \
  http://localhost:8080/api/cluster

// UPDATE

curl --header "Content-Type: application/json" \
  --request PUT \
  --data '{"id":3, "name":"New Cluster Three", "url":"https://192.168.44.10", "connected": "true"}' \
  http://localhost:8080/api/cluster

// DELETE

curl --header "Content-Type: application/json" --request DELETE   http://localhost:8080/api/cluster/3

// CONNECT

curl --header "Content-Type: application/json" \
  --request POST \
  --data '' \
  http://localhost:8080/api/cluster/1/connect

```

### Service Curl Tests

```
// GET ALL

curl --header "Content-Type: application/json" http://localhost:8080/api/service

// GET ONE

curl --header "Content-Type: application/json" http://localhost:8080/api/service/2

// CREATE

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{ "name":"New Cluster", "url":"https://192.168.44.10",  "status": "Active"}' \
  http://localhost:8080/api/service

// UPDATE

curl --header "Content-Type: application/json" \
  --request PUT \
  --data '{"id":3, "name":"New Cluster", "url":"https://192.168.44.10", "status": "Inactive"}' \
  http://localhost:8080/api/service

// DELETE

curl --header "Content-Type: application/json" --request DELETE   http://localhost:8080/api/service/3

```


### Factory Curl Tests

```
// GET ALL

curl --header "Content-Type: application/json" http://localhost:8080/api/factory

// GET ONE

curl --header "Content-Type: application/json" http://localhost:8080/api/factory/2

// CREATE

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{ "name":"New Factory", "url":"https://192.168.44.10", "status": "Active"}' \
  http://localhost:8080/api/factory

// UPDATE

curl --header "Content-Type: application/json" \
  --request PUT \
  --data '{"id":3, "name":"New Factory", "url":"https://192.168.44.10", "status": "Inactive"}' \
  http://localhost:8080/api/factory

// DELETE

curl --header "Content-Type: application/json" --request DELETE   http://localhost:8080/api/factory/3

```

### Docker deployment

docker build -t scp-rest-svr .
docker run --rm -d -p 8080:8080 scp-rest-svr
curl http://localhost:8080/api/cluster

docker rmi $(docker images -f dangling=true -q)

### minikube

$ docker push registry.local/scp-rest-svr:latest
$ kubcectl create -f deployment.yaml
$ kubcectl create -f nodeport.yaml

$ minikube service --url scp-rest-svr
http://192.168.64.4:30180

$  curl http://192.168.64.4:30180/api/cluster
[{"id":1,"name":"Cluster One","url":"http://192.168.0.20","connected":"0001-01-01T00:00:00Z"},{"id":2,"name":"Cluster Two","url":"http://192.168.0.42","connected":"0001-01-01T00:00:00Z"}]

from inside scp-rest-svr can get pod internal ip using dns name
$ nslookup   scp-rest-svr.default.svc.cluster.local
Server:		10.96.0.10
Address:	10.96.0.10:53

Name:	scp-rest-svr.default.svc.cluster.local
Address: 10.98.114.9
