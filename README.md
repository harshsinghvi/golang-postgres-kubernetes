# golang-postgres-kubernetes

## practice

- postgressql - indexing, explain querry
- GoLang - APIs, concurrency
- Kubernetes EKS - Autoscaling, load balencing,
- API LoadTesting - Apache Benchmark

## resources

- Linkedin Learning (kubernetes, golang courses)
- udemy hnsr course db
- https://go.dev/doc/tutorial/web-service-gin
- https://dev.to/ramu_mangalarapu/building-rest-apis-in-golang-go-gin-with-persistence-database-postgres-4616
- https://www.coding-bootcamps.com/blog/build-containerized-applications-with-golang-on-kubernetes.html
- https://docs.aws.amazon.com/eks/latest/userguide/horizontal-pod-autoscaler.html
- https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/
- https://dev.to/asizikov/using-github-container-registry-with-kubernetes-38fb ghcr.io kubernetes
- https://aws.amazon.com/blogs/containers/using-alb-ingress-controller-with-amazon-eks-on-fargate/ fargarte exose services

- https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html alb imp
- https://docs.aws.amazon.com/eks/latest/userguide/alb-ingress.html eks ingress imp
- https://artifacthub.io/packages/helm/metrics-server/metrics-server install Matrics server
- https://docs.aws.amazon.com/eks/latest/userguide/metrics-server.html matrics server 
- https://levelup.gitconnected.com/how-to-deploy-a-multi-container-two-tier-go-application-in-eks-fargate-6266494f5bcf go and postgres eks

## TODOS

- Golang API
- Deploy go API to Kubernetes

- test autoscaling using Apache benchmark

- setup CI/CD pipeline
- Connect external postgress to it
- deploy postgress to Kubernetes
- autoscale postgress deployment

## AWS Resources created (tags: pingsafe-test)
- EKS IAM Role
- EKS Cluster
- ECS Cluster
- vpc subnets

docker login -u AWS -p $(aws ecr get-login-password --region ap-south-1) 194505915562.dkr.ecr.ap-south-1.amazonaws.com


## Build docker image for eks faragete 
// login ghcr docker
export CR_PAT=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
echo $CR_PAT | docker login ghcr.io -u harshsinghvi --password-stdin

docker buildx build --platform=linux/amd64 -t ghcr.io/harshsinghvi/golang-postgres-kubernetes:latest .

docker push ghcr.io/harshsinghvi/golang-postgres-kubernetes:latest

```

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

aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-0f8e0a7152ce63763 \
    --security-groups $security_group_id

aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-0e2824fc49bdcd202 \
    --security-groups $security_group_id

aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-028474cfc7ca5c2c5 \
    --security-groups $security_group_id

aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-060db0728f89d0203 \
    --security-groups $security_group_id

aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-00782a5c917b7ae74 \
    --security-groups $security_group_id

aws efs create-mount-target \
    --file-system-id $file_system_id \
    --subnet-id subnet-0a7fb187cb42744b1 \
    --security-groups $security_group_id
