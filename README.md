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

### Dev evolution

1) first option is to manually create a new CRD yaml file and add it using kubectl create.
   Then create instances of the new type.
   Then we can use the dynamic go client to read and write these types.
   We need to specify the fields in a map, there are no generated go types for the CRD type.

   see the first section
   Defining and creating a Custom Resource
   in https://www.martin-helmich.de/en/blog/kubernetes-crd-client.html

2) I tried using to generate go types to match the CRD in step one and my own clientset as per this tutorial
https://www.martin-helmich.de/en/blog/kubernetes-crd-client.html
it just uses controller-gen tool, which is part of the Kubebuilder framework:
So this approach uses strong typing using generated go types and the static go-client

Add required annotations then run

$ controller-gen object paths=./api/v1/

this adds deep copy methods for all types in v1 folder into zz_generated.deepcopy.go
NOTE: these methods are for use with the static client.
Then the article goes on to write a client-set for the new type to do CRUD opertions on it

A note about clients, creating your own client-set for static types is manually intensive.
it works but is very painful compared to just using the client-go dynamic client.
With this article you install the CRDS yourself and then write custom types.
Kubebuilder and scp-operator will do all of these steps together
you can select Y for Create Resource [y/n]  and N to generate controller (operator).
See section 3

See kubebuilder quick start
https://book.kubebuilder.io/quick-start.html
https://www.openshift.com/blog/kubernetes-operators-best-practices



3)  Next I noticed the sample-operator builds crd types and also generates lister functions to get the custom crd.
https://github.com/kubernetes/sample-controller
https://itnext.io/building-an-operator-for-kubernetes-with-the-sample-controller-b4204be9ad56
do this one with kubebuilder
https://itnext.io/building-an-operator-for-kubernetes-with-kubebuilder-17cbd3f07761
so rather than just trying to generate go types for my crd, I will do the sample-controller steps in this project

kubebuilder init --domain my.domain
kubebuilder create api --group webapp --version v1 --kind Scpcluster
kubebuilder create api --group webapp --version v1 --kind ManagedOperator

// install CRDs into cluster
make install

See config/samples for a uncustomized crd template to create an instance of the type

This calls controller-gen under the covers


TODO
List real operators - need new CRD ManagedOperator
List real instances
Add other personas to match tanzu screen snapshot.
think about scp and tdm interestion
Hide cert stuff for cluster

Out of Scope
- installing operator

Name of operator
List of its CRDs
 - key fields like status

 - CRD to create a resource ?
   input fields and type for ui/cli

- how to identify services it has created 
  CRD type
  CRD key fields like, name, status



