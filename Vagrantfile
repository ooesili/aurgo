Vagrant.configure('2') do |config|
  config.vm.box = 'archlinux'

  config.vm.provision 'ansible' do |ansible|
    ansible.playbook = 'vagrant/playbook.yml'
  end

  config.vm.synced_folder '.', '/go/src/github.com/ooesili/aurgo'
end
