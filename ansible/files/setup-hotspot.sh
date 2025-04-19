#!/bin/bash

set -e

# This file is part of PartyHall appliance software
# Its licence applies
# Learn more at https://github.com/PartyHall/partyhall

# This script aims at setting up hostapd + iptables
# so that the appliance provides an hotspot that
# shares the ethernet internet connection

if [ $# -ne 4 ]; then
    echo "Usage: $0 <ETH_INTERFACE> <WIFI_INTERFACE> <SSID> <PASSWORD>"
    exit 10
fi

ETH_IF="$1"
WIFI_IF="$2"
SSID="$3"
AP_PWD="$4"
HOSTAPD_CONF="/etc/hostapd/hostapd.conf"

#region Check that interfaces exist
if ! ip link show "$ETH_IF" >/dev/null 2>&1; then
    echo "Error: Ethernet interface $ETH_IF not found"
    exit 2
fi

if ! ip link show "$WIFI_IF" >/dev/null 2>&1; then
    echo "Error: WiFi interface $WIFI_IF not found"
    exit 3
fi
#endregion

#region Validating SSID/Password
if ! [[ "$SSID" =~ ^[a-zA-Z0-9\ ]{1,32}$ ]]; then
    echo "Invalid SSID (1-32 alphanumerical characters + spaces)"
    exit 5
fi

if ! [[ "$AP_PWD" =~ ^[a-zA-Z0-9]{12,63}$ ]]; then
    echo "Invalid password (12-63 alphanumerical characters)"
    exit 6
fi
#endregion

systemctl stop hostapd
systemctl stop dnsmasq

#region Setting up interface IP

# Cleaning old hotspot interfaces, just in case we change iface
# for the new hotspot
PREVIOUS_IPS=$(ip -br addr show | grep "192\.168\.203\.1/24" | awk '{print $1}')

for INTERFACE in $PREVIOUS_IPS; do
    ip link set dev "$INTERFACE" down 1>/dev/null || true
    ip addr flush dev "$INTERFACE" 2>/dev/null || true
done

# Just in case
sleep 0.5

# Clearing old ip assigned to the new interface
ip addr flush dev "$WIFI_IF" 2>/dev/null || true

# Setting the base IP for the new interface
ip addr add 192.168.203.1/24 broadcast 192.168.203.255 dev "$WIFI_IF"

# Finally we start the interface
ip link set dev "$WIFI_IF" up

if ! ip addr show "$WIFI_IF" | grep -q "192.168.203.1"; then
  echo "Failed to assign IP to $WIFI_IF"
  exit 9
fi
#endregion

#region Flush existing rules
iptables -F wifiApRules
iptables -t nat -F wifiApNatRules
#endregion

#region Configure rules
iptables -A wifiApRules -i "$WIFI_IF" -o "$ETH_IF" -j ACCEPT
iptables -A wifiApRules -i "$ETH_IF" -o "$WIFI_IF" -m state --state RELATED,ESTABLISHED -j ACCEPT
iptables -t nat -A wifiApNatRules -o "$ETH_IF" -j MASQUERADE
#endregion

echo "Hotspot rules updated for $WIFI_IF -> $ETH_IF"

#region Setup hostapd config file
cat > "$HOSTAPD_CONF" <<EOF
# Basic Configuration

interface=$WIFI_IF
driver=nl80211
country_code=FR
channel=6

# SSID and Authentication
ssid=$SSID
wpa_passphrase=$AP_PWD

ignore_broadcast_ssid=0
auth_algs=1
wpa=2
wpa_key_mgmt=WPA-PSK
rsn_pairwise=CCMP

# Settings
hw_mode=g
ieee80211n=1
ht_capab=[HT40+][SHORT-GI-20][DSSS_CCK-40]
macaddr_acl=0
ieee80211d=1
wmm_enabled=1
EOF
#endregion

#region Configure dnsmasq
cat > "/etc/dnsmasq.d/01-interface.conf" <<EOF
interface=$WIFI_IF
bind-interfaces
dhcp-range=192.168.203.2,192.168.203.20,255.255.255.0,24h
dhcp-option=option:router,192.168.203.1
dhcp-option=option:dns-server,192.168.203.1
EOF
#endregion

echo "Updated dnsmasq config"

if ! systemctl start hostapd; then
    echo "Failed to restart hostapd service"
    exit 7
fi

if ! systemctl start dnsmasq; then
    echo "Failed to restart dnsmasq service"
    exit 8
fi

echo "Wifi access point restarted"