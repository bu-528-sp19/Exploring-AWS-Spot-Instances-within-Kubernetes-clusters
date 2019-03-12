Clone the microservices-demo repo

<pre>git clone https://github.com/suryatej77/microservices-demo
cd microservices-demo</pre>

Start Minikube

<pre> minikube start --memory 4096 </pre>

Deploy the Sock Shop application on Minikube

<pre>cd deploy/kubernetes/
kubectl create namespace sock-shop

kubectl apply -f complete-demo.yaml </pre>

Wait for all the Sock Shop services to start:

<pre> kubectl get pods --namespace="sock-shop" </pre>

Once the application is deployed, navigate to
<br>
http://[minikube_ip]:30001

Get minikube's ip by running:

<pre>minikube ip</pre>
