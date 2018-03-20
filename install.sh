#!/usr/bin/env bash

# This script requires sshpass and parallel
#
passwd="bwcon"
export SSHPASS="${passwd}"

export SERVICE="bam_agent"
export BINARY="${SERVICE}-linux-arm"
export INSTALL_DIR="/usr/bin"

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
	sshpass -e scp -o StrictHostKeychecking=no ${BINARY}.gz root@$ipaddr:${INSTALL_DIR}
	sshpass -e scp -o StrictHostKeychecking=no bam_agent.service root@$ipaddr:/etc/systemd/system
	sshpass -e ssh -o StrictHostKeychecking=no root@$ipaddr "systemctl daemon-reload; systemctl stop ${SERVICE}; rm -f ${INSTALL_DIR}/${BINARY}; gunzip ${INSTALL_DIR}/${BINARY}.gz; chmod ugo+x ${INSTALL_DIR}/${BINARY}; systemctl enable ${SERVICE}; systemctl start ${SERVICE}"
	echo "----------- ${ipaddr} finish"
}

export -f dowork

parallel --no-notice dowork ::: ${WORKERS}
