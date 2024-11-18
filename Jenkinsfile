pipeline {
    environment {
        DOCKER_REGISTRY = "your-registry"
        BUILD_TAG = "v${BUILD_NUMBER}"
        DOCKER_CREDENTIALS_ID = 'docker-registry-credentials'
        KUBE_CONFIG_ID = 'minikube-config'  // ID of the kubeconfig credential in Jenkins
    }
    
    agent any
    
    stages {
        stage('Build and Push Images') {
            steps {
                script {
                    docker.withRegistry("https://${DOCKER_REGISTRY}", DOCKER_CREDENTIALS_ID) {
                        parallel(
                            "Product Service": {
                                dir('product-service') {
                                    def productImage = docker.build("${DOCKER_REGISTRY}/product-service:${BUILD_TAG}")
                                    productImage.push()
                                }
                            },
                            "Inventory Service": {
                                dir('inventory-service') {
                                    def inventoryImage = docker.build("${DOCKER_REGISTRY}/inventory-service:${BUILD_TAG}")
                                    inventoryImage.push()
                                }
                            },
                            "Order Service": {
                                dir('order-service') {
                                    def orderImage = docker.build("${DOCKER_REGISTRY}/order-service:${BUILD_TAG}")
                                    orderImage.push()
                                }
                            },
                            "API Gateway": {
                                dir('api-gateway') {
                                    def gatewayImage = docker.build("${DOCKER_REGISTRY}/api-gateway:${BUILD_TAG}")
                                    gatewayImage.push()
                                }
                            }
                        )
                    }
                }
            }
        }
        
        stage('Deploy to Minikube') {
            steps {
                withKubeConfig([credentialsId: KUBE_CONFIG_ID]) {
                    sh '''
                        # Create namespace
                        kubectl apply -f k8s/namespaces/namespaces.yaml

                        # Create ConfigMaps
                        kubectl apply -f k8s/configmaps/config.yaml

                        # Generate deployment files
                        mkdir -p generated-k8s
                        for file in k8s/deployments/*.yaml; do
                            envsubst < $file > "generated-k8s/$(basename $file)"
                        done
                        
                        # Apply Kubernetes configurations
                        kubectl apply -f k8s/services/services.yaml
                        kubectl apply -f generated-k8s/
                        kubectl apply -f k8s/ingress/ingress.yaml
                        
                        # Wait for deployments
                        kubectl -n microservices rollout status deployment/product-service
                        kubectl -n microservices rollout status deployment/inventory-service
                        kubectl -n microservices rollout status deployment/api-gateway
                    '''
                }
            }
        }
        
        stage('Verify Deployment') {
            steps {
                withKubeConfig([credentialsId: KUBE_CONFIG_ID]) {
                    sh '''
                        echo "Service Status:"
                        kubectl get svc
                        
                        echo "\nPod Status:"
                        kubectl get pods
                        
                        echo "\nDeployment Status:"
                        kubectl get deployments
                    '''
                }
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
        success {
            echo "Deployment to Minikube completed successfully!"
        }
        failure {
            echo "Deployment to Minikube failed!"
        }
    }
}