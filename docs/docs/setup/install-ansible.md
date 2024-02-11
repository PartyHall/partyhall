---
sidebar_position: 1
---

# Installing ansible

## What is ansible ?

The setup process is based on ansible, a tool commonly used to deploy software to server.

For our use case, the main benefits are that you can set everything up in one command, and you can repeat the command, be it for updating your appliance or because you messed around with files and want to go back to a clean state.


You'll need to install ansible on your main computer, more info can be found on [their docs](https://docs.ansible.com/ansible/latest/installation_guide/index.html).

## Preparing your computer

Clone the repository and go to the ansible folder:
```sh
$ git clone https://github.com/PartyHall/partyhall
$ cd partyhall/ansible
```

The setup depends on multiple ansible packages so you'll need to install them:
```sh
$ ansible-galaxy install -r requirements.yaml
$ ansible-galaxy collection install -r requirements.yaml
```

Note that those are two different commands, because of the way ansible is working, but you need both.

If everything went well, you can go to the next step.