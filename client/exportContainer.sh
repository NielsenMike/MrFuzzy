#!/bin/bash

echo '#Running docker container export!'

containers=$(docker container ls --quiet)

for container in $containers
do
	echo "Start exporting $container"
	docker export $container > $container.tar
	echo "Finished exporting"
done

