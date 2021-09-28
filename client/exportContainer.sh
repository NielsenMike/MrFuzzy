#!/bin/bash

echo 'Running docker container export!'

containers=$(docker container ls --quiet)
now=$(date)

echo "$now : Script Started" >> /home/pi/logs/logs.txt

for container in $containers
do
        echo "Start exporting $container"
        docker export $container > $container.tar
        if [[ $? != 0 ]]; then
                echo "$now : Error creating Tar files" >> /home/pi/logs/logs.txt
        fi
        mv $container.tar /home/mf_exporter/tarfiles/
        if [[ $? != 0 ]]; then
                echo "$now : Error moving Tar files to directory" >> home/pi/logs/logs.txt
        fi
        echo "Finished"
done

echo "$now : Script Finished" >> /home/pi/logs/logs.txt
echo " " >> /home/pi/logs/logs.txt

