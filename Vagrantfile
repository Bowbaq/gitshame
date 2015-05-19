# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.box_url = "https://cloud-images.ubuntu.com/vagrant/trusty/current/trusty-server-cloudimg-amd64-vagrant-disk1.box"

  config.vm.synced_folder ".", "/home/vagrant/go/src/github.com/Bowbaq/gitshame"

  config.vm.network "forwarded_port", guest: 3000, host: 3000

  config.vm.provider "virtualbox" do |v|
    v.memory = 4096
    v.cpus = 4
  end

  config.vm.provision "shell", inline: <<-SHELL
    set -x

    sudo apt-get -qy update
    sudo apt-get -qy upgrade
    sudo apt-get -qy install git mercurial bzr postgresql postgresql-contrib
    sudo gem install sass coffee-script

    sudo -u postgres psql -c "CREATE USER gitshame WITH PASSWORD 'gitshame';"
    sudo -u postgres createdb -O gitshame gitshame

    GO_VERSION=1.4.2

    wget https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz
    tar xf godeb-amd64.tar.gz
    ./godeb install $GO_VERSION
    rm godeb-amd64.tar.gz *.deb

    echo 'export GOPATH=/home/vagrant/go' >> /home/vagrant/.profile
    echo 'export PATH=$PATH:$GOPATH/bin' >> /home/vagrant/.profile
    source /home/vagrant/.profile

    cd go/src/github.com/Bowbaq/gitshame

    go get github.com/tools/godep
    go get github.com/codegangsta/gin
    godep restore

    chown -R vagrant:vagrant $GOPATH
  SHELL
end
