```
###############################################################
Welcome to the Redwoods fuzzing Suite
###############################################################
The Idea behind this project is to have a base for automated fuzzing of your code packages.
Fuzzing is a great way to find critical errors in your code that otherwise would have remained hidden.

```

# How to use it

```

How to use:
In order to use redwoods you first need to create a config. You can do that either through:

```
$ redwoods -c
```

which will take you through a wizzard or use the

```
$ redwoods defaultconfig
```

to create a redwoods-cfg.json in you working directory.

NOTE: This program will create subfolders like "fuzz", "analyze" and "workspace" so its best to run it
in a directory of your choice apart from a project.

You find addionnal help with:

```
$ redwoods [command] -h
```

Usage:
  redwoods [flags]
  redwoods [command]

Available Commands:
  analyze        Analyse the AST of the source for enhanced information
  analyzefuzz    Find the fuzzable packages in the given repository reduced to fuzzing packages
  buildfuzz      Build the fuzz-build archives
  completion     generate the autocompletion script for the specified shell
  defaultconfig  Creates the default config so it can be modified
  dependencies   Download all dependencies needed for this project
  gomodanalyzer  Analyse the go.mod file for all packages
  help           Help about any command
  init           Creates Workfolders and downloads the project form git
  outputanalysis Analyses the output after fuzzing
  runfuzz        Build the fuzz-build archives
  workflow       Run predefined workflows automatically

Flags:
  -h, --help          help for redwoods
  -c, --set config    set new configuration
  -s, --show config   show current configuration
  -v, --version       show current version of CLI
```

## Features

- automatic detection of fuzzing tests
- complete analysis of all functions and packages available in a project
- analysis of the go.mod file
- building and running fuzztests with go-fuzz


## Requirements:

- At this point the program requires the fuzzing tests to be in files that end with _fuzz.go. This makes it super easy for the parser to pick up. It also does check for the //+build gofuzz but only of the file is suffixed with fuzz.go

For instance fuzzing the XOR Package of hashicorps vault could create the following file to test their XOR Implementation:

xor_fuzz.go
```
//+build gofuzz

package xor

func Fuzz_xor(data []byte) int {
	_, _ = XORBase64(string(data), string(data))
	return 1
}
```

## Using Containers to do you work

This Project also makes it possible to do all the fuzzing inside of a container requiring you to install no dependencies at all.

By Using the command:

```
$ make docker-run
```

which will build, tag and run a version of the docker container on your host which will execute all steps of the build process inside of a container regardless of your host system.
If you are using private packages, all you have to do is mount your git config inside of the container as a shared volume or you can do the same thing using kubernetes.

Attention:
- the docker-run command will mount the redwoods-cfg from the current directory into the container. the (pwd)/fuzz dir and local redwoods-cfg.json will also be mounted.


## Running in interactive mode

The Program provides an interactive dockermode inwhich the program will be built copied and installed on a golang container and be ready to use.

```
$ make docker-run-interactive
```

at which point you can use the redwoods suite as it were installed locally.

TIP:

You can create and easy way to fuzz local packages without having to install the dependencies on your local machine by using docker.
In the case you would like to keep the results for further analysis, mount a volume to your /fuzz/ and /analysis/ as well as your redwoods-cfg.json

## Installing and uninstalling

In order to install redwoods into your GOBIN directory simply run "make install" in the root of the project.

```
$ make install
```

In oder to unsintall redwoods from yoru GOBIN directoy simply run "make uninstall" in the root of the package

```
make uninstall
```

## The makefile

Redwoods is accompanied by a makefile in order to build docker images, build the source itself or install it into your gobin directoy. 
For more help use:

```
$ make help
```

```
###############################################################
Welcome to the Redwoods fuzzing Suite
###############################################################

Usage:
  make <target>
  help             Display this help.

Development
  fmt              Run go fmt against code.
  vet              Run go vet against code.
  test             Run tests.
  run              Run.

Core
  clean            remove previous binaries
  build            build a version of the app, pass Buildversion, Comit and projectname as build arguments
  install          install redwoods into gobin
  uninstall        uninstall redwoods from gobin

Container Deployment
  docker-build     Build the docker image and tag it with the current version and :latest
  docker-run       Build the docker image and tag it and run it in docker
  docker-run-interactive  run an interactive container
  docker-push     push your image to the docker hub  
  ```

## Versioning and build information

Version and build information is being passed through the makefile. Please make changes to your own makefile to change the container repository, your github path as well as the intended version number your would like to build. Redwoods uses and internal versioning package that can easily be extended. It can be found under /pkgs/


## Aknowledgements:

This project is possible thanks to the <https://github.com/dvyukov/go-fuzz> Project

## Intersting stuff
- go-fuzz requires gcc
- go fuzz logs to stderr instead of stdout
- go ast could be used for fuzz test definition


# Requirements

These are the assumptions this system is so far build upon:

- The program assumes to be running in a working golang environment
- required packages for the system are: sudo make bash apk add curl git nano gcc if you for instance start with the golang/alpine image
- Read and Write acces to the folder and subfolders
- An Example of this can be found in the dockerfiles

# Infrastructure assumptions

- The Program runs on a dedicated ec2 instance or your local host
- The system has a working git setup
- if VPNs are neccesairy they need to be set up

## Infra requirments

- On a Container crash the process wil  have to be restarted: this is important later for k8s resource limits on pods

# Kubernetes plans

NOTE: Future talk

Given the fact all above is implemented the container will be able to be mounted in k8s too.
This will however require acces to an image repository to push the images and an implementation for shared results.
State needs to be communicated through the uploaders (s3, http api), but this make it very attractive to implement this proccess in to the CI Pipeline.
A first draft for this idea can be found in the docs section.