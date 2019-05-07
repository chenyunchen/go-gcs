def getFailureCount() {
    def curBuild = currentBuild.getPreviousBuild()
    def count = 0
    while (curBuild != null && curBuild.result != 'SUCCESS'){
        count++
        curBuild = curBuild.getPreviousBuild()
    }
    return count
}
def sendSlackFailureMessage(String stage) {

    def message =
        "<http://192.168.11.49:8080/job/FileManager/|FileManager> » stage `${stage}` failed " +
        "<${env.JOB_URL}|${env.GIT_BRANCH}> » " +
        "<${env.RUN_DISPLAY_URL}|#${env.BUILD_NUMBER}>."

    def failCount = getFailureCount()
    if (failCount > 0) {
        message += " This job has failed `${failCount+1}` times."
    }

    slackSend channel: '#backend-ci', color: 'danger', message: message

    if (env.gitlabActionType == "PUSH") {
        message += " <!here>"
        slackSend channel: '#jello_backend', color: 'danger', message: message
    }
}

def sendSlackFixedMessage() {
    def message =
        "<http://192.168.11.49:8080/job/FileManager/|FileManager> » " +
        "<${env.JOB_URL}|${env.GIT_BRANCH}> » " +
        "<${env.RUN_DISPLAY_URL}|#${env.BUILD_NUMBER}> is fixed " +
        "after ${getFailureCount()} failed builds."

    slackSend channel: '#backend-ci', color: 'good', message: message
}

pipeline {
    agent any
    options {
         gitLabConnection('Gitlab')
    }
    post {
        fixed {
            sendSlackFixedMessage()
        }
        always {
            script {
                if (env.gitlabActionType == 'MERGE') {
                    junit '*.xml'
                }
            }
        }
    }
    stages {
        stage('Prepare') {
            when {
                expression { -> env.gitlabActionType == 'MERGE'}
            }
            steps {
                gitlabCommitStatus('Prepare') {
                    script {
                        sh "go get -u github.com/tebeka/go2xunit"
                    }
                }
            }
            post {
                failure {
                    sendSlackFailureMessage 'Prepare'
                }
            }
        }
        stage('Testing') {
            parallel {
                stage('Unit Test') {
                    when {
                        expression { -> env.gitlabActionType == 'MERGE'}
                    }
                    steps {
                        gitlabCommitStatus('Unit Test') {
                            script {
                                sh 'make src.test | /Users/jello/go/bin/go2xunit -output junit.xml'
                            }
                        }
                    }
                    post {
                        failure {
                            sendSlackFailureMessage 'Unit Test'
                        }
                    }
                }
                stage('Build Test') {
                    when {
                        expression { -> env.gitlabActionType == 'MERGE'}
                    }
                    steps {
                        gitlabCommitStatus('Build Test') {
                            script {
                                sh 'eval $(minikube docker-env) && make dockerfiles.build-local'
                            }
                        }
                    }
                    post {
                        failure {
                            sendSlackFailureMessage 'Build Test'
                        }
                    }
                }
            }
        }
        stage('Deploy Local') {
            when {
                expression { -> env.gitlabActionType == 'MERGE'}
            }
            steps {
                gitlabCommitStatus('Deploy Local') {
                    script {
                        sh 'make apps.install-local'
                    }
                }
            }
            post {
                failure {
                    sendSlackFailureMessage 'Deploy Local'
                }
            }
        }
        stage('Build GCR Image') {
            when {
                expression { -> env.gitlabActionType == 'PUSH'}
            }
            steps {
                gitlabCommitStatus('Build GCR Image') {
                    script {
                        sh "export PATH=$PATH:/Users/jello/google-cloud-sdk/bin && gcloud config set account ${env.GCR_ACCOUNT}"
                        sh "export PATH=$PATH:/Users/jello/google-cloud-sdk/bin && make dockerfiles.build-${env.ENV}"
                    }
                }
            }
            post {
                failure {
                    sendSlackFailureMessage 'Build GCR Image'
                }
            }
        }
        stage('Deploy GKE') {
            when {
                expression { -> env.gitlabActionType == 'PUSH'}
            }
            steps {
                gitlabCommitStatus('Deploy GKE') {
                    script {
                        sh "export PATH=$PATH:/Users/jello/google-cloud-sdk/bin && gcloud config set account ${env.GKE_ACCOUNT}"
                        sh "export PATH=$PATH:/Users/jello/google-cloud-sdk/bin && make apps.install-${env.ENV}"
                    }
                }
            }
            post {
                failure {
                    sendSlackFailureMessage 'Deploy'
                }
            }
        }
    }
}
