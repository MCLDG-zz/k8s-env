<h2>Kubernetes Environment Application</h2>

An application that can be run on Kubernetes on Amazon or Google Cloud that shows the environment variables within a container and the IP Address and POD details of the POD handling the request.

Is useful for showing load balancing and phased rollouts

<h3>Pre-requisites</h3>
You'll need a Kubernetes cluster running on AWS. Make sure your AWS config and credentials points to the region where you want to create the Kubernetes cluster. Instructions to create one are here: http://kubernetes.io/docs/getting-started-guides/aws/ or here: http://kubernetes.io/docs/admin/multiple-zones/ (for multi-zone clusters)

Access to the Docker registry on AWS, and the repositories storing the Docker images:
cpa/aml_k8s-env_front
185711092606.dkr.ecr.us-west-2.amazonaws.com/cpa/aml_k8s-env_front
cpa/aml_k8s-env_back
185711092606.dkr.ecr.us-west-2.amazonaws.com/cpa/aml_k8s-env_back

These are both under the CPA AML account:
Account: 1857-1109-2606

To access the registry you need to login to the AWS ECR service (do this from the same server as you would run 'kubectl'):
aws ecr get-login --region us-west-2
then paste the resulting 'docker login' statement 

<h3>Starting PODS and Services</h3>
Start and stop the pods and services using start-all.sh and stop-all.sh. update-all.sh can be used if you want to apply changes to the running cluster (such as a rolling update where you roll out another version of a container)

<h3>Building and pushing Docker images</h3>
After making changes to the application in the 'container' directory, run the associated 'docker-build.sh' script. Before tagging and pushing, edit the 'docker-tag-push.sh' and make sure it points to the appropriate registry and repository, then change the version number. Then execute the script. This will push the Docker image to the registry

