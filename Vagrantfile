client_vagrantfile = File.expand_path('../dokku-src/Vagrantfile', __FILE__)
load client_vagrantfile


Vagrant::configure("2") do |config|
  config.vm.define "dokku", primary: true do |vm|
    vm.vm.synced_folder File.expand_path('../bin-build', __FILE__), "/var/lib/dokku/plugins/available/dokku-enterprise"
    vm.vm.provision :shell do |s|
      s.inline = <<-EOT
        sudo dokku plugin:enable dokku-enterprise
        apt-get install nginx-extras -y
      EOT
    end
  end
end