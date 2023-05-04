pipeline {
    agent any
    stages{
        stage('Build docker image'){
            steps{
                script{
                    echo 'Starting build...'
                    git url: 'https://github.com/VedRatan/GO-CRUD.git', branch: 'main'
                    sh "docker build -t lusciousmaestro/api_image:${env.BUILD_NUMBER}.0 ."
                    echo 'Build completed'
                }
            }
        }
        stage('Push image to Hub'){
            steps{
                script{
                  withCredentials([string(credentialsId: 'dockerhub-pwd', variable: 'dockerhubpwd')]) {
                  sh 'docker login -u lusciousmaestro -p ${dockerhubpwd}'
                    }
                  sh "docker push lusciousmaestro/api_image:${env.BUILD_NUMBER}.0"
                }
            }
        }
        
        stage('Version') {
            steps {
                echo 'Starting update...'
                sh "sed -i 's|image: lusciousmaestro/api_image:.*|image: lusciousmaestro/api_image:${env.BUILD_NUMBER}.0|' docker-compose.yaml"
                sh "cat docker-compose.yaml"
                echo 'Update completed'   
            }
        }
        
        stage('Update the service'){
            steps{
                sh 'docker stop ci-cd_app_1'
                sh 'docker rm ci-cd_app_1'
                sh 'docker compose up -d --force-recreate --no-deps app'
                sh 'docker network disconnect ci-cd_default ci-cd_app_1'
                sh 'docker network connect flipr-assignment_default ci-cd_app_1'
            }
        }
    }
}