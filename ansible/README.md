# PartyHall deploy scripts

This repository contains the ansible scripts required to setup a PartyHall booth

## Usage
Flash an SD card with Raspbian lite (64 BITS !!!) and boot it up with an ethernet connection.

Do the first time setup:
- Create a `pi` account with a known password
- `sudo raspi-config`
    - Set your locale / timezone / keyboard layout accordingly
    - Enable console autologin
    - Enable SSH
    - Reboot 

Clone the repository:
```
$ git clone https://github.com/partyhall/partyhall.git
$ cd partyhall/ansible
```

Fill the inventory correctly:
```
$ nvim inventories/hosts
partyhall ansible_user=pi ansible_password=[[ YOUR RPI PASSWORD ]] ansible_host=[[ YOUR RPI ADDRESS ]] ansible_port=22
```

You can add as many hosts as you have partyhall built. You then need to setup the config in the `inventories/host_vars/{HOST_NAME}.yml` file.

Additional settings are possible, defaults values are set in `inventories/group_vars/all.yml`. You should not edit this file, rather copy the values in the host file and update the value to the ones you want.

Once done you can process with the ansible script:
```sh
$ ansible-galaxy install -r requirements.yaml
$ ansible-playbook -i inventories/hosts setup.yaml
```

Reboot your Pi and you should be good to go!