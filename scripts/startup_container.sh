#!/bin/bash

required_args=2
if [ "$#" -ne "$required_args" ]; then
    echo "Usage: $0 <arg1> <arg2>"
    exit 1
fi

uname=$(uname -m)

ora_pass=$2
ora_db=$1

image="gvenzl/oracle-xe:18"
arch="amd64"
args="-e ORACLE_PASSWORD=$ora_pass -e ORACLE_DATABASE=$ora_db"
container_name="mentor-ora"

if [ "$uname" = "arm64" ]; then
  arch="arm64"
fi

echo "Using arch=$arch"
echo "args: ORACLE_PASSWORD=$ora_pass, ORACLE_DATABASE=$ora_db"

echo "Checking if colima is running..."
echo "Colima must be running: 'colima start --arch x86_64 --memory 4'"
status=$(colima status 2>&1 | awk -F 'msg=' '{print $2}' | awk -F '"' '{print $2}')

if [ "$status" = "colima is not running" ]; then
    echo "Colima is not running."
    echo "Using docker runtime."
fi

echo "Using container image=$image"

oracle_stopped=$(docker ps --filter "status=exited" --filter "name=$container_name" --quiet 2> /dev/null | wc -l)
if [ "$oracle_stopped" -ne 0 ]; then
    docker rm -f mentor-ora
fi

oracle_started=$(docker ps --filter "name=$container_name" --quiet 2> /dev/null | wc -l)
if [ "$oracle_started" -eq 0 ]; then
    docker run --name mentor-ora -v $(pwd)/tmp:/tmp -d $args -p 1521:1521 $image
fi

search_string="DATABASE IS READY TO USE"

loading_animation() {
    local -a frames=('-' '\' '|' '/')
    local i=0
    while ! docker logs "$container_name" | grep -q "$search_string"; do
        printf "\r[${frames[$i]}] Waiting for container to start..."
        ((i = (i + 1) % ${#frames[@]}))
        sleep 0.2
    done
    printf "\rOracle container started!                \n"
}

loading_animation


