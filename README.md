# A Simple Microservices Application writen in Go and gRPC to demonstrate Jenkins CI/CD Pipeline

## Services
- Product Service
- Inventory Service
- Order Service
- API Gateway

## CI/CD Pipeline
- Jenkins Pipeline to build, test, and deploy the application to Docker containers
- SonarQube to perform static code analysis
- Docker to containerize the application
- Kubernetes to deploy the application

## Needed Environment Tools
- Docker
- Docker Compose
- Jenkins
- SonarQube
- Kubernetes (Minikube)

## Jenkins Plugins
- SonarQube Scanner
- Docker Pipeline
- Kubernetes CLI

## Setup Jenkins
- Install Jenkins Plugins
- Setup SonarQube server and add to Jenkins, setup webhook
- Setup Kubernetes Cluster via Minikube and add token to Jenkins Credentials
- Setup Docker
- Create Jenkins Pipeline
- Link to git repository
- Setup webhook to trigger Jenkins Pipeline
