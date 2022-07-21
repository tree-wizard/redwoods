# Intro

Thanks for taking the time to review this. We can call this project ‘redwoods’, until I come up with a name.

Let me describe how I use go-fuzz at work, in my local terminal:
I clone the teams project repo and install go-fuzz
$ git clone git@github.net:team/project.git
$ go get -u github.com/dvyukov/go-fuzz/go-fuzz@latest github.com/dvyukov/go-fuzz/go-fuzz-build@latest
Then cd into the package directory
$ cd team/project/pkg/crypto/
I then write new tests for each function; each file is tagged with a ‘// +build gofuzz’ header. Currently I save them in a fuzz/ dir; ie team/project/pkg/crypto/fuzz. But I am still working on the best organization method so these may change.
After writing a test I usually build the fuzzers locally and run for a few minutes. 
$ go-fuzz-build
The go-fuzz-build command creates a .zip binary and I specify the function
$ go-fuzz -bin <package_fuzz.zip> -func <fuzzAES>

It is out of scope for now but I would like to add support for f.fuzz fuzzers which use a built-in framework and are saved in the same repo as the code. 
$ cd team/project/pkg/crypto/
$ gotip test -fuzz FuzzComputeHeader

Notes: 
1) go-fuzz tests are in <filename>_fuzz.go and f.fuzz() are <filename>_fuzz_test.go
2) Currently my go-fuzz and f.fuzz() test functions often share the same name
3) We have to start and stop the tests, inserting different function names manually.



## Component 1:
The goal of this component is to help me understand and scope the code for fuzzing. I envisioned this all as part of one tool, but we can actually make this a standalone utility.

I would like to download the repo on my machine and run something like:
$ git clone git@github.net:team/project.git
$ redwood -package /repo/pkg/crypto

Or better yet a smarter version that could also call git directly and would discover the packages automatically:
$ redwood -repo /repo/
$ redwood -git git@git@github.net:team/project.git

It would print the name and CLOC of every package, the count and names of all ‘normal’ functions, and the count and name of all unit tests and existing fuzzers: Something like:

repo/pkg/crypto
3450 lines of go
23 functions:
Parsejson 
Convertnum
Ect…
12 Unit tests:
TestParseJson
ect...
2 go-fuzz functions:
FuzzParseJson

repo/pkg/types/
5650 lines of go
13 functions:
Parse... 
10 Unit tests:
TestParseInt
1 go-fuzz functions:
FuzzParseJson

If you can add the cyclomatic complexity in a column next to the function list that would be phenomenal. This utility would allow me to understand the code and see what unit tests I can use to create fuzz functions. We would probably use the same fuzz_test_find() function to gather the list of fuzzers to use in the next steps.


## Component 2:
I broadly class this component as building and running of fuzzers.

I write the fuzz tests locally so I’m living in the package dir and always run the test for a short while to ensure it is working.
$ go-fuzz-build
$ go-fuzz -bin parse_fuzz.zip -func FuzzJsonParseKey

This is where I really need help. Right now I SCP that parse_fuzz.zip binary to an ec2 server and ssh in to run the‘$ go-fuzz -bin parse_fuzz.zip -func FuzzJsonParseKey’ command.

How can I deploy the code and tests to the fuzzing infra? First thoughts are I create a new branch on the repo, push the code changes with the tests, and clone the repo branch when building the machines. 

I do think I will be modifying the tests locally or on the server in case they are not working as efficiently or have some issues, like low coverage. Also, in order to pull code from our repos we need to be on the VPN and have the git SSH key so it’s kind of messy, at least for me but I’m open to all ideas. 

My goals would be to have a command like
redwood deploy -repo /path/ 
Which builds the binary and creates the k8’s to start running.

Another issue is the past tests and corpus data.There is no state that keeps track of how long the fuzzer has been running but there is a working directory that keeps track of crashes, suppressions, and is building a corpus of test inputs. By default it’s kept in the same dir as where you run but we can specify with -workdir=fuzz-workdir

As long as you don’t delete that working directory, you can stop and re-start go-fuzz and it’ll pick up pretty much where it left off. 

My projects have about 20-30 tests right now, with plans for many more as I build them out. My simple/dumb solution was to alternate between the tests for ~50 hours on a single machine. The ultimate goal is to get at least 1000 hours of cpu fuzzing time per function.

The tool has the ability to utilize several machines with a coordinator process. Documented here, maybe that can be used.

I would love to spin up multiple machines splitting the load between something like 5-10 servers. 
 


## Component 3:
I classify this as the analysis stage.

After running for several hours I’d like to have stats of the progress so far. The stdout output looks like this:

2015/04/25 12:39:53 workers: 500, corpus: 186 (42s ago), crashers: 3,
 	restarts: 1/8027, execs: 12009519 (121224/sec), cover: 2746, uptime: 1m39s

I’d like to save the hours ran, crashes, and suppressions so I can perform data basic analysis. So that if functionA had a crash line 10 and is fixed. We will be able to see that the function had another crash in line 15. 

I do think we can get all of it from the /working-dir though. Maybe just save that in a volume?

I really like your plan of a separate fuzz directory.
Future Plans:
I would eventually like a web front end that shows the output of Component 1 and cumulative running stats of the fuzzer, like hours, crashes, and suppressions.

Also would eventually like to add a ‘Create Jira Ticket’ button so we can create tickets of the crashes directly from the web app.
