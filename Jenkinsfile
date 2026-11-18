pipeline {
    agent any

    environment {
        DOCKER_IMAGE = "roofiahmad/store-ns-pipeline"
        IMAGE_TAG    = "${env.BUILD_NUMBER}" 
        APP_DIR      = "app-repo"
        MANIFEST_DIR = "k8s-manifests"
    }

    stages {
        stage('1. Clone Repositories') {
            steps {
                dir("${APP_DIR}") {
                    git branch: 'main',
                        credentialsId: 'github-roofiahmad',
                        url: 'git@github.com:Roofiahmad/store-notification-service.git'
                }

                dir("${MANIFEST_DIR}") {
                    git branch: 'main',
                        credentialsId: 'github-roofiahmad',
                        url: 'git@github.com:Roofiahmad/store-manifest.git',
                        poll: false,
                        changelog: false
                }
            }
        }

        stage('2. Trivy Dependency Scan (3rd-Party Packages)') {
            steps {
                dir("${APP_DIR}") {
                    sh """
                        trivy fs \
                        --scanners vuln \
                        --severity HIGH,CRITICAL \
                        --skip-version-check \
                        --format template \
                        --template '{{ range . }}{{ \$target := .Target }}{{ range .Vulnerabilities }}[{{ .Severity }}] {{ \$target }} | {{ .PkgName }} ({{ .InstalledVersion }} -> {{ .FixedVersion }}) | {{ .VulnerabilityID }}\n{{ end }}{{ end }}' \
                        --exit-code 1 \
                        --no-progress \
                            .
                        """
                }
            }
        }

        stage('3. Build & Package Spring Boot') {
            steps {
                dir("${APP_DIR}") {
                    sh "docker build -t ${DOCKER_IMAGE}:${IMAGE_TAG} ."
                }
            }
        }

        stage('4. Trivy Security Scan') {
            steps {
                    sh """
                        trivy image \
                          --severity HIGH,CRITICAL \
                          --scanners vuln \
                          --format template \
                          --template '{{ range . }}{{ \$target := .Target }}{{ range .Vulnerabilities }}[{{ .Severity }}] {{ \$target }} | {{ .PkgName }} ({{ .InstalledVersion }} -> {{ .FixedVersion }}) | {{ .VulnerabilityID }}\n{{ end }}{{ end }}' \
                          --exit-code 1 \
                          --no-progress \
                          ${DOCKER_IMAGE}:${IMAGE_TAG}
                    """
                }
        }

        stage('5. Push to Docker Hub') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-credentials',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh 'echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin'
                    sh "docker push ${DOCKER_IMAGE}:${IMAGE_TAG}"
                }
            }
        }

        stage('6. Update Kubernetes Deployment YAML') {
            steps {
                dir("${MANIFEST_DIR}") {
                    sh """
                        sed -i "s|image:.*|image: ${DOCKER_IMAGE}:${IMAGE_TAG}|g" notification-service.yaml
                        git diff notification-service.yaml
                    """
                    
                            sshagent(credentials: ['github-roofiahmad']) {
                        sh """
                            git config user.email "jenkins@yourdomain.com"
                            git config user.name "Jenkins CI"
                            git add notification-service.yaml
                            git commit -m "chore: update image tag to ${DOCKER_IMAGE}:${IMAGE_TAG}"
                            git push origin HEAD:main
                        """
                    }
                }
            }
        }
    }

    post {
        always {
            sh "docker rmi ${DOCKER_IMAGE}:${IMAGE_TAG} || true"
            sh "docker image prune -a --force || true"
            sh 'docker logout || true'
        }
    }
}