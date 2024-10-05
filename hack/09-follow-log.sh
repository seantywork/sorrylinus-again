#!/bin/bash

function usage_helper(){

    echo "usage:"
    echo "  again : follow sorrylinus-again log"
    echo "  hub   : follow sorrylinus-hub log"
}

if [ -z "$1" ]
then 

    echo "argument is empty"

    usage_helper
    
    exit 1
fi



if [ "$1" == "again" ]
then 

    TARGET_POD=$(sudo kubectl -n frank get pods --no-headers -o custom-columns=":metadata.name" | grep 'sorrylinus-again')
    
    echo "your follow target name: ${TARGET_POD}"
    
    sudo kubectl -n frank logs "${TARGET_POD}" --follow


elif [ "$1" == "hub" ]
then 

    TARGET_POD=$(sudo kubectl -n frank get pods --no-headers -o custom-columns=":metadata.name" | grep 'sorrylinus-hub')
    
    echo "your follow target name: ${TARGET_POD}"
    
    sudo kubectl -n frank logs "${TARGET_POD}" --follow

else 

    echo "wrong argument: $1"

    usage_helper

    exit 1

fi