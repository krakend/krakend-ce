#!groovy

@Library("ucglobal") _

import java.time.LocalDateTime
import java.time.format.DateTimeFormatter
import java.util.*

optiva.addJobParam(stringParam(name: 'VERSION', defaultValue: '1.0.0',
                            description: 'Semantic version for the produced artifacts'))

// Closure makeKrakenD() {
//     return {
//         withCredentials([usernamePassword(credentialsId: 'github-app-main',
//                                   usernameVariable: 'GITHUB_APP',
//                                   passwordVariable: 'GITHUB_ACCESS_TOKEN')]) {
//             sh("make docker GITHUB_TOKEN=${GITHUB_APP}:${GITHUB_ACCESS_TOKEN} VERSION=${VERSION}")
//         }
//     }
// }

// def krakendStage() {
//     optiva.addJobStage(optiva.PRIORITY_AUTO, 'Build KrakenD', makeKrakenD())
// }


def tag

boolean isOnPullRequest() {
    return (env.CHANGE_ID != null
            && env.CHANGE_BRANCH != null
            && env.CHANGE_TARGET != null)
}


// TODO : THIS WILL NEED TO BE UPDATED APPROPRIATELY
def setTag() {
    def version = params.VERSION
    if (isOnPullRequest()) {
        def buildNumber = String.format( "*%03d*" , currentBuild.number );
        tag = "$version-SNAPSHOT-${LocalDateTime.now().format(DateTimeFormatter.ofPattern('yyyyMMdd'))}${buildNumber}"
    }else {
        if (env.GIT_BRANCH != 'main') {
            tag = "$version.${currentBuild.number}"
        }else {
            tag = "$version"
        }
    }
}


Closure pushImagesToHarborClosure() {
    return {
        
        withCredentials([usernamePassword(credentialsId: 'github-app-main',
                                  usernameVariable: 'GITHUB_APP',
                                  passwordVariable: 'GITHUB_ACCESS_TOKEN')]) {

                withCredentials([usernamePassword(credentialsId: 'harbor-credentials', passwordVariable: 'DOCKER_PASSWORD',
                    usernameVariable: 'DOCKER_USERNAME')]) {
                    env.JAVA_HOME = env.JDK_111
                    sh("export DOCKER_USERNAME=${DOCKER_USERNAME}")
                    sh("export DOCKER_PASSWORD=${DOCKER_PASSWORD}")

                    sh('export GOLANG_VERSION=1.19.3')
                    sh('export ALPINE_VERSION=1.19.3')


                    docker.withRegistry('https://harbor.optiva.com', 'harbor-credentials') {
                        def krakenImage = docker.build("oce/krakend:${tag}", "--no-cache --build-arg GITHUB_TOKEN=${GITHUB_APP}:${GITHUB_ACCESS_TOKEN} --build-arg GOLANG_VERSION=${GOLANG_VERSION} --build-arg ALPINE_VERSION=${ALPINE_VERSION} .")
                        krakenImage.push()
                        krakenImage.push('latest')

                    }
                }
        }
    }
}

def krakendDockerStage() {
    optiva.addJobStage(optiva.PRIORITY_AUTO, 'Build and Push KrakenD', pushImagesToHarborClosure())
}




def stages() {
    setTag()
    krakendDockerStage()
}

stages()

optiva.addJobFinalStage(optiva.PRIORITY_AUTO, 'post') {

    if (currentBuild.currentResult == "SUCCESS") {
        def dateString = LocalDateTime.now().format(DateTimeFormatter.ofPattern("yyyyMMdd"))
        currentBuild.displayName = "#${currentBuild.number}-${dateString}"
    }
    if (params.KEEP_SLAVE) {
        sleep(15 * 60)
    }

}

optiva.customPipeline()
