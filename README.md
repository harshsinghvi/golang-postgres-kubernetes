# golang-postgres-kubernetes

## practice

- postgressql - indexing, explain querry
- GoLang - APIs, concurrency
- Kubernetes EKS - Autoscaling, load balencing,
- API LoadTesting - Apache Benchmark

## resources

- <https://go.dev/doc/tutorial/web-service-gin>
- <https://dev.to/ramu_mangalarapu/building-rest-apis-in-golang-go-gin-with-persistence-database-postgres-4616>
- <https://www.coding-bootcamps.com/blog/build-containerized-applications-with-golang-on-kubernetes.html>
- <https://docs.aws.amazon.com/eks/latest/userguide/horizontal-pod-autoscaler.html>
- <https://stackgres.io/features/>
- <https://dev.to/asizikov/using-github-container-registry-with-kubernetes-38fb> ghcr.io kubernetes

- <https://aws.amazon.com/blogs/containers/using-alb-ingress-controller-with-amazon-eks-on-fargate/> fargarte exose services
- <https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html> alb imp
- <https://docs.aws.amazon.com/eks/latest/userguide/alb-ingress.html> eks ingress imp
- <https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/> HPA

- <https://artifacthub.io/packages/helm/metrics-server/metrics-server> install Matrics server
- <https://docs.aws.amazon.com/eks/latest/userguide/metrics-server.html> matrics server
- <https://levelup.gitconnected.com/how-to-deploy-a-multi-container-two-tier-go-application-in-eks-fargate-6266494f5bcf> go and postgres eks

- golang postgres api <https://medium.com/@cavdy/creating-restful-api-using-golang-and-postgres-part-2-542aac86e2bd> <https://medium.com/@cavdy/creating-restful-api-using-golang-and-postgres-part-1-58fe83c6f1ee>

## TODOS

- Golang API
- Deploy go API to Kubernetes

- test autoscaling using Apache benchmark

- setup CI/CD pipeline
- Connect external postgress to it
- deploy postgress to Kubernetes
- autoscale postgress deployment

## K8S procedure

1. eksctl faragete cluster `eksctl create cluster --name cluster --region ap-south-1 --fargate`
1. cluster ALB ingress <https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html>
1. setup matrics server (for HPA) from YML
1. sertup efs (elastic file storage) <https://github.com/kubernetes-sigs/aws-efs-csi-driver/blob/master/docs/efs-create-filesystem.md> get file_system_id and replace volumeHandle: fs-1234567899 in database.yml
1. ghcr secrets for image <https://dev.to/asizikov/using-github-container-registry-with-kubernetes-38fb> replace required fields in secrets.yml

1. deploy services (yml files) yml files includes HPA

## commands

```bash
kubectl rollout restart deployment/name # to update image
kubectl get ingress # ingress exposed url
kubectl port-forward statefulset.apps/postgres 5432:5432
kubectl exec --stdin --tty pod/postgres-0 -- /bin/bash
aws eks update-kubeconfig --region ap-south-1 --name cluster
```

## GHCR image build and push

`https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry`

```bash
docker buildx build --platform=linux/amd64 -t golang-postgres-kubernetes .
docker tag golang-postgres-kubernetes ghcr.io/harshsinghvi/golang-postgres-kubernetes:latest
docker push ghcr.io/harshsinghvi/golang-postgres-kubernetes:latest
```

## ELB and ingress SETUP

```bash
ACCOUNT_ID= # aws sts get-caller-identity
AWS_EKS_CLUSTER_NAME=cluster
AWS_EKS_CLUSTER_REGION=ap-south-1

AWS_EKS_CLUSTER_VPC_ID=$(aws eks describe-cluster \
    --name $AWS_EKS_CLUSTER_NAME \
    --query "cluster.resourcesVpcConfig.vpcId" \
    --output text)

# AWS_EKS_CLUSTER_VPC_ID= # console>cloudformations

curl -O https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.5.4/docs/install/iam_policy.json

aws iam create-policy \
    --policy-name AWSLoadBalancerControllerIAMPolicy \
    --policy-document file://iam_policy.json

eksctl utils associate-iam-oidc-provider --region=ap-south-1 --cluster=cluster --approve

eksctl create iamserviceaccount \
  --cluster=cluster \
  --namespace=kube-system \
  --name=aws-load-balancer-controller \
  --role-name AmazonEKSLoadBalancerControllerRole \
  --attach-policy-arn=arn:aws:iam::194505915562:policy/AWSLoadBalancerControllerIAMPolicy \
  --approve

helm repo add eks https://aws.github.io/eks-charts

aws sts get-caller-identity

helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system \
  --set clusterName=cluster \
  --set serviceAccount.create=false \
  --set serviceAccount.name=aws-load-balancer-controller \
  --set region=ap-south-1 \
  --set vpcId=vpc-07ae5f71518dd2545
  
kubectl get deployment -n kube-system aws-load-balancer-controller 

                    # during upgrade 
                    kubectl apply -k "github.com/aws/eks-charts/stable/aws-load-balancer-controller/crds?ref=master"

                    helm upgrade aws-load-balancer-controller eks/aws-load-balancer-controller \
                    -n kube-system \
                    --set clusterName=cluster \
                    --set serviceAccount.create=false \
                    --set serviceAccount.name=aws-load-balancer-controller \
                    --set region=ap-south-1 \
                    --set vpcId=vpc-07ae5f71518dd2545
```

