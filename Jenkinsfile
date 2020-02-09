pipeline {

    agent any

    environment {
        registry = "cipheredbytes/cloud-native-restful-go"
        registryCredential = "dockerhub"
    }

    stages {
        stage('Build docker image') {
            steps {
                sh "docker build --no-cache -f docker/app/Dockerfile -t $registry:${env.BUILD_NUMBER} ."
            }
        }
        stage('Publish docker image') {
            steps {
                withDockerRegistry([ credentialsId: registryCredential, url: "" ]) {
                    sh "docker push $registry:${env.BUILD_NUMBER}"
                }
            }
        }
        stage('Clean docker image') {
            steps {
                sh "docker rmi $registry:${env.BUILD_NUMBER}"
            }
        }
    }
}
