---
sidebar_position: 2
---

# Pre-configuration

Most of the config is done before setting up the appliance. If you want to update the config, you can just run the command that deploy afterward.

We will need to modify a few files before running the deploy scripts.

## The host file

The file `inventories/hosts` lists all the machines you want to setup PartyHall on. For most of us, it will be only one.

There is not much in this file, you just need to update the `ansible_user` variable if you did not name your user `partyhall` and set the IP address to the one of your appliance.

## The config

I'm using ansible in a weird way. I'm using the *per-host* config file to store the settings you should keep for yourself. This includes passwords and API-keys. And the *group* config file to setup the other config.

I'm not using `ansible-vault` because this is still a public repository, and I don't want to have my secret being published there, encrypted or not.

If you are trying to setup PartyHall on multiple appliance at the same time, you should probably move those variable back to the next file. You can also put there variables from the other file you want to apply at one host only.

**Note**: You can do however you want, probably for easiness sake move everything to the `group_vars/all.yml` file. Its easier to backup that way.

Here's the rundown of available options

### General config
#### ethernet_interface
The ethernet interface name your appliance has. This settings is important, failing to set it correcly will result in being locked out (even through SSH!)

*Accepted values*
- Any networked interface (e.g. enp2s0, see `ip -br -c a`)

#### grub_resolution
The resolution accepted by your bootloader, depends on your screen's resolution. This is used for the splashscreen. In case of doubt, keep 720p

*Accepted values*
- Any valid resolution (e.g. `1280x720`)

#### hotspot_enabled
Enable or disable the Wifi access point.

*Accepted values*
- true
- false

#### hotspot_interface
The wifi interface name your appliance has. This is the device that will be used to broadcast your wifi access point.

*Accepted values*
- Any networked interface (e.g. wlan0, see `ip -br -c a`)

#### hotspot_ssid / hotspot_password / hotspot_wifi_channel

The info about the wifi hotspot, setting its name (ssid) and its password.

#### ntp_enabled
Synchronize the time with an NTP server

*Accepted values*
- true
- false

#### ntp_timezone
The timezone to sync to

*Accepted values*
- Any valid timezone, e.g. `Europe/Paris`

#### partyhall_admin_fullname

This is the displayed name for the admin user. That's the one used during the karaoke for "sung by".

**Note that the account is only created the first time and won't be updated** (This is subject to change, a PR will fix this later)

*Accepted values*
- Any string

#### partyhall_admin_username
This is the username for the admin account that will be created.

**Note that the account is only created the first time and won't be updated** (This is subject to change, a PR will fix this later)

*Accepted values*
- Any valid username string

#### partyhall_admin_password
This is the password for the admin account that will be created.

**Note that the account is only created the first time and won't be updated** (This is subject to change, a PR will fix this later)

*Accepted values*
- Any string

#### partyhall_debug_mode
Debug mode is used to bypass security features preventing the booth mode to be used outside of the appliance. It should only be set while contributing to PartyHall's source-code. Never on an actual appliance.

*Accepted values*:
- true
- false

#### partyhall_guests_allowed
This allows or disallow "guests". Those are people logging-in to the UI from their smartphone without an admin account. For now the only thing they can do is searching song for the karaoke, add them to the queue / change the queue order and use media-player controls.

*Accepted values*:
- true
- false

#### partyhall_language
This sets the language used by the appliance. Some errors coming from the backend are still not translated but the frontend should be ok.

*Accepted values*
- en
- fr

### partyhall_spotify_client_id / partyhall_spotify_secret
This is the client id / client secret given to you by spotify to reach their APIs. This is used to fetch the cover of the songs when you're adding one in the Admin.

To get your client id and secret, go to [https://developer.spotify.com/](https://developer.spotify.com/) and request one.

*Accepted values*
- Valid client id / client secret for the Spotify API

#### spotify_name
This sets the name of the appliance on Spotify Connect (Bottom left icon that choose which device you want to play your music on when using Spotify).

*Accepted values*
- Any valid spotifyd device name (Did not find docs about the valid characters)

### Photobooth-related config
#### partyhall_photobooth_has_hardware_flash
Feature not used yet. This will later be used to trigger a flash from an hardware flash (ledstrip with relay on the side of the screen).

*Accepted values*
- true
- false

#### partyhall_photobooth_default_timer
The amount of seconds between the button being pressed and the image being taken (The timer is shown on screen)

*Accepted values*
- Any integer value

#### partyhall_photobooth_unattended_inverval
The amounts of minutes between each unattended pictures.

*Accepted values*
- Any integer value
- **Recommended**: 5 (depends on the duration of the event, for a 10h night-long party, 5 is a good value, decrease for shorted ones).

#### partyhall_photobooth_webcam_width / partyhall_photobooth_webcam_height

The resolution of the webcam

*Accepted values*
- Any integer value

### Karaoke-related config
#### partyhall_karaoke_amt_songs_per_page
Temporary variable, this will be gone once I settled on the value, the pagination size returned by the API

*Accepted values*
- Any integer value

#### partyhall_karaoke_pre_play_timer
The duration in seconds while the title of the next song and the singer is displayed on the appliance

*Accepted values*
- Any integer value
- **Recommended**: 5