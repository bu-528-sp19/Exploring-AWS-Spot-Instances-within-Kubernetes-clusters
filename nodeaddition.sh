pwd
cd /users/kevin
export KOPS_STATE_STORE=s3://ccproject-bucket2
kops get ig --name= nodes -o yaml > nodes.yml
python addnodes.py  
kops replace ig -f nodes.yml
kops update cluster --yes
