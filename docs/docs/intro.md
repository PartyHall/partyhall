---
sidebar_position: 1
---

# PartyHall

Let's discover the software side of PartyHall

## Getting Started

You're done building your booth (that's a lie, you'll find a lot of things to add to it) and now you need to setup the software. Let's do this together !

## The RPi side-note

While I started using this software, I now moved on to have better performances to an old motherboard (I'm using a i5-3570k with its internal graphics).

As of the last time I've tested it, it worked ok-ish on a Raspberry Pi 3 but you'll be better off with a newer one.

Also note that the install scripts are NOT updated to be used on the RPi so they will probably won't work. I need to change them but I'm short on time for now so it might not work for you.

## Concepts

### Export
Once your party is over, you can request an export from the Admin UI. This will generate a zip-file containing all the pictures from the photobooth, a timelapse if you did not disable it, and any module's thing to export.

### Hotspot
The appliance, when configured during the install, will emit a Wifi hotspot that you can connect to.

It does NOT share the access to the internet (yet?).

This is available even when the appliance is not connected to ethernet.

You can access the UI from it through the url [http://partyhall.lan](http://partyhall.lan) or the hostname you chose if this is not the same.

### PHK files
Those are PartyHall Karaoke files. These are song-packs to load directly. They basically are zip files with a custom structure.

### Timelapse
This is a short video running at 6fps that recaps the full party. The pictures in it are the [unattended](#unattended) pictures.

### Unattended
The photobooth has an unattended mode, this will automatically take pictures every 5 minutes (configurable). Those are used to generate a [timelapse](#timelapse) in the [export](#export).