# Exploring AWS Spot Instances within Kubernetes clusters

## **The goal of this project is to explore how we can use spot instances to reduce the cost of running kubernetes clusters.**

[![Watch the video](https://img.youtube.com/vi/X1dKqediJmE/maxresdefault.jpg)](https://youtu.be/X1dKqediJmE)


 Click here to watch our video where we discuss about the whole project

## ****Background****

  

- [Amazon EC2 Spot Instances]([https://aws.amazon.com/ec2/spot/](https://aws.amazon.com/ec2/spot/)  "Amazon EC2 Spot Instances") offer spare compute capacity available in the AWS cloud at steep discounts compared to On-Demand instances. Spot Instances enable you to optimize your costs on the AWS cloud and scale your application's throughput up to 10X for the same budget. The Spot price is determined by temporary trends in supply and demand and the amount of On-Demand capacity on a particular instance size, family, Availability Zone, and AWS Region

  

- [Kubernetes]([https://kubernetes.io/](https://kubernetes.io/)  "Kubernetes") is a portable, extensible open-source platform for managing containerized workloads and services, that facilitates both declarative configuration and automation. It has a large, rapidly growing ecosystem. Kubernetes services, support, and tools are widely available.

  

  

## ****The Motivation****

  

- Developers use On Demand EC2 instances to use with Kubernetes clusters as it is risk free and easier to manage compared to Spot Instances.

- But using such On Demand instances results in high costs of running an application.

- Using Amazon’s Elastic Kubernetes Service (EKS) one can manage Spot Instances, but Amazon charges $150 per month for running the control panel in addition to cost of running nodes which host the application.

- This additional cost can be avoided if you're administering the Kubernetes Cluster.

  

## ****The Solution****

- We create a controller which runs in an infinite while loop and sits on top of the Kubernetes.

- This controller performs these actions

	- Manages an application that runs on spot instances within a Kubernetes cluster in a cost effective way

	- Optimizes the number of spot instances dynamically

	- Maintains the Service Level Agreement (SLA) for the application

**
<p align="center">
<img src="https://lh6.googleusercontent.com/EsOKlESF_YkQ5rV7xiNm-COTyKvpnuzWd4sBZOuMzjB2YukxieLIbEEEI5h4qsDS5EQs62vqmtFXFoPTUtZLHWFLsmaQdlfsrDsnKkjiOnO-qJzcrdInZpKXRJqHOsYpMjkOQsQubwA" alt="alt text" width="400" height="400">

**

## ****Users****

  

- This will be used by cluster admins and developers using the Kubernetes Cluster and want to optimize their infrastructure cost for various applications or projects.

- It does not target non-expert users who are unfamiliar with Kubernetes or AWS.
## ****Tools Used****
#### [Kops]([https://github.com/kubernetes/kops](https://github.com/kubernetes/kops)) (Kubernetes Operations) 
-   Helps to create, destroy and upgrade clusters from command line
    
#### [Prometheus]([https://prometheus.io/](https://prometheus.io/)) (Monitoring tool)
-   Collects metrics from configured targets at given intervals
    
-   Displays the results
    

#### [Sock-Shop]([https://microservices-demo.github.io/docs/index.html](https://microservices-demo.github.io/docs/index.html)) (Demo Application)

-   The application is the user-facing part of an online shop that sells socks.
    
-   It is intended to aid the demonstration and testing of microservices and cloud native technologies.

## **Workflow**

**![](https://lh4.googleusercontent.com/IcNoiCC4YxYVly8smWrcYvX5d8TiUxoqIMBNt-AuovAL1X1Gnzv7A5bvlkdFL3RR3N1PLwYvrbK6d3bqPTtw3VRGe_OnzU8rEfphyXUIGRRXWc5sAWmcZf_NrnjYjWouZauB5r2mdY4)**

- Application is deployed on kubernetes cluster with on demand as and spot instances
- Prometheus pulls metrics from application
- Controller takes thresholds as inputs
- Controller then compares the threshold values with metrics from prometheus and takes decisions
- These decisions invoke actions that performed through kops on the cluster

## **SLA**
-   A service level agreement, or SLA, is a common term for formal service commitments that are made to customers by service providers.
    
-   The target SLA for our test application (Sock-shop) is set to have availability of 99% (“two nines”)
    
-   Availability is defined by these Service Level Indicators (SLIs)
    

	-   **Throughput**: More than 96%
    
	-   **Error Rate**: Less than 4%
    
	-   **Latency**: Less than 1.5 seconds

## **Budget**
-   Budget is the quota which we have in terms of SLA that allows our application to perform below set thresholds.
    
-   Budget is increased after a specific time interval.
    
-   Anytime our test application performs below a certain threshold, we reduce some amount from the budget
    
-   For the test application, we have defined the budget to be 15 mins per day.

## **Algorithm**
-   Our test application requires minimum 2 nodes to be on demand instances.
    
-   Based on network traffic (load) , our algorithm considers three cases
	1.   High Network Traffic ( more than 600 requests per second )
    
	2.   Medium Network Traffic ( between 100 and 600 requests per second )
    
	3.   Low Network Traffic ( less than 100 requests per second )


### Flowchart of Controller

**![](https://lh4.googleusercontent.com/34-OSES1K_7zYtevjMRKbStTlkkddgOknzLc2JELlZozVg3xcQcK_kW322JUXB2tHfanM3Tngg0bjcyxyBvgxH6AjLtSYDK69glnmdTVp4Cb27LLQbO5309DUBzQQ6T1098cPoJ5P0k)**


#### 1. High Network Traffic
-   If network_traffic > 600:
    

	-   If application performance is low for 10 mins:
    

		-   If budget available:
    

			-   Spin up a Spot instance scale up the application
    

		-   Else:
    

			-   Spin up an On demand and scale up the application
    

	-   Else
    

		-   Application is running smoothly, do nothing
#### **2. Moderate Network Traffic**
-   If 100 < network_traffic < 600:
	-   If application performance is high for 180 mins:
		-   If budget available:    
			-   Scale down
			-   Delete on demand
			-   Spin up a spot
    
			-   Scale up
    
		-   Else:
    
			-   Budget not available, do nothing

#### ****3. Low Network Traffic****
-   If network_traffic < 100:
    

	-   If application performance stable from 180 mins:
    

		-   If budget available :
    

			-   Scale down
    
			-   If On demand nodes > 2
    

				-   Delete an On-Demand Node
    

			-   Else
    

				-   Delete Spot instance

### ****Spot Termination****
- Before we spin up a Spot instance, we have to set a max bidding price.

- Spot Instances can be taken away at any time due to following reasons:

	-   **Price**: The Spot price is greater than your maximum bidding price.
    
	-   **Capacity**: If there are not enough unused EC2 instances to meet the demand for Spot Instances, Amazon EC2 interrupts Spot Instances. The order in which the instances are interrupted is determined by Amazon EC2.
-   Controller constantly checks for ‘to be terminated’ flag of Spot Instances.
    
	-   If to_be_terminated = True:
	    

		-   Scale down the application
	    
		-   Drain the spot instance
		    
		-   If budget available :
		    

			-   Spin up a spot instance
		    

		-   Else
		    

			-   Spin up an On demand Instance

## ****Input****
  
Config file that contains these parameters:
- SLA of the application provided by the client such as the availability of the application
- External IP of the prometheus deployment



## ****Milestones****

- Sprint 1:
	- Set up the prerequisites (AWS-CLI, kops, Kubectl)
	- Created a Kubernetes cluster using kops
	- Mock EC2 spot instances that simulated termination from AWS
- Sprint 2:
	- Tested multiple demo applications on minikube
	- Finalize Sock-shop as the demo application
- Sprint 3:
	- Used Jmeter to perform load testing on the application deployed to the cluster
	- Observed the impact on application due to scaling pods or nodes
	- Obtained kubernetes metrics using a Go based controller
- Sprint 4:
	- Monitored metrics using prometheus
	- Defined a basic SLA for the application
	- Created and deleted on demand and spot instances from the Go controller
- Sprint 5:
	- Defined the complete SLA for the project with the SLIs
	- Created the algorithm to make decisions

## **Cost Savings**
**![](https://lh4.googleusercontent.com/hj7fvXAeR_dV-26yvbPAV5xTXujmF4W0fMQ_W7NlB3AT2-q1TZHnb3efZFEPfAa3S-EsEaIHneYBsdR8LOvGYXJcz1VeYgrW2EQegshv0PL3Z6UsA8gsW08ePlCVAWfFg_MWlHhxJsA)**

- We started the application with 2 on-demand nodes of type ```t3.medium```
- We scaled up the cluster using only on-demand nodes in one case, and using only spot instances in the second case.
- We observed that the savings increased as we increased the number of nodes.
# Running the application

  

## Step 1: Prerequisite installations

1. [Install kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

2. [Install kops](https://github.com/kubernetes/kops)

3. [Install AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)


## Step 2: Configure AWS credentials

[Configure AWS](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html)

## Step 3: Setting up the cluster

#### Follow these steps to create a cluster 

1. Use the following command to create a S3 bucket for storing kops cluster data
```
aws s3 mb s3://<bucket-name>
export KOPS_STATE_STORE=s3://<bucket-name>
```
2. Use the following command to create a cluster
```
kops create cluster <cluster-name> --zones us-east-1a --yes --master-size=<machine-type> --node-size=<machine-type> --node-count=<count>
```
**Note: ```<cluster-name>.k8s.local``` creates a cluster without DNS.**  

Set the ```<machine-type>```  and ```<count>``` to the desired option.  

We set ```<machine-type> => t3.medium``` and ```<count> => 2``` for our test application.

3. Create a new instance group for Spot Instances
```
kops create ig <instance-group-name>
```
4. Edit the yaml file that opens in vim as follows
```
apiVersion: kops/v1alpha2

kind: InstanceGroup

metadata:

  creationTimestamp: 2019-05-01T20:54:19Z

  labels:

    kops.k8s.io/cluster: <cluster-name>

  name: <instance-group-name>

spec:

  image: kope.io/k8s-1.11-debian-stretch-amd64-hvm-ebs-2018-08-17

  machineType: <machine-type>

  maxPrice: <max-bidding-price>

  maxSize: 0
  minSize: 0
  nodeLabels:

    kops.k8s.io/instancegroup: <instance-group-name>

  role: Node

  subnets:

  - us-east-1a

```

Edit the ```minSize``` and ```maxSize``` both to 0, as we don't want the cluster to begin on Spot Instances.

This is done to make sure that the cluster is set-up without any risk of the nodes being deleted when deploying the application and the metrics.

Add a new parameter ```maxPrice:<max-bidding-price>``` above ```maxSize```.  

Set the ```<max-bidding-price>``` as the max price you want to bid for a spot instance.

Set the ```<machine-type>``` to the desired option.

We used ```t3.medium``` for our test application.

5. Use the kops update command to update the cluster
```
kops update cluster <cluster-name> --yes
```
6. Use kops validate to check if cluster is ready
```
kops validate cluster
``` 

## Step 4: Deploying required applications


#### Follow these steps to deploy the test application Sock-shop

1. Clone the Sock shop repository  

```git clone https://github.com/microservices-demo/microservices-demo.git```

2. Deploy the sock-shop along with prometheus application to your cluster using following commands  
```
cd microservices-demo/deploy/kubernetes/
kubectl create namespace sock-shop
kubectl apply -f complete-demo.yaml
kubectl create -f ./deploy/kubernetes/manifests-monitoring

```

3. Expose the deployments  
```
kubectl expose deployment front-end --type=LoadBalancer -n sock-shop --name=front-end-deployment

kubectl expose deployment prometheus-deployment --type=LoadBalancer -n monitoring --name=prometheus-deployment
```

4. Note down the exposed external-ip of the application
```
kubectl get service -n sock-shop
```
You will see an output similar to this example:
```
NAME           TYPE           CLUSTER-IP       EXTERNAL-IP                                                              PORT(S)          AGE
carts          ClusterIP      100.70.116.237   <none>                                                                   80/TCP           1d
carts-db       ClusterIP      100.70.224.159   <none>                                                                   27017/TCP        1d
catalogue      ClusterIP      100.65.171.188   <none>                                                                   80/TCP           1d
catalogue-db   ClusterIP      100.65.237.207   <none>                                                                   3306/TCP         1d
front-end      LoadBalancer   100.64.107.5     a471afe9b6ce911e9b71e02004e2872a-463552420.us-east-1.elb.amazonaws.com   8079:30416/TCP   14h
orders         ClusterIP      100.65.108.161   <none>                                                                   80/TCP           1d
orders-db      ClusterIP      100.67.156.105   <none>                                                                   27017/TCP        1d
payment        ClusterIP      100.69.170.8     <none>                                                                   80/TCP           1d
queue-master   ClusterIP      100.66.55.164    <none>                                                                   80/TCP           1d
rabbitmq       ClusterIP      100.70.49.108    <none>                                                                   5672/TCP         1d
shipping       ClusterIP      100.68.54.197    <none>                                                                   80/TCP           1d
user           ClusterIP      100.66.202.190   <none>                                                                   80/TCP           1d
user-db        ClusterIP      100.64.14.164    <none>                                                                   27017/TCP        1d
```
5. Note down the exposed external-ip of prometheus 
```
kubectl get service -n monitoring
```
You will see an output similar to this example:
```
NAME                    TYPE           CLUSTER-IP       EXTERNAL-IP                                                               PORT(S)          AGE
kube-state-metrics      ClusterIP      100.64.7.247     <none>                                                                    8080/TCP         1d
prometheus              NodePort       100.68.102.152   <none>                                                                    9090:31090/TCP   1d
prometheus-deployment   LoadBalancer   100.64.119.122   af1b1f0726c5011e9b71e02004e2872a-1377931422.us-east-1.elb.amazonaws.com   9090:30544/TCP   1d
```

View the website using ```http://<sock-shop-external-ip>:8079```
View Prometheus Graphs using ```http://<prometheus-external-ip>:9090```

## Step 5: Running the controller

1. Edit prometheus endpoint in the go program using the external ip of your prometheus deployment
2. Store the kops scripts inside bin folder of your Go installation.
3. Run the program

# Members

## Mentorss

[Daniel McPherson](https://github.com/danmcp)

[ravi gudimetla](https://github.com/ravisantoshgudimetla)

## Contributors

[Aditya Kadam](https://github.com/adi6496)

[kevin rodrigues](https://github.com/kevinrodrigues13)

[Nikhil Singh](https://github.com/Nikjz)

[Suryateja Gudiguntla](https://github.com/suryatej77)
