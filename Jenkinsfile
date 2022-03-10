pipeline {
  agent any

  environment {
    GITHUB_TOKEN = credentials('github-credential')
    IMAGE_NAME = 'barezazad/gopher'
    IMAGE_VERSION = '0.0.1'
  }

  stages {
    stage('build image') {
      steps {
        sh 'docker build . -t $IMAGE_NAME:$IMAGE_VERSION -f Dockerfile'
      }
    }

    stage('login to GHCR') {
      steps {
        sh 'echo "${{ secrets.TOKEN }}" |  docker login ghcr.io -u barezazad --password-stdin'
      }
    }

    stage('tag image') {
      steps {
        sh 'docker tag $IMAGE_NAME:$IMAGE_VERSION ghcr.io/$IMAGE_NAME:$IMAGE_VERSION'
      }
    }

    stage('push image') {
      steps {
        sh 'docker push ghcr.io/$IMAGE_NAME:$IMAGE_VERSION'
      }
    }
  }

  post {
    always {
      sh 'docker logout'
    }
  }
}
