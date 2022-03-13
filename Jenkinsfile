pipeline {
    agent {
        docker { image 'ubuntu:latest' }
    }

    environment {
    GITHUB_CREDENTIAL = credentials('github-credential')
    PROJECT_PATH = 'https://github.com/barezazad/gopher.git'
    IMAGE_NAME = 'gopher'
    }

    stages {
      stage('Git Clone') {
        steps {
            git '$PROJECT_PATH'
        }
      }

      stage('Build') {
        steps {
            sh 'docker build -t  ghcr.io/barezazad/$IMAGE_NAME:latest .'
        }
      }

      stage('Login') {
        steps {
            sh 'echo $GITHUB_CREDENTIAL_PSW | docker login ghcr.io -u $GITHUB_CREDENTIAL_USR --password-stdin'
        }
      }

      stage('Push') {
        steps {
            sh 'docker push ghcr.io/barezazad/$IMAGE_NAME:latest'
        }
      }
    }

    post {
      always {
          sh 'docker logout'
      }
    }
}
