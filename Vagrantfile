Vagrant.configure(2) do |config|

  config.vm.box = "http://opscode-vm-bento.s3.amazonaws.com/vagrant/virtualbox/opscode_centos-7.0_chef-provisionerless.box"

  config.vm.provision "shell", privileged: false, inline: <<-EOF
    sudo yum -y install rpmdevtools mock
    rpmdev-setuptree

    export BASE=/vagrant
    ln -s $BASE/ec2-metadata-environment.spec $HOME/rpmbuild/SPECS/
    ln -s $BASE/pkg/linux_amd64/update-ec2-metadata-environment $HOME/rpmbuild/SOURCES
    ln -s $BASE/ec2-metadata-environment.service $HOME/rpmbuild/SOURCES

    rpmbuild -bb rpmbuild/SPECS/ec2-metadata-environment.spec

    find $HOME/rpmbuild -type d -name "RPMS" -exec cp -r {} $BASE/ \\;
  EOF

end
