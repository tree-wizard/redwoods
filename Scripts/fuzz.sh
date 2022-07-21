#!/bin/bash

set -e

function build_fuzz_archive() {
  PROJECTPACKAGES=$1
  FUZZ_BASE_WORKDIRECTORY=$2
  FUZZ_WORKSPACE=$3

  for pkg in $PROJECTPACKAGES; do
    pushd $FUZZ_WORKSPACE/$pkg
    for file in *fuzz.go; do
      NAME=$(echo $file | sed 's/\_fuzz.go$//1')
      echo "Building zip file for $pkg/$NAME" "go-fuzz-build -func" "Fuzz_$NAME" "-o" "$FUZZ_BASE_WORKDIRECTORY$pkg-$NAME-fuzz.zip" " in " "$FUZZ_WORKSPACE/$pkg"
      #prepare workdir
      mkdir -p $FUZZ_BASE_WORKDIRECTORY/$pkg/
      go-fuzz-build -func "Fuzz_$NAME" -o "$FUZZ_BASE_WORKDIRECTORY/$pkg/$NAME-fuzz.zip" "$FUZZ_WORKSPACE/$pkg"
    done

    popd
  done
}

# timeout is a cross platform alternative to the GNU timeout command that
# unfortunately isn't available on macOS by default. see https://github.com/lightningnetwork/lnd/blob/master/scripts/fuzz.sh
timeout() {
    time=$1
    $2 &
    pid=$!
    sleep $time
    kill -s SIGINT $pid    
}

function run_fuzzer() {
  PROJECTPACKAGES=$1
  FUZZER_RUN_TIME=$2
  TIMEOUT_RUN_TIME=$3
  NOOFPROCESSES=$4
  FUZZ_BASE_WORKDIRECTORY=$5
  FUZZ_WORKSPACE=$6
  # For all of our defined packages
  for pkg in $PROJECTPACKAGES; do
    # prepare iteartion
    pushd $FUZZ_WORKSPACE/$pkg
    # for every fuzz.go file in our project
    for file in *fuzz.go; do
      NAME=$(echo $file | sed 's/\_fuzz.go$//1')
      WORKDIR=$FUZZ_BASE_WORKDIRECTORY/$pkg/$NAME
      # create workdir
      mkdir -p $WORKDIR
      mkdir -p $WORKDIR/Output
      mkdir -p $WORKDIR/Work
      # TODO: We should check if the build exists
      echo "Running fuzzer $pkg-$NAME-fuzz.zip with $NOOFPROCESSES processors for $FUZZER_RUN_TIME seconds with a timeout of $TIMEOUT_RUN_TIME"
      COMMAND="go-fuzz -bin=$FUZZ_BASE_WORKDIRECTORY/$pkg/$NAME-fuzz.zip -workdir=$WORKDIR/Work -procs=$NOOFPROCESSES -timeout=$TIMEOUT_RUN_TIME "
       echo "$COMMAND"
      timeout "$FUZZER_RUN_TIME" "$COMMAND"  &> $WORKDIR/Output/fuzzoutput.md
      echo "/endfuzz"
    done

    popd
  done
}
function run_fuzzer_docker() {
  PROJECTPACKAGES=$1
  FUZZER_RUN_TIME=$2
  TIMEOUT_RUN_TIME=$3
  NOOFPROCESSES=$4
  FUZZ_BASE_WORKDIRECTORY=$5
  FUZZ_WORKSPACE=$6
  # For all of our defined packages
  for pkg in $PROJECTPACKAGES; do
    # prepare iteartion
    pushd $FUZZ_WORKSPACE/$pkg
    # for every fuzz.go file in our project
    for file in *fuzz.go; do
      NAME=$(echo $file | sed 's/\_fuzz.go$//1')
      #prepare work directories on local machine
      WORKDIR=$FUZZ_BASE_WORKDIRECTORY/$pkg/$NAME/Work
      OUTPUTDIR=$FUZZ_BASE_WORKDIRECTORY/$pkg/$NAME/Output
      # create workdir
      mkdir -p $WORKDIR
      mkdir -p $OUTPUTDIR
      # TODO: We should check if the build exists
      echo "Spawning fuzzer $pkg-$NAME-fuzz.zip with $NOOFPROCESSES processors for $FUZZER_RUN_TIME seconds with a timeout of $TIMEOUT_RUN_TIME in docker"
      sudo docker run --name redwoods-0.0.1-livefuzz-$NAME -v $OUTPUTDIR:/app/redwoods/fuzz/$pkg/$NAME/Output -v $WORKDIR:/app/redwoods/fuzz/$pkg/$NAME/Work --rm redwoods:0.0.1 /bin/bash /app/redwoods/Scripts/fuzz.sh run $pkg $FUZZER_RUN_TIME  $TIMEOUT_RUN_TIME $NOOFPROCESSES /app/redwoods/fuzz /app/redwoods/workspace/vault  &
    done

    popd
  done
}
function usage() {
  echo "these Scripts are included in the Redwoods Fuzzing Suite and should not be executed alone"
}

# run input as subcommand
SUBCOMMAND=$1
shift

# we switch depending on entry
case $SUBCOMMAND in
buildarchives)
  build_fuzz_archive "$@"
  ;;
run)
  run_fuzzer "$@"
  ;;
fuzz_docker)
  run_fuzzer_docker "$@"
  ;;
*)
  usage
  exit 1
  ;;
esac