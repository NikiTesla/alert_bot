#!/bin/bash

ARCTICDEM_DIR=$1
FOUNDATION_DIR=$2
SOVIETDEM_DIR=$3

echo "articdem $ARCTICDEM_DIR"
echo "foundation $FOUNDATION_DIR"
echo "sovietdem $SOVIETDEM_DIR"

export HOST=localhost
export PORT=2704

foundation_file=$(find $FOUNDATION_DIR -maxdepth 2 -type f -iregex ".*/basSETSM.*/SETSM.*dem.tif")
articdem_files=$(find $ARCTICDEM_DIR -maxdepth 2 -type f -iregex ".*/SETSM.*/SETSM.*dem.tif")
sovietdem_file=$(find $SOVIETDEM_DIR -maxdepth 1 -type f -iregex ".*.tif")

success_counter=0
failed_counter=0

for articdem_file in $articdem_files; do
    codem $foundation_file $articdem_file || res=1
    if [[ $res == "1" ]]; then
        failed_counter=$((failed_counter+1))
        curl -d "$articdem_file failed. Total failed: $failed_counter. Total success $success_counter" $HOST:$PORT/result
    else
        success_counter=$((success_counter+1))
        curl -d "$articdem_file succeeded. Total failed: $failed_counter. Total success $success_counter" $HOST:$PORT/result
    fi
done

