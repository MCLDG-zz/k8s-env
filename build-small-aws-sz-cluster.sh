export KUBERNETES_PROVIDER=aws
export KUBE_AWS_ZONE=us-west-2a
export NUM_NODES=2
export MASTER_SIZE=m3.medium
export NODE_SIZE=m3.medium
export AWS_S3_REGION=us-west-2
export AWS_S3_BUCKET=cx-aml-kubernetes-artifacts
export INSTANCE_PREFIX=k8s-cx

curl -sS https://get.k8s.io | bash

