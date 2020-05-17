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

function build_images() {
  echo "building images..."

  cd $DIR &&
  docker-compose build
}

function remove_images() {
  echo "removing images..."

  cd $DIR &&
  sudo docker-compose down
}

function run_containers() {
  echo "running containers..."

  sudo docker-compose up &&
  git daemon --export-all --enable=upload-pack --reuseaddr --base-path=. &
}

function update_containers() {
  # コンテナ内のリポジトリを更新 & コマンド再起動
  echo "updating containers..."

  # アプリ一時停止
  sudo docker-compose exec api-server $REPO_DIR/dev-container.sh stop
  sudo docker-compose exec api-server $REPO_DIR/dev-container.sh wait

  cd $DIR &&
  sudo docker-compose exec api-server $REPO_DIR/dev-container.sh pull
  sudo docker-compose exec api-server $REPO_DIR/dev-container.sh build

  # アプリ再開
  sudo docker-compose exec api-server $REPO_DIR/dev-container.sh resume
}

function stop_containers() {
  echo "stopping containers..."
  sudo docker-compose stop
}


case $MODE in
"build")
  build_images
  ;;
"remove")
  remove_images
  ;;
"run")
  run_containers
  ;;
"stop")
  stop_containers
  ;;
"update")
  update_containers
  ;;
*)
  echo "options: build, remove, run, stop, update"
esac