## EFS Setup

```bash
AWS_EKS_CLUSTER_NAME=cluster
AWS_EKS_CLUSTER_REGION=ap-south-1

vpc_id=$(aws eks describe-cluster \
    --name $AWS_EKS_CLUSTER_NAME \
    --query "cluster.resourcesVpcConfig.vpcId" \
    --output text)

cidr_range=$(aws ec2 describe-vpcs \
    --vpc-ids $vpc_id \
    --query "Vpcs[].CidrBlock" \
    --output text \
    --region $AWS_EKS_CLUSTER_REGION)


security_group_id=$(aws ec2 create-security-group \
    --group-name MyEfsSecurityGroup \
    --description "My EFS security group" \
    --vpc-id $vpc_id \
    --output text)

aws ec2 authorize-security-group-ingress \
    --group-id $security_group_id \
    --protocol tcp \
    --port 2049 \
    --cidr $cidr_range

file_system_id=$(aws efs create-file-system \
    --region ap-south-1 \
    --performance-mode generalPurpose \
    --query 'FileSystemId' \
    --output text)


aws ec2 describe-subnets \
    --filters "Name=vpc-id,Values=$vpc_id" \
    --query 'Subnets[*].{SubnetId: SubnetId,AvailabilityZone: AvailabilityZone,CidrBlock: CidrBlock}' \
    --output table

# run for each subnet
aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-09555c7ce2147f642  \
    --security-groups $security_group_id
aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-019b6e706b2823a7b  \
    --security-groups $security_group_id
aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-0324d7a94eb3afd09  \
    --security-groups $security_group_id
aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-04d07f3812cf78123  \
    --security-groups $security_group_id
aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-0ee5c658df8ef377c  \
    --security-groups $security_group_id
aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-0360ff2918bf5fceb  \
    --security-groups $security_group_id
```

## AUTOSCALE LOGS HPA

`kubectl get hpa --watch`

