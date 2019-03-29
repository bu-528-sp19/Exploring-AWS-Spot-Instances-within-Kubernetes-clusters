cd Desktop/CS640_PA1_Kevin_Rodrigues/microservices-demo/deploy/kubernetes/
kubectl create namespace sock-shop
kubectl apply -f complete-demo.yaml
cd ~/metrics-server
kubectl create -f deploy/.


