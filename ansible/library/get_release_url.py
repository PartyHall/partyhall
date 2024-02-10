#!/usr/bin/python

from ansible.module_utils.basic import *
import requests
import platform

# -linux-xxxx-default.tar.gz = spotifyd
CUSTOM_ARCH_NAMES = {
    'x86_64': ['amd64', '-linux-default.tar.gz'],
    'aarch64': ['arm64', '-linux-armhf-default.tar.gz'],
}

def main():

    fields = {
        'repository': { "required": True, "type": "str" }
    }

    module = AnsibleModule(argument_spec=fields)

    arch = platform.uname()[4]
    if arch in CUSTOM_ARCH_NAMES.keys():
        arch = CUSTOM_ARCH_NAMES[arch]

    resp = requests.get(f'https://api.github.com/repos/{module.params["repository"]}/releases/latest')
    
    if resp.status_code != 200:
        module.fail_json(changed=False, msg=f'Failed to get the Github releases: {resp.status_code}')

    data = resp.json()

    for asset in data['assets']:
        for curr_arch in arch:
            if asset['name'].endswith(curr_arch):
                module.exit_json(changed=False, meta=asset['browser_download_url'])

    module.fail_json(changed=False, msg=f'Binary not available for this architecture: {arch}')

if __name__ == '__main__':
    main()