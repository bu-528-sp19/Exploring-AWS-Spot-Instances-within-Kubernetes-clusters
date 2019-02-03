# Exploring-AWS-Spot-Instances-within-Kubernetes-clusters
## **The goal of this project is to explore how we can use spot instances to reduce the cost of running kubernetes clusters.** 

**Background** 
- Amazon EC2 Spot Instances offer spare compute capacity available in the AWS cloud at steep discounts compared to On-Demand instances. Spot Instances enable you to optimize your costs on the AWS cloud and scale your application's throughput up to 10X for the same budget.
- Kubernetes is a portable, extensible open-source platform for managing containerized workloads and services, that facilitates both declarative configuration and automation. It has a large, rapidly growing ecosystem. Kubernetes services, support, and tools are widely available.

 ## **1. Vision and goals**
 - To provide a cost effective way of running a Kubernetes cluster using spot instances which are economical than the default On-Demand instances
 - To recommend the client an optimal percentage of total Kubernetes clusters to be used as spot instances to reduce the overall cost. 
 
 ## **2. Users**
 - This will be used by Cluster Admins and Cluster Developers using the Kubernetes Cluster for various application or projects.
 
 It does not target end users.
 
## **3. Scope and features**
- Able to handle various types of applications like JBoss, Spark, Hadoop, etc.
- Can handle sporadic or any distributed workloads.
- Primarily focuses on scaling of EC2 instances.

## **4. Solution**
- To create controller logic which will be continously looking for a better state.
  + The controller looks for available spot instances which are cheaper and are capable of processing the workload.
  + The controller always looks for a cheaper spot instance than the current spot instance while also considering time rquired for migration.
    

(A better state is when a controller finds a spot instance that is cheaper and provides the same computing power required and also time of migration is lesser than the processing time)
