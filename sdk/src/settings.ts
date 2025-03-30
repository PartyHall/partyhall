import { AudioDevice, AudioDevices } from './models/audio';
import { SDK } from './index';

export default class Settings {
    sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
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
        const resp = await this.sdk.post(`/api/webapp/settings/audio-devices/${device.id}/volume`, { volume });

        return AudioDevices.fromJson(await resp.json());
    }
}
