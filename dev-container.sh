#!/bin/bash

# DIR=`pwd`
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )" # このスクリプトの場所
MODE=$1

REPO_DIR=/home/dev/localchat
MAIN_BRANCH=master
APP=cmd

FLAG_FILE=/home/wait
CMD=/go/bin/cmd
CMD_PID=/home/cmd.pid

function build_app() {
  echo "building $CMD ..."
  git reset --hard origin/$MAIN_BRANCH &&
  go install ./${APP}
}
function run_app() {
  echo "running $CMD ..."

  # turn on bash's job control
  set -m

  if [ -e $CMD_PID ]; then
    rm $CMD_PID
  fi

  while true; do
    while [ -e $FLAG_FILE ]; do
      echo "run.sh: waiting..."
      sleep 3
    done

    $REPO_DIR/stop.sh
    $CMD &
    echo $! > $CMD_PID &&
    fg %1

    sleep 3
  done
}
function stop_app() {
  if [ -e $CMD_PID ]; then
    echo "stopping $CMD ..."
    kill `cat $CMD_PID`
  fi
}
function resume_app() {
  echo "resume $CMD ..."
  rm $FLAG_FILE
}
function wait_app() {
  echo "waiting $CMD ..."
  touch $FLAG_FILE
}
function pull_app() {
  echo "pulling in $REPO_DIR ..."
  cd $REPO_DIR &&
  git fetch origin $MAIN_BRANCH &&
  git reset --hard origin/$MAIN_BRANCH
}


case $MODE in
"build")
  build_app
  ;;
"run")
  run_app
  ;;
"stop")
  stop_app
  ;;
"resume")
  resume_app
  ;;
"wait")
  wait_app
  ;;
"pull")
  pull_app
  ;;
*)
  echo "options: build, remove, run, stop, update"
esac


