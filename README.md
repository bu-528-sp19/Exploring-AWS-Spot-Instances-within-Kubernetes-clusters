# Exploring AWS Spot Instances within Kubernetes clusters

## **The goal of this project is to explore how we can use spot instances to reduce the cost of running kubernetes clusters.**

  

**Background**

- [Amazon EC2 Spot Instances]([https://aws.amazon.com/ec2/spot/](https://aws.amazon.com/ec2/spot/)  "Amazon EC2 Spot Instances") offer spare compute capacity available in the AWS cloud at steep discounts compared to On-Demand instances. Spot Instances enable you to optimize your costs on the AWS cloud and scale your application's throughput up to 10X for the same budget. The Spot price is determined by temporary trends in supply and demand and the amount of On-Demand capacity on a particular instance size, family, Availability Zone, and AWS Region

- [Kubernetes]([https://kubernetes.io/](https://kubernetes.io/)  "Kubernetes") is a portable, extensible open-source platform for managing containerized workloads and services, that facilitates both declarative configuration and automation. It has a large, rapidly growing ecosystem. Kubernetes services, support, and tools are widely available.

  

## **1. Vision and goals**

- To provide a cost-effective way of running a Kubernetes cluster using spot instances which are more economical than the default On-Demand instances

- To reduce the cost in such a way that the Service Level Agreement (SLA) is maintained

- To recommend the client an optimal percentage of total Kubernetes clusters to be used as spot instances to reduce the overall cost

- Kubernetes has built-in logic for self-healing which makes it easier and ideal target for the development of such a solution

## **2. Users**

- This will be used by cluster admins and developers using the Kubernetes Cluster and want to optimize their infrastructure cost for various applications or projects

- It does not target non-expert users who are unfamiliar with Kubernetes or AWS

## **3. Scope and features**

- Reads a user-defined config file that defines the parameters of the distribution of Spot Instances and runs the program

- Maintain the SLA for the application while using the Spot Instances

- Able to handle various types of applications like JBoss, Spark, Hadoop, etc.

- Can handle sporadic or any distributed workloads

- Primarily focuses on scaling of EC2 instances

- Can recover from node termination gracefully while maintaining the SLA for all applications

  

## **4. Input**

- Config file that contains parameters like SLA by the client including the up time of the application and the percentage of the cluster to be used as spot instances

  

## **5. Output**

- Stable running application

- Actual costs saved

- Recommended percentage of clusters to optimize cost savings and performance

  

## **6. Solution**

- Kubernetes distrubutes the workload onto the cluster nodes equally, with even load on every node

- Kubernetes is capable of rescheduling workloads when there are issues on infrastructure

- The proposed solution looks for ideal combination of On-Demand and Spot Instance nodes for the workload that the cluster is handling

- This combination is either pre-decided by the config file provided by the client or by the program we create

- The program thus proposed is a controller program that continuously looks for a better state that the application can exist in.

- This controller:

+ Runs constantly maintaining the SLA

+ Looks for available spot instances which are cheaper and can process the workload

+ Can dynamically handle any termination of node, either On-Demand or spot instance

+ Recovers from termination utilizing the cluster autoscaler component of Kubernetes to either evacuate and shift to a new On-demand node or decide to look for another available spot instance

+ May also look for a cheaper spot instance than the current spot instance and decides while considering time required for migration, in order to provide additional cost savings

+ Looks at signals(CPU Usage, Node CPU capacity, Memory usage, Node Memory capacity, Requests, Limits and pod health) from Spot Instances and workloads on Kubernetes to decide when to transfer a load from On-Demand node to Spot Instance node

+ Can handle any type of workload distribution for a variety of applications

- The program will never compromise the working of the application as defined by the SLA in order to reduce the costs. That is, it will not shift the load to a spot instance just for the sake of it.

- The program will also calculate the cost savings that the client makes, and can suggest an optimum percentage of workload to be run on Spot Instances to maximize savings while running no risk of stepping out of line of the SLA

  

### Proposed System

![Proposed System]([https://user-images.githubusercontent.com/20182350/52174269-67e5ca80-275f-11e9-95a4-4e592fee92cc.JPG](https://user-images.githubusercontent.com/20182350/52174269-67e5ca80-275f-11e9-95a4-4e592fee92cc.JPG))

  

  

## **7. Acceptance criteria**

- Controller can handle a single application with a constant load with 10% of the clusters as spot instances to reduce the effective cost

  

**Stretch goals**

- Dynamically scale the usage of spot instances

- Capable of handling various applications with sporadic workloads

  

## **8. Release planning**

- Release #1:

Create a controller (in GoLang) can call the AWS SDK

- Release #2:

Watch the AWS spot instance events to determine when a node is going away, figure out how to recover from such events with your Kubernetes cluster (Ex: evacuate the existing node and add a new node)

  

- Release #3: Use an example application (JBoss, Spark, etc.), simulate load, and trigger spot instance events to move the application and maintain SLA

  

  

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
```
kubectl get service -n sock-shop
```
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
```
kubectl get service -n monitoring
```
```
NAME                    TYPE           CLUSTER-IP       EXTERNAL-IP                                                               PORT(S)          AGE
kube-state-metrics      ClusterIP      100.64.7.247     <none>                                                                    8080/TCP         1d
prometheus              NodePort       100.68.102.152   <none>                                                                    9090:31090/TCP   1d
prometheus-deployment   LoadBalancer   100.64.119.122   af1b1f0726c5011e9b71e02004e2872a-1377931422.us-east-1.elb.amazonaws.com   9090:30544/TCP   1d
```
