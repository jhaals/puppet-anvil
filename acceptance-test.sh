#!/bin/bash

CWD=$(pwd)

ROOT_PATH=./tmp-test/accep

type puppet >/dev/null 2>&1 || { echo >&2 "I require puppet but it's not installed.  Aborting."; exit 1; }

mkdir -p $ROOT_PATH/out $ROOT_PATH/modules/puppetlabs/apache $ROOT_PATH/modules/puppetlabs/concat $ROOT_PATH/modules/puppetlabs/stdlib


MODULEPATH=$ROOT_PATH/modules PORT=8080 ./build/puppet-anvil > /dev/null 2>&1  &
PID=$!

if [ ! -f $ROOT_PATH/puppetlabs-apache-1.5.0.tar.gz 2>&1 >/dev/null ]; then
	wget -q https://forgeapi.puppetlabs.com/v3/files/puppetlabs-apache-1.5.0.tar.gz -O $ROOT_PATH/puppetlabs-apache-1.5.0.tar.gz >/dev/null
	wget -q https://forgeapi.puppetlabs.com/v3/files/puppetlabs-concat-1.2.3.tar.gz -O $ROOT_PATH/puppetlabs-concat-1.2.3.tar.gz >/dev/null
	wget -q https://forgeapi.puppetlabs.com/v3/files/puppetlabs-stdlib-4.6.0.tar.gz -O $ROOT_PATH/puppetlabs-stdlib-4.6.0.tar.gz >/dev/null
fi
curl -s -X PUT http://localhost:8080/admin/module/puppetlabs-apache-1.5.0.tar.gz -T $ROOT_PATH/puppetlabs-apache-1.5.0.tar.gz >/dev/null
curl -s -X PUT http://localhost:8080/admin/module/puppetlabs-concat-1.2.3.tar.gz -T $ROOT_PATH/puppetlabs-concat-1.2.3.tar.gz >/dev/null
curl -s -X PUT http://localhost:8080/admin/module/puppetlabs-stdlib-4.6.0.tar.gz -T $ROOT_PATH/puppetlabs-stdlib-4.6.0.tar.gz >/dev/null


puppet module install puppetlabs/apache --module_repository http://127.0.0.1:8080 --modulepath $ROOT_PATH/out >/dev/null
RC=$?
if [ $RC == 0 ]; then
	echo PASS
else
	echo FAIL
fi


sudo rm -rf $ROOT_PATH/modules $ROOT_PATH/out

kill $PID
exit $RC
