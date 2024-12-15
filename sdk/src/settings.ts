import { AudioDevice, AudioDevices } from './models/audio';
import { PhEvent, SDK } from './index';

export default class Settings {
    sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async setMode(mode: string): Promise<void> {
        await this.sdk.post(`/api/webapp/settings/mode/${mode}`);
    }

    public async setEvent(id: number): Promise<PhEvent | null> {
        const resp = await this.sdk.post(`/api/webapp/settings/event/${id}`);
        const data = await resp.json();

        return PhEvent.fromJson(data);
    }

    public async showDebug(): Promise<void> {
        await this.sdk.post('/api/webapp/settings/debug');
    }

    public async getAudioDevices(): Promise<AudioDevices | null> {
        const resp = await this.sdk.get('/api/webapp/settings/audio-devices');
        const data = await resp.json();

        return AudioDevices.fromJson(data);
    }

    public async setAudioDevices(source: number, sink: number): Promise<AudioDevices | null> {
        const resp = await this.sdk.post('/api/webapp/settings/audio-devices', {
            source_id: source,
            sink_id: sink,
        });
        const data = await resp.json();

        return AudioDevices.fromJson(data);
    }

    public async setAudioDeviceVolume(device: AudioDevice, volume: number) {
        const resp = await this.sdk.post(`/api/webapp/settings/audio-devices/${device.id}/volume`, { volume });
        const data = await resp.json();

        return AudioDevices.fromJson(data);
    }
}
