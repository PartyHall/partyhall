import { PhKaraoke, SDK, VolumeType } from './index';

export default class Karaoke {
    private sdk: SDK;

    public constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async songProgress(timecode: number) {
        await this.sdk.put(`/api/karaoke/timecode`, {
            timecode: timecode,
        });
    }

    public async togglePlay(status: boolean) {
        await this.sdk.put(`/api/karaoke/playing_status`, { status });
    }

    public async setVolume(type: VolumeType, volume: number): Promise<PhKaraoke> {
        const resp = await this.sdk.put(`/api/karaoke/volume`, { type, volume });
        const data = await resp.json();

        return PhKaraoke.fromJson(data);
    }
}
