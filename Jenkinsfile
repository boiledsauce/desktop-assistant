pipeline {
    agent any

    environment {
        DOCKER_IMAGE = "localhost/desktop-assistant-image:${env.BUILD_NUMBER}"
        PROJECT_DIR = "/home/wille/projects/desktop-assistant"
        K8S_DIR = "/home/wille/projects/desktop-assistant/k8s"
        KIND_CLUSTER_NAME = "kind-cluster"
    }

    stages {
        stage('Create Kubernetes Cluster') {
            steps {
                script {
                    def clusterExists = sh(script: 'kind get clusters | grep -w ${KIND_CLUSTER_NAME} || true', returnStdout: true).trim()
                    if (!clusterExists) {
                        echo "Creating Kind cluster ${KIND_CLUSTER_NAME}"
                        sh 'kind create cluster --name ${KIND_CLUSTER_NAME} --config ${K8S_DIR}/kind-config.yaml'
                    } else {
                        echo "Kind cluster ${KIND_CLUSTER_NAME} already exists"
                    }
                }
            }
        }

        stage('Prepare Environment') {
            steps {
                dir("${PROJECT_DIR}") {
                    echo "Current directory: ${pwd()}"
                    sh 'ls -la'
                }
            }
        }

        stage('Build Docker image') {
            steps {
                dir("${PROJECT_DIR}") {
                    sh 'docker build -t ${DOCKER_IMAGE} .'
                    sh 'kind load docker-image ${DOCKER_IMAGE} --name ${KIND_CLUSTER_NAME}'
                }
            }
        }

        stage('Update Kubernetes Deployment') {
            steps {
                script {
                    // Load YAML as a text instead of structured YAML to avoid formatting issues
                    def deploymentYAML = readFile("${K8S_DIR}/desktop-assistant.yaml").replaceAll(
                        "(image:\\s+).*",
                        "\$1${DOCKER_IMAGE}"
                    )
                    writeFile file: "${K8S_DIR}/desktop-assistant.yaml", text: deploymentYAML

                    sh 'kubectl apply -f ${K8S_DIR}/pv-and-pvc.yaml'
                    sh 'kubectl apply -f ${K8S_DIR}/desktop-assistant-config.yaml'
                    sh 'kubectl apply -f ${K8S_DIR}/desktop-assistant.yaml'
                }
            }
        }
    }
}
