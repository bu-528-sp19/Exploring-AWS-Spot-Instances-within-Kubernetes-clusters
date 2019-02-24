Install and setup Hyper-V hypervisor for Windows. (Can also use Virtualbox or any such hypervisor supported by Minikube)

Install and setup Minikube and Kubectl with Hyper-V.

Run these commands to setup a demo php application to be deployed.

Start Minikube
<pre>minikube start</pre>

Apply the Redis Master Deployment from the redis-master-deployment.yaml file
<pre>kubectl apply -f redis-master-deployment.yaml</pre>

Query the list of Pods to verify that the Redis Master Pod is running
<pre>kubectl get pods</pre>

Apply the Redis Master Service from the redis-master-service.yaml file
<pre>kubectl apply -f redis-master-service.yaml</pre>

Query the list of Services to verify that the Redis Master Service is running
<pre>kubectl get service</pre>

Apply the Redis Slave Deployment from the redis-slave-deployment.yaml file
<pre>kubectl apply -f redis-slave-deployment.yaml</pre>

Query the list of Pods to verify that the Redis Slave Pods are running
<pre>kubectl get pods</pre>

Apply the Redis Slave Service from the redis-slave-service.yaml file
<pre>kubectl apply -f redis-slave-service.yaml</pre>

Query the list of Services to verify that the Redis slave service is running
<pre>kubectl get services</pre>

Apply the frontend Deployment from the frontend-deployment.yaml file
<pre>kubectl apply -f frontend-deployment.yaml</pre>

Query the list of Pods to verify that the three frontend replicas are running
<pre>kubectl get pods -l app=guestbook -l tier=frontend</pre>

Apply the frontend Service from the frontend-service.yaml file
<pre>kubectl apply -f frontend-service.yaml</pre>

Query the list of Services to verify that the frontend Service is running
<pre>kubectl get services</pre>

Run the following command to get the IP address for the frontend Service
<pre>minikube service frontend --url</pre>

Copy the IP address to your browser to view the guestbook

Run the following command to scale up the number of frontend Pods
<pre>kubectl scale deployment frontend --replicas=5</pre>

Query the list of Pods to verify the number of frontend Pods running
<pre>kubectl get pods</pre>

Run the following commands to delete all Pods, Deployments, and Services.
<pre>kubectl delete deployment -l app=redis
kubectl delete service -l app=redis
kubectl delete deployment -l app=guestbook
kubectl delete service -l app=guestbook</pre>
