pwd
cd /users/kevin/go/bin
export KOPS_STATE_STORE=s3://ccproject-bucket11
kops get ig --name= nodes -o yaml > nodes.yaml
python SubNodes.py  
kops replace ig -f nodes.yaml
kops update cluster --yes
