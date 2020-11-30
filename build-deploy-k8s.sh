docker build -t registry.local/scp-rest-svr:latest .

kubectl delete deployment scp-rest-svr
kubectl create -f deployment.yaml
