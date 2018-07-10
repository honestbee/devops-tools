#!/bin/bash
function replace() {
    sed -i "/old:$/,/^$/d" $3
    sed -i "s;$1;$2;g" $3
}

function commit() {
    cd $1
    currentBranch=`git symbolic-ref --short HEAD`
    newBranch="2018-07-10-devops-remove-old-registry"
    git checkout -b $newBranch $currentBranch
    git add .drone.yml
    git commit -m "remove old registry"
    git push origin $newBranch
    hub pull-request -b $currentBranch -m "remove old registry" >> ../../pr.txt
    cd ../..
}

while read line; do
    folder=`echo $line | cut -d"," -f1`
    newRegistry=`echo $line | cut -d"," -f3`
    droneFile="working/$folder/.drone.yml"
    if [[ ! -f $droneFile ]]; then
        continue
    fi
    oldRegistry=`cat $droneFile | grep -oP "registry.honestbee.com/[a-zA-Z]+/?[a-zA-z]+" | uniq`
    if [[ "$oldRegistry" != "" ]]; then
        replace $oldRegistry $newRegistry $droneFile
        commit working/$folder
    fi
done < 'repos.csv'

