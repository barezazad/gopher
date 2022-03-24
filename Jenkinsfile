pipeline {
    agent any

    environment {
    GITHUB_CREDENTIAL = credentials('github-credential')
    PROJECT = 'https://github.com/barezazad/gopher.git'
    CONTAINER_REGISTRY = 'barezazad'
    IMAGE_NAME = 'gopher'
    }

    stages {

      stage('Login') {
        steps {
            sh 'echo $GITHUB_CREDENTIAL_PSW | docker login ghcr.io -u $GITHUB_CREDENTIAL_USR --password-stdin'
        }
      }

      stage('Clone & Build') {
        steps {
            sh 'docker build $PROJECT -t  ghcr.io/$CONTAINER_REGISTRY/$IMAGE_NAME:latest '
        }
      }

      stage('Push') {
        steps {
            sh 'docker push ghcr.io/$CONTAINER_REGISTRY/$IMAGE_NAME:latest'
        }
      }

      stage('Remove the built image') {
        steps {
            sh 'docker image rm ghcr.io/$CONTAINER_REGISTRY/$IMAGE_NAME:latest'
            sh 'docker image rm -f $(docker images -f dangling=true -q)'
        }
      }

    }

    post {
      always {
          sh 'docker logout'
      }
    }
} 