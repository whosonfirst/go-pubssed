#!/bin/sh

export PATH="${PATH}:/usr/local/bin"

PYTHON=`which python`
GOLANG=`which go`

if [ ! -x ${PYTHON} ]
then
    echo "Missing Python binary"
    exit 1
fi

if [ ! -x ${GOLANG} ]
then
    echo "Missing go binary"
    exit 1
fi

WHOAMI=`${PYTHON} -c 'import os, sys; print os.path.realpath(sys.argv[1])' $0`

SYSTEMD=`dirname ${WHOAMI}`
GO_SWIM=`dirname ${SYSTEMD}`

USER="pubssed"
GROUP="pubssed"

PUBBSED_SERVICE="/lib/systemd/system/pubssed-server.service"

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit 1
fi

if getent passwd ${USER} > /dev/null 2>&1; then
    echo "${USER} user account already exists"
else
    useradd ${USER} -s /sbin/nologin -M
fi

if [ ! -d /usr/local/etc/swim ]
then
    mkdir -p /usr/local/etc/swim
    chown swim /usr/local/etc/swim
    chmod 700 /usr/local/etc/swim
fi

cd ${GO_SWIM}
${GOLANG} build -mod vendor -o /usr/local/bin/pubssed-server cmd/pubssed-server/main.go
cd -

for SERVICE in ${PUBSSED_SERVICE}
do

    SERVICE_FNAME=`basename ${SERVICE}`

    if [ -f ${SERVICE} ]
    then
	cp ${SYSTEMD}/${SERVICE_FNAME}.example ${SERVICE}
	sudo chmod 644 ${SERVICE}

	echo ""
	echo "system stuff installed - you will still need to run the following, manually:"
	echo "	sudo systemctl daemon-reload"
	echo "	sudo systemctl restart ${SERVICE_FNAME}"

    else
	cp ${SYSTEMD}/${SERVICE_FNAME}.example ${SERVICE}
	sudo chmod 644 ${SERVICE}

	echo ""
	echo "system stuff installed - you will still need to run the following, manually:"
	echo "	sudo systemctl enable ${SERVICE_FNAME}"
	echo "	sudo systemctl start ${SERVICE_FNAME}"

    fi

done
    
exit 0
