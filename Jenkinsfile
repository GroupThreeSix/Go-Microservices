pipeline {
    environment {
        DOCKER_REGISTRY = "tuilakhanh"
        BUILD_TAG = "v${BUILD_NUMBER}-${GIT_COMMIT[0..7]}"
        DOCKER_CREDENTIALS_ID = 'dockercerd'
        KUBE_CONFIG_ID = 'k8s-config'
        KUBE_CLUSTER_NAME = 'minikube'
        KUBE_CONTEXT_NAME = 'minikube'
        KUBE_SERVER_URL = 'https://192.168.39.206:8443'
        REPORT_DIR = 'reports'
        SONAR_PROJECT_BASE_DIR = '.'
        SONAR_SCANNER_OPTS = '-Xmx2048m'
    }
    
    agent any
    
    stages {
        stage('Generate Protobuf') {
            steps {
                sh '''
                    chmod +x scripts/gen-proto.sh
                    ./scripts/gen-proto.sh
                '''            }
        }
        
        stage('Lint') {
            steps {
                script {
                    sh "mkdir -p ${REPORT_DIR}"
                    
                    def services = ['product-service', 'inventory-service', 'order-service', 'api-gateway']
                    services.each { service ->
                        dir(service) {
                            sh """
                                golangci-lint run --out-format checkstyle ./... > ../${REPORT_DIR}/${service}-lint.xml || true
                            """
                        }
                    }
                }
                
                recordIssues(
                    tools: [
                        checkStyle(pattern: "${REPORT_DIR}/*-lint.xml", reportEncoding: 'UTF-8')
                    ],
                    qualityGates: [[threshold: 100, type: 'TOTAL', unstable: true]],
                    healthy: 50,
                    unhealthy: 100,
                    minimumSeverity: 'WARNING'
                )
            }
        }
        
        stage('SonarQube Analysis') {
            steps {
                script {
                    def scannerHome = tool 'SonarScanner'
                    def services = ['product-service', 'inventory-service', 'order-service', 'api-gateway']
                    
                    // Run all analyses first
                    withSonarQubeEnv('sq1') {
                        services.each { service ->
                            dir(service) {
                                sh """
                                    ${scannerHome}/bin/sonar-scanner \
                                    -Dsonar.projectKey=${service} \
                                    -Dsonar.projectName=${service} \
                                    -Dsonar.sources=. \
                                    -Dsonar.exclusions=**/*_test.go,**/vendor/**,**/proto/** \
                                    -Dsonar.go.coverage.reportPaths=coverage.out \
                                    -Dsonar.go.tests.reportPaths=test-report.json \
                                    -Dsonar.qualitygate.wait=false
                                """
                            }
                        }
                    }
                    
                    // Then wait for all quality gates
                    timeout(time: 5, unit: 'MINUTES') {
                        services.each { service ->
                            def qg = waitForQualityGate projectKey: service
                            if (qg.status != 'OK') {
                                error "Quality gate failed for ${service}: ${qg.status}"
                            }
                            echo "Quality gate passed for ${service}"
                        }
                    }
                }
            }
        }
        
        stage('Build and Deploy Services') {
            parallel {
                stage('Product Service') {
                    when {
                        anyOf {
                            changeset "product-service/**/*"
                            changeset "shared/**/*"
                            changeset "proto/**/*"
                        }
                    }
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
                    when {
                        anyOf {
                            changeset "inventory-service/**/*"
                            changeset "shared/**/*"
                            changeset "proto/**/*"
                        }
                    }
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
                    when {
                        anyOf {
                            changeset "order-service/**/*"
                            changeset "shared/**/*"
                        }
                    }
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
                    when {
                        anyOf {
                            changeset "api-gateway/**/*"
                            changeset "shared/**/*"
                        }
                    }
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
        
        stage('Verify Deployments') {
            steps {
                verifyDeployments()
            }
        }
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
    withKubeConfig(clusterName: KUBE_CLUSTER_NAME, contextName: KUBE_CONTEXT_NAME, credentialsId: KUBE_CONFIG_ID, serverUrl: KUBE_SERVER_URL) {
        sh """
            # First time setup - apply base resources if they don't exist
            kubectl apply -f k8s/base/namespace.yml || true
            kubectl apply -f k8s/base/config.yaml || true
            kubectl apply -f k8s/base/services.yaml || true
            
            # Update only the specific service
            cd k8s/overlay/services/${serviceName}
            
            # Update the image tag in kustomization.yaml
            kustomize edit set image ${serviceName}=${DOCKER_REGISTRY}/${serviceName}:${BUILD_TAG}
            
            # Apply changes only for this service
            kustomize build . | kubectl apply -f -
            
            # Wait for deployment
            kubectl -n microservices rollout status deployment/${serviceName}
        """
    }
}

def verifyDeployments() {
    withKubeConfig(clusterName: KUBE_CLUSTER_NAME, contextName: KUBE_CONTEXT_NAME, credentialsId: KUBE_CONFIG_ID, serverUrl: KUBE_SERVER_URL) {
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