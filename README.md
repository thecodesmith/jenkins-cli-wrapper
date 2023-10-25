# Jenkins CLI Wrapper

_A simple and fast CLI to develop Jenkins builds on the command line_

## Features

- [ ] Repository-aware: project metadata enables seamless pipeline operations for multibranch pipelines
- [ ] Branch-aware pipeline replay and pipeline logs
- [ ] Replay uses Jenkinsfile in current directory
- [x] Jenkinsfile linting: defaults to Jenkinsfile in current directory
- [ ] Open pipeline in browser

## Usage

    jenkinsw [command] [options]

    Commands:
      context  Configure multiple Jenkins servers and switch between them
      help     Display help info for wrapper commands
      init     Download jenkins-cli.jar from Jenkins server and initialize API token
      lint     Lint a Declarative Jenkinsfile
      logs     Display the logs for a multibranch pipeline job
      replay   Replay a multibranch pipeline job
      version  Display version info for the Jenkins server, CLI and wrapper

    jenkinsw lint  # runs declarative-linter on Jenkinsfile in current directory
    jenkinsw lint -j foo/Jenkinsfile  # runs declarative-linter on Jenkinsfile specified by path

    jenkinsw context list
    jenkinsw context add
    > Jenkins host URL:
    > Jenkins user:
    > Jenkins API key:

## Development

Install the Cobra CLI:

    make setup

Generate a new CLI command:

    cobra-cli add <subcommand> -a '<your name>' -l mit --viper

Generate a new CLI subcommand:

    cobra-cli add <subcommand> -p <parent command> -a '<your name>' -l mit --viper
