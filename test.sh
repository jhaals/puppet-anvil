#!/bin/bash

CWD=$(pwd)


type puppet >/dev/null 2>&1 || { echo >&2 "I require puppet but it's not installed.  Aborting."; exit 1; }


sudo rm -rf ./test/out ./test/modules
mkdir -p test/out test/modules/puppetlabs/apache test/modules/puppetlabs/concat test/modules/puppetlabs/stdlib

go build -o puppet-anvil && echo "'puppet-anvil' built"




MODULEPATH=./test/modules PORT=8080 ./puppet-anvil &
PID=$!

if [ ! -f ./test/puppetlabs-apache-1.5.0.tar.gz 2>&1 >/dev/null ]; then
	wget https://forgeapi.puppetlabs.com/v3/files/puppetlabs-apache-1.5.0.tar.gz -O ./test/puppetlabs-apache-1.5.0.tar.gz
	wget https://forgeapi.puppetlabs.com/v3/files/puppetlabs-concat-1.2.3.tar.gz -O ./test/puppetlabs-concat-1.2.3.tar.gz
	wget https://forgeapi.puppetlabs.com/v3/files/puppetlabs-stdlib-4.6.0.tar.gz -O ./test/puppetlabs-stdlib-4.6.0.tar.gz
fi
curl -i -X PUT http://localhost:8080/admin/module/puppetlabs-apache-1.5.0.tar.gz -T ./test/puppetlabs-apache-1.5.0.tar.gz
curl -i -X PUT http://localhost:8080/admin/module/puppetlabs-concat-1.2.3.tar.gz -T ./test/puppetlabs-concat-1.2.3.tar.gz
curl -i -X PUT http://localhost:8080/admin/module/puppetlabs-stdlib-4.6.0.tar.gz -T ./test/puppetlabs-stdlib-4.6.0.tar.gz


puppet module install puppetlabs/apache --module_repository http://127.0.0.1:8080 --modulepath ./test/out

kill $PID
