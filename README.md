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
