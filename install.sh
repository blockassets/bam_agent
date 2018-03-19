#!/usr/bin/env bash

# This script requires sshpass and parallel
#
passwd="bwcon"
export SSHPASS="${passwd}"

export SERVICE="bam_agent-linux-arm"

if [ -e "./workers.txt" ] ; then
	WORKERS=`cat ./workers.txt`
fi

if [ -z "${WORKERS}" ] ; then
	echo "Need some workers to install to!"
	exit 1
fi

dowork() {
	ipaddr=$1
	echo "----------- ${ipaddr} start"
	sshpass -e scp -o StrictHostKeychecking=no ${SERVICE}.gz root@$ipaddr:/usr/bin
	sshpass -e scp -o StrictHostKeychecking=no ${SERVICE}.service root@$ipaddr:/etc/systemd/system
	sshpass -e ssh -o StrictHostKeychecking=no root@$ipaddr "rm -f /usr/bin/${SERVICE}; gunzip /usr/bin/${SERVICE}.gz; chmod ugo+x /usr/bin/${SERVICE}; systemctl enable ${SERVICE}; systemctl stop ${SERVICE}; systemctl start ${SERVICE}"
	echo "----------- ${ipaddr} finish"
}

export -f dowork

parallel --no-notice dowork ::: ${WORKERS}
