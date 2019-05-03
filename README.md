# Exploring AWS Spot Instances within Kubernetes clusters
## **The goal of this project is to explore how we can use spot instances to reduce the cost of running kubernetes clusters.** 

**Background** 
 - [Amazon EC2 Spot Instances](https://aws.amazon.com/ec2/spot/ "Amazon EC2 Spot Instances") offer spare compute capacity available in the AWS cloud at steep discounts compared to On-Demand instances. Spot Instances enable you to optimize your costs on the AWS cloud and scale your application's throughput up to 10X for the same budget. The Spot price is determined by temporary trends in supply and demand and the amount of On-Demand capacity on a particular instance size, family, Availability Zone, and AWS Region
 - [Kubernetes](https://kubernetes.io/ "Kubernetes") is a portable, extensible open-source platform for managing containerized workloads and services, that facilitates both declarative configuration and automation. It has a large, rapidly growing ecosystem. Kubernetes services, support, and tools are widely available.

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
![Proposed System](https://user-images.githubusercontent.com/20182350/52174269-67e5ca80-275f-11e9-95a4-4e592fee92cc.JPG)


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



To run our application you need to install these prerequisites first:  
 1.Install kubectl from https://kubernetes.io/docs/tasks/tools/install-kubectl/  
 2.Install kops from https://github.com/kubernetes/kops  
 3.Install AWS CLI from https://docs.aws.amazon.com/cli/latest/userguide/install-linux-al2017.html  
 
  

 
