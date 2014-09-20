deb:
	GOOS=linux GOARCH=amd64 go build -o usr/bin/puppet-anvil
	fpm -f -n puppet-anvil -s dir -t deb \
		--workdir debian \
		--version `git describe --tags --long` \
		--deb-upstart debian/upstart/puppet-anvil \
		--after-install debian/postinst usr/bin/
	rm -r usr