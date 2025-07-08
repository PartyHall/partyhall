import { ApSettings, InterfacesSettings } from './models/interfaces';
import { AudioDevice, AudioDevices } from './models/audio';
import NexusSettings from './models/nexus';
import { SDK } from './index';
import SpotifySettings from './models/spotify';

export default class Settings {
    sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async getButtonMappings(): Promise<Record<number, string>> {
        const resp = await this.sdk.post('/api/settings/button-mappings');
        return await resp.json();
    }

    public async getButtonMappingsActions(): Promise<string[]> {
        const resp = await this.sdk.get('/api/settings/button-mappings/actions');
        return await resp.json();
    }

    public async setButtonMappings(mappings: Record<number, string>): Promise<Record<number, string>> {
        const resp = await this.sdk.put('/api/settings/button-mappings', mappings);
        return await resp.json();
    }

    public async getNexus() {
        const resp = await this.sdk.get('/api/settings/nexus');
        return NexusSettings.fromJson(await resp.json());
    }

    public async setNexus(baseUrl: string, hwid: string, apiKey: string, bypassSsl: boolean) {
        const resp = await this.sdk.put('/api/settings/nexus', {
            base_url: baseUrl,
            hardware_id: hwid,
            api_key: apiKey,
            bypass_ssl: bypassSsl,
        });

        return NexusSettings.fromJson(await resp.json());
    }

    public async getSpotify() {
        const resp = await this.sdk.get('/api/settings/spotify');

        return SpotifySettings.fromJson(await resp.json());
    }

    public async setSpotify(enabled: boolean, name: string) {
        await this.sdk.put('/api/settings/spotify', {
            enabled,
            name,
        });
    }

    public async setWebcam(width: number, height: number) {
        await this.sdk.put('/api/settings/webcam', {
            width,
            height,
        });
    }

    public async setUnattended(enabled: boolean, interval: number) {
        await this.sdk.put('/api/settings/unattended', {
            enabled,
            interval,
        });
    }

    public async setFlash(powered: boolean, brightness: number) {
        await this.sdk.put('/api/settings/flash', {
            powered,
            brightness,
        });
    }

    public async getAudioDevices(): Promise<AudioDevices | null> {
        const resp = await this.sdk.get('/api/settings/audio-devices');

        return AudioDevices.fromJson(await resp.json());
    }

    public async setAudioDevices(source: number, sink: number): Promise<AudioDevices | null> {
        const resp = await this.sdk.post('/api/settings/audio-devices', {
            source_id: source,
            sink_id: sink,
        });

        return AudioDevices.fromJson(await resp.json());
    }

    public async setAudioDeviceVolume(device: AudioDevice, volume: number) {
        const resp = await this.sdk.put(`/api/webapp/settings/audio-devices/${device.id}`, { volume });

        return AudioDevices.fromJson(await resp.json());
    }

    public async getWirelessAp() {
        const resp = await this.sdk.get('/api/settings/ap');

        return InterfacesSettings.fromJson(await resp.json());
    }

    public async setWirelessAp(
        wiredInterface: string,
        wirelessInterface: string,
        enabled: boolean,
        ssid: string,
        password: string
    ) {
        const resp = await this.sdk.put('/api/settings/ap', {
            wired_interface: wiredInterface,
            wireless_interface: wirelessInterface,
            enabled,
            ssid,
            password,
        });

        return ApSettings.fromJson(await resp.json());
    }
}
