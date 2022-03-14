# Jenkins CLI Wrapper

_A simple and fast CLI to develop Jenkins builds on the command line_

## Features

[ ] Repository-aware: project metadata enables seamless pipeline operations for multibranch pipelines
[ ] Branch-aware pipeline replay and pipeline logs
[ ] Replay uses Jenkinsfile in current directory
[ ] Jenkinsfile linting: defaults to Jenkinsfile in current directory
[ ] Open pipeline in browser

## Usage

    jenkinsw [command] [options]

    Commands:
      init     Download jenkins-cli.jar from Jenkins server and initialize API token
      lint     Lint a Declarative Jenkinsfile
      replay   Replay a multibranch pipeline job
      logs     Display the logs for a multibranch pipeline job
      version  Display version info for the Jenkins server, CLI and wrapper
      help     Display help info for wrapper commands
