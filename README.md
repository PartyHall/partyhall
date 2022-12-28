# PartyHall

## Install
Flash an SD card with Raspbian lite (64 BITS !!!) and boot it up with an ethernet connection.

Do the first time setup:
- Create a `pi` account with a known password
- `sudo raspi-config`
    - Set your locale / timezone / keyboard layout accordingly
    - Enable console autologin
    - Enable SSH
    - Reboot 

If you have an older Pi (Running 32 bit operating system) you'll have a few more steps (Expanding storage, ...) along with recompiling the partyhall for your architecture. This is not supported and no help will be given. I'm not even sure it works on 32bits tbh I had no thought about it while coding it.

Clone the repository:
```
$ git clone https://github.com/partyhall/partyhall-v3.git
$ cd partyhall-v3/ansible
```

Fill the inventory correctly:
```
$ nvim inventories/hosts
partyhall ansible_user=raspberry ansible_password=[[ YOUR RPI PASSWORD ]] ansible_host=[[ YOUR RPI ADDRESS ]] ansible_port=22
```

You can add as many hosts as you have partyhall built. You then need to setup the config in the `ansible/inventories/host_vars/{HOST_NAME}.yml` file.

Additional settings are possible, defaults values are set in `ansible/inventories/group_vars/all.yml`. You should not edit this file, rather copy the values in the host file and update the value to the ones you want.

Once done you can process with the ansible script:
```sh
$ ansible-galaxy install -r requirements.yaml
$ ansible-playbook -i inventories/hosts setup.yaml
```

Reboot your Pi and you should be good to go!