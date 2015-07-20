# Puppet Anvil

Puppet Anvil is a minimal implementation of the Puppet Forge and does not require any database.
Puppet modules can then be downloaded using the [Puppet module tool](https://docs.puppetlabs.com/puppet/latest/reference/modules_installing.html#installing-from-another-module-repository) or [librarian-puppet](http://librarian-puppet.com/)


##### Build and run from source

    make build
    PORT=1337 MODULEPATH=/var/lib/puppet-anvil/modules ./build/puppet-anvil

    Starting Puppet Anvil on port 1337 serving modules from /var/lib/puppet-anvil/modules

Modules must be stored in the following directory structure `user/module/user-module-version.tar.gz`
example:

    /yourmoduledir/puppetlabs/apache/puppetlabs-apache-1.1.0.tar.gz

#### Usage with Puppet
A custom module_repository can be specified in the puppet config file.

    module_repository=http://my-forge.com/

Or directly on command line

    ~ puppet module install puppetlabs/apache --module_repository http://127.0.0.1:8080 --modulepath modules
    Notice: Preparing to install into /Users/jhaals/modules ...
    Notice: Created target directory /Users/jhaals/modules
    Notice: Downloading from http://127.0.0.1:8080 ...
    Notice: Installing -- do not interrupt ...
    /Users/jhaals/modules
    └─┬ puppetlabs-apache (v1.1.0)
      ├── puppetlabs-concat (v1.1.0)
      └── puppetlabs-stdlib (v4.2.2)

This project is inspired by [simple-puppet-forge](https://github.com/dalen/simple-puppet-forge)

#### Hacking

	make deps
	make test
	make build
	
	# acceptance-test uses your installed copy of puppet to excercise the anvil server
	make acceptance-test