```text
go-todo-api-hpa   Deployment/go-todo-api   15%/30%   1         10        2          19m
go-todo-api-hpa   Deployment/go-todo-api   14%/30%   1         10        1          19m
go-todo-api-hpa   Deployment/go-todo-api   15%/30%   1         10        1          19m
go-todo-api-hpa   Deployment/go-todo-api   14%/30%   1         10        1          20m
go-todo-api-hpa   Deployment/go-todo-api   15%/30%   1         10        1          20m
go-todo-api-hpa   Deployment/go-todo-api   14%/30%   1         10        1          20m
go-todo-api-hpa   Deployment/go-todo-api   22%/30%   1         10        1          21m
go-todo-api-hpa   Deployment/go-todo-api   26%/30%   1         10        1          22m
go-todo-api-hpa   Deployment/go-todo-api   26%/30%   1         10        1          22m
go-todo-api-hpa   Deployment/go-todo-api   27%/30%   1         10        1          22m
go-todo-api-hpa   Deployment/go-todo-api   25%/30%   1         10        1          22m
go-todo-api-hpa   Deployment/go-todo-api   26%/30%   1         10        1          23m
go-todo-api-hpa   Deployment/go-todo-api   26%/30%   1         10        1          23m
go-todo-api-hpa   Deployment/go-todo-api   27%/30%   1         10        1          23m
go-todo-api-hpa   Deployment/go-todo-api   26%/30%   1         10        1          23m
go-todo-api-hpa   Deployment/go-todo-api   36%/30%   1         10        1          24m
go-todo-api-hpa   Deployment/go-todo-api   37%/30%   1         10        2          24m
go-todo-api-hpa   Deployment/go-todo-api   38%/30%   1         10        2          24m
go-todo-api-hpa   Deployment/go-todo-api   38%/30%   1         10        2          24m
go-todo-api-hpa   Deployment/go-todo-api   37%/30%   1         10        2          25m
go-todo-api-hpa   Deployment/go-todo-api   29%/30%   1         10        2          25m
go-todo-api-hpa   Deployment/go-todo-api   22%/30%   1         10        2          25m
go-todo-api-hpa   Deployment/go-todo-api   21%/30%   1         10        2          26m
go-todo-api-hpa   Deployment/go-todo-api   21%/30%   1         10        2          26m
go-todo-api-hpa   Deployment/go-todo-api   21%/30%   1         10        2          27m
go-todo-api-hpa   Deployment/go-todo-api   21%/30%   1         10        2          28m
go-todo-api-hpa   Deployment/go-todo-api   21%/30%   1         10        2          28m
go-todo-api-hpa   Deployment/go-todo-api   21%/30%   1         10        2          28m
go-todo-api-hpa   Deployment/go-todo-api   21%/30%   1         10        2          29m
go-todo-api-hpa   Deployment/go-todo-api   67%/30%   1         10        2          29m
go-todo-api-hpa   Deployment/go-todo-api   73%/30%   1         10        4          29m
go-todo-api-hpa   Deployment/go-todo-api   75%/30%   1         10        5          29m
go-todo-api-hpa   Deployment/go-todo-api   74%/30%   1         10        5          30m
go-todo-api-hpa   Deployment/go-todo-api   73%/30%   1         10        5          30m
go-todo-api-hpa   Deployment/go-todo-api   72%/30%   1         10        5          30m
go-todo-api-hpa   Deployment/go-todo-api   23%/30%   1         10        5          30m
go-todo-api-hpa   Deployment/go-todo-api   8%/30%    1         10        5          31m
go-todo-api-hpa   Deployment/go-todo-api   6%/30%    1         10        5          31m
go-todo-api-hpa   Deployment/go-todo-api   3%/30%    1         10        5          31m
go-todo-api-hpa   Deployment/go-todo-api   3%/30%    1         10        5          31m
go-todo-api-hpa   Deployment/go-todo-api   3%/30%    1         10        5          32m
go-todo-api-hpa   Deployment/go-todo-api   28%/30%   1         10        5          32m
go-todo-api-hpa   Deployment/go-todo-api   40%/30%   1         10        5          32m
go-todo-api-hpa   Deployment/go-todo-api   40%/30%   1         10        7          32m
go-todo-api-hpa   Deployment/go-todo-api   39%/30%   1         10        7          33m
go-todo-api-hpa   Deployment/go-todo-api   43%/30%   1         10        7          33m
go-todo-api-hpa   Deployment/go-todo-api   36%/30%   1         10        7          33m
go-todo-api-hpa   Deployment/go-todo-api   22%/30%   1         10        7          33m
go-todo-api-hpa   Deployment/go-todo-api   15%/30%   1         10        7          34m
go-todo-api-hpa   Deployment/go-todo-api   13%/30%   1         10        7          34m
go-todo-api-hpa   Deployment/go-todo-api   7%/30%    1         10        7          34m
go-todo-api-hpa   Deployment/go-todo-api   6%/30%    1         10        7          34m
go-todo-api-hpa   Deployment/go-todo-api   7%/30%    1         10        7          35m
go-todo-api-hpa   Deployment/go-todo-api   19%/30%   1         10        7          35m
go-todo-api-hpa   Deployment/go-todo-api   32%/30%   1         10        7          35m
go-todo-api-hpa   Deployment/go-todo-api   33%/30%   1         10        7          36m
go-todo-api-hpa   Deployment/go-todo-api   42%/30%   1         10        7          36m
go-todo-api-hpa   Deployment/go-todo-api   38%/30%   1         10        7          36m
go-todo-api-hpa   Deployment/go-todo-api   24%/30%   1         10        7          36m
go-todo-api-hpa   Deployment/go-todo-api   17%/30%   1         10        7          37m
go-todo-api-hpa   Deployment/go-todo-api   16%/30%   1         10        7          37m
go-todo-api-hpa   Deployment/go-todo-api   9%/30%    1         10        7          37m
go-todo-api-hpa   Deployment/go-todo-api   9%/30%    1         10        7          37m
go-todo-api-hpa   Deployment/go-todo-api   10%/30%   1         10        7          38m
go-todo-api-hpa   Deployment/go-todo-api   19%/30%   1         10        7          38m
go-todo-api-hpa   Deployment/go-todo-api   24%/30%   1         10        7          38m
go-todo-api-hpa   Deployment/go-todo-api   31%/30%   1         10        7          39m
go-todo-api-hpa   Deployment/go-todo-api   44%/30%   1         10        7          39m
go-todo-api-hpa   Deployment/go-todo-api   45%/30%   1         10        7          39m
go-todo-api-hpa   Deployment/go-todo-api   34%/30%   1         10        7          39m
go-todo-api-hpa   Deployment/go-todo-api   21%/30%   1         10        7          40m
go-todo-api-hpa   Deployment/go-todo-api   22%/30%   1         10        7          40m
go-todo-api-hpa   Deployment/go-todo-api   13%/30%   1         10        7          40m
go-todo-api-hpa   Deployment/go-todo-api   30%/30%   1         10        7          41m
go-todo-api-hpa   Deployment/go-todo-api   40%/30%   1         10        7          41m
go-todo-api-hpa   Deployment/go-todo-api   38%/30%   1         10        7          41m
go-todo-api-hpa   Deployment/go-todo-api   36%/30%   1         10        7          42m
go-todo-api-hpa   Deployment/go-todo-api   32%/30%   1         10        7          42m
go-todo-api-hpa   Deployment/go-todo-api   31%/30%   1         10        7          42m
go-todo-api-hpa   Deployment/go-todo-api   29%/30%   1         10        7          43m
go-todo-api-hpa   Deployment/go-todo-api   25%/30%   1         10        7          43m
go-todo-api-hpa   Deployment/go-todo-api   20%/30%   1         10        7          44m
go-todo-api-hpa   Deployment/go-todo-api   19%/30%   1         10        7          44m
go-todo-api-hpa   Deployment/go-todo-api   19%/30%   1         10        7          44m
go-todo-api-hpa   Deployment/go-todo-api   21%/30%   1         10        7          44m
go-todo-api-hpa   Deployment/go-todo-api   19%/30%   1         10        7          45m
go-todo-api-hpa   Deployment/go-todo-api   16%/30%   1         10        7          45m
go-todo-api-hpa   Deployment/go-todo-api   16%/30%   1         10        7          45m
go-todo-api-hpa   Deployment/go-todo-api   17%/30%   1         10        7          45m
go-todo-api-hpa   Deployment/go-todo-api   14%/30%   1         10        7          46m
go-todo-api-hpa   Deployment/go-todo-api   15%/30%   1         10        7          46m
go-todo-api-hpa   Deployment/go-todo-api   18%/30%   1         10        7          46m
go-todo-api-hpa   Deployment/go-todo-api   14%/30%   1         10        7          46m
go-todo-api-hpa   Deployment/go-todo-api   14%/30%   1         10        7          47m
go-todo-api-hpa   Deployment/go-todo-api   18%/30%   1         10        7          47m
go-todo-api-hpa   Deployment/go-todo-api   14%/30%   1         10        7          47m
go-todo-api-hpa   Deployment/go-todo-api   13%/30%   1         10        7          47m
go-todo-api-hpa   Deployment/go-todo-api   15%/30%   1         10        7          48m
go-todo-api-hpa   Deployment/go-todo-api   15%/30%   1         10        7          48m
go-todo-api-hpa   Deployment/go-todo-api   13%/30%   1         10        7          48m
go-todo-api-hpa   Deployment/go-todo-api   13%/30%   1         10        5          48m
go-todo-api-hpa   Deployment/go-todo-api   15%/30%   1         10        4          49m
go-todo-api-hpa   Deployment/go-todo-api   16%/30%   1         10        4          49m
go-todo-api-hpa   Deployment/go-todo-api   15%/30%   1         10        4          49m
go-todo-api-hpa   Deployment/go-todo-api   18%/30%   1         10        4          49m
go-todo-api-hpa   Deployment/go-todo-api   18%/30%   1         10        4          50m
go-todo-api-hpa   Deployment/go-todo-api   18%/30%   1         10        3          50m
go-todo-api-hpa   Deployment/go-todo-api   20%/30%   1         10        3          50m
go-todo-api-hpa   Deployment/go-todo-api   24%/30%   1         10        3          50m
go-todo-api-hpa   Deployment/go-todo-api   24%/30%   1         10        3          51m
go-todo-api-hpa   Deployment/go-todo-api   22%/30%   1         10        3          51m
go-todo-api-hpa   Deployment/go-todo-api   23%/30%   1         10        3          52m
go-todo-api-hpa   Deployment/go-todo-api   22%/30%   1         10        3          52m
go-todo-api-hpa   Deployment/go-todo-api   23%/30%   1         10        3          52m
go-todo-api-hpa   Deployment/go-todo-api   20%/30%   1         10        3          52m
go-todo-api-hpa   Deployment/go-todo-api   8%/30%    1         10        3          53m
go-todo-api-hpa   Deployment/go-todo-api   0%/30%    1         10        3          53m
```
