#!groovy

@Library("ucglobal") _

import java.time.LocalDateTime
import java.time.format.DateTimeFormatter
import java.util.*

optiva.addJobParam(stringParam(name: 'VERSION', defaultValue: '2.1.3',
                            description: 'Semantic version for the produced artifacts'))
                            
optiva.addJobParam(stringParam(name: 'GOLANG_VERSION', defaultValue: '1.19.3'))
optiva.addJobParam(stringParam(name: 'ALPINE_VERSION', defaultValue: '3.16'))

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

    echo "'version=$version'"
    echo "'buildNumber=${currentBuild.number}'"
    echo "'branch=${scm.branches[0].name}'"
    
    if (isOnPullRequest()) {
        def buildNumber = String.format( "%03d" , currentBuild.number );
        tag = "$version-SNAPSHOT-${LocalDateTime.now().format(DateTimeFormatter.ofPattern('yyyyMMdd'))}${buildNumber}"
    }else {
        
        if (env.GIT_BRANCH == 'develop') {
            tag = "$version-wip.${currentBuild.number}"
        } else  if (env.GIT_BRANCH.startsWith('release')) {
            tag = "$version-rc.${currentBuild.number}"
        } else  if (env.GIT_BRANCH.startsWith('hotfix')) {
            tag = "$version-hotfix.${currentBuild.number}"
        } else  if (env.GIT_BRANCH == 'main') {
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

                    sh "echo 'Version=$tag'"
                    docker.withRegistry('https://harbor.optiva.com', 'harbor-credentials') {
                        def krakenImage = docker.build("oce/krakend:${tag}", "--no-cache --build-arg GITHUB_TOKEN=${GITHUB_APP}:${GITHUB_ACCESS_TOKEN} --build-arg GOLANG_VERSION=${params.GOLANG_VERSION} --build-arg ALPINE_VERSION=${params.ALPINE_VERSION} .")
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
