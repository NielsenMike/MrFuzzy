#!/bin/bash

#Find & Send  Files
FILES="/home/mf_exporter/tarfiles/*"
now=$(date)

echo "$now : Script Started" >> /home/mf_exporter/logs/logs.txt

if [[ -z "$(ls -A /home/mf_exporter/tarfiles/)" ]]; then
        echo "$now : Tarfiles Directory Empty, Ending..." >> /home/mf_exporter/logs/logs.txt
else
        echo "not empty"
        for tars in $FILES
        do
                echo "Sending $tars file"
                scp $tars mf_exporter@wiproh21-mnielsen.el.eee.intern:/home/mf_exporter/raspberry1
                if [[ $? != 0 ]]; then
                        echo "$now : scp $? : Error sending files" >> home/mf_exporter/logs/logs.txt
                fi 
                echo "File Sent!"
done
fi

#Activate the Hashing Process on Server
curl -X POST http://wiproh21-mnielsen.el.eee.intern:8080/writeDataIntoDB

echo "$now : Script Finished " >> /home/mf_exporter/logs/logs.txt
echo " " >> /home/mf_exporter/logs/logs.txt
