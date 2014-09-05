deb:
	GOOS=linux GOARCH=amd64 go build -o usr/bin/go-puppet-forge
	fpm -f -n go-puppet-forge -s dir -t deb \
		--workdir debian \
		--version `git describe --tags --long` \
		--deb-upstart debian/upstart/go-puppet-forge \
		--after-install debian/postinst usr/bin/
	rm -r usr