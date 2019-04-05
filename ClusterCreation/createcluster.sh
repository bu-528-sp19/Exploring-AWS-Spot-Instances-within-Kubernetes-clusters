aws s3 mb s3://ccproject-bucket2
export KOPS_STATE_STORE=s3://ccproject-bucket2
kops create cluster ccproject.cluster2.k8s.local --ssh-public-key  /Users/kevin/KevinSsh.pub --zones us-east-1a --yes --master-size=t2.micro --node-size=t2.micro --node-count=1
kops update cluster ccproject.cluster2.k8s.local --yes
