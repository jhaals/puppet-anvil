# Puppet Anvil

Puppet Anvil is a minimal implementation of the Puppet Forge and does not require any database.
Puppet modules can then be downloaded using the [Puppet module tool](https://docs.puppetlabs.com/puppet/latest/reference/modules_installing.html#installing-from-another-module-repository) or [librarian-puppet](http://librarian-puppet.com/)

### Run in Docker Container
Pull container `docker pull jhaals/puppet-anvil`

Serve modules on port 8080 from /var/lib/modules

    docker run -v /var/lib/modules:/modules -p 8080:8080 jhaals/puppet-anvil

Modules must be stored in the following directory structure `user/module/user-module-version.tar.gz`
example:

    /yourmoduledir/puppetlabs/apache/puppetlabs-apache-1.1.0.tar.gz

##### Build and run from source

    go build -o puppet-anvil
    PORT=1337 MODULEPATH=/var/lib/puppet-anvil/modules ./puppet-anvil

    Starting Puppet Anvil on port 1337 serving modules from /var/lib/puppet-anvil/modules

_You can create a .deb package for Ubuntu using `make deb`. the fpm gem is required._

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

#### Managing your repo
There are two ways to manage your module artifacts:

* Manually land files in your module directory

	cp puppetlabs-apache-1.5.0.tar.gz /var/lib/puppet-anvil/modules/puppetlabs/apache/puppetlabs-apache-1.5.0.tar.gz

* Use the supplied `admin/module` endpoint

	curl -s -X PUT http://localhost:8080/admin/module/puppetlabs-apache-1.5.0.tar.gz -T ./puppetlabs-apache-1.5.0.tar.gz


