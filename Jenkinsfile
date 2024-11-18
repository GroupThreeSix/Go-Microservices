pipeline {
    environment {
        DOCKER_REGISTRY = "tuilakhanh"
        BUILD_TAG = "v${BUILD_NUMBER}-${GIT_COMMIT[0..7]}"
        DOCKER_CREDENTIALS_ID = 'dockercerd'
        // KUBE_CONFIG_ID = 'minikube-config'
    }
    
    agent any
    
    stages {
        stage('Generate Protobuf') {
            steps {
                sh '''
                    chmod +x scripts/gen-proto.sh
                    ./scripts/gen-proto.sh
                '''
            }
        }
        
        stage('Build and Deploy Services') {
            parallel {
                stage('Product Service') {
                    // when {
                    //     anyOf {
                    //         changeset "product-service/**/*"
                    //         changeset "shared/**/*"
                    //         changeset "proto/**/*"
                    //     }
                    // }
                    stages {
                        stage('Build Product Service') {
                            steps {
                                buildAndPushImage('product-service')
                            }
                        }
                        stage('Deploy Product Service') {
                            steps {
                                deployService('product-service')
                            }
                        }
                    }
                }
                
                stage('Inventory Service') {
                    // when {
                    //     anyOf {
                    //         changeset "inventory-service/**/*"
                    //         changeset "shared/**/*"
                    //         changeset "proto/**/*"
                    //     }
                    // }
                    stages {
                        stage('Build Inventory Service') {
                            steps {
                                buildAndPushImage('inventory-service')
                            }
                        }
                        stage('Deploy Inventory Service') {
                            steps {
                                deployService('inventory-service')
                            }
                        }
                    }
                }
                
                stage('Order Service') {
                    // when {
                    //     anyOf {
                    //         changeset "order-service/**/*"
                    //         changeset "shared/**/*"
                    //     }
                    // }
                    stages {
                        stage('Build Order Service') {
                            steps {
                                buildAndPushImage('order-service')
                            }
                        }
                        stage('Deploy Order Service') {
                            steps {
                                deployService('order-service')
                            }
                        }
                    }
                }
                
                stage('API Gateway') {
                    // when {
                    //     anyOf {
                    //         changeset "api-gateway/**/*"
                    //         changeset "shared/**/*"
                    //     }
                    // }
                    stages {
                        stage('Build API Gateway') {
                            steps {
                                buildAndPushImage('api-gateway')
                            }
                        }
                        stage('Deploy API Gateway') {
                            steps {
                                deployService('api-gateway')
                            }
                        }
                    }
                }
            }
        }
        
        // stage('Verify Deployments') {
        //     steps {
        //         verifyDeployments()
        //     }
        // }
    }
    
    post {
        always {
            cleanWs()
        }
        success {
            echo "Pipeline completed successfully!"
        }
        failure {
            echo "Pipeline failed!"
        }
    }
}

// Helper functions
def buildAndPushImage(String serviceName) {
    script {
        docker.withRegistry("", DOCKER_CREDENTIALS_ID) {
            dir(serviceName) {
                def serviceImage = docker.build("${DOCKER_REGISTRY}/${serviceName}:${BUILD_TAG}")
                serviceImage.push()
                serviceImage.push('latest')
            }
        }
    }
}

def deployService(String serviceName) {
    // withKubeConfig([credentialsId: KUBE_CONFIG_ID]) {
    //     sh """
    //         # Create namespace if not exists
    //         kubectl apply -f k8s/namespace.yml
    //         kubectl apply -f k8s/config.yaml
            
    //         # Generate deployment files
    //         mkdir -p generated-k8s
    //         envsubst < k8s/${serviceName}-services.yaml > generated-k8s/${serviceName}.yaml
            
    //         # Apply Kubernetes configurations
    //         kubectl apply -f k8s/services.yaml
    //         kubectl apply -f generated-k8s/${serviceName}.yaml
            
    //         # Wait for deployment
    //         kubectl -n microservices rollout status deployment/${serviceName}
    //     """
    // }
}

def verifyDeployments() {
    withKubeConfig([credentialsId: KUBE_CONFIG_ID]) {
        sh '''
            echo "Services Status:"
            kubectl get svc -n microservices
            
            echo "\nPods Status:"
            kubectl get pods -n microservices
            
            echo "\nDeployments Status:"
            kubectl get deployments -n microservices
        '''
    }
}