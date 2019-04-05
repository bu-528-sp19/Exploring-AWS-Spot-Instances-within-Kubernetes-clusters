export KOPS_STATE_STORE=s3://ccproject-bucket2
kops delete cluster ccproject.cluster2.k8s.local --yes
