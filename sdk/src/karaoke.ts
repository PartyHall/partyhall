import { Collection, PhKaraoke, PhSong, PhSongSession, SDK, VolumeType } from './index';

export default class Karaoke {
    private sdk: SDK;

    public constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async getCollection(
        page: number | null,
        search?: string | null,
        formats: string[] = [],
        hasVocals: boolean | null = null
    ): Promise<Collection<PhSong>> {
        if (!page) {
            page = 1;
        }

        const query = new URLSearchParams();
        query.set('page', `${page}`);
        if (search) {
            query.set('search', search);
        }

        if (formats.length > 0) {
            query.set('formats', formats.join(','));
        }

        if (hasVocals !== null) {
            query.set('has_vocals', hasVocals ? 'true' : 'false');
        }

        const resp = await this.sdk.get(`/api/webapp/songs?${query.toString()}`);
        const data = await resp.json();

        const events = Collection.fromJson(data, (x) => {
            const event = PhSong.fromJson(x);
            if (!event) {
                throw 'Failed to parse event';
            }

            return event;
        });
        if (!events) {
            throw 'Failed to parse event';
        }

        return events;
    }

    public async addToQueue(
        songId: string,
        singerName: string,
        directPlay: boolean = false
    ): Promise<PhSongSession | null> {
        const resp = await this.sdk.post(`/api/webapp/session`, {
            song_id: songId,
            display_name: singerName,
            direct_play: directPlay,
        });

        const data = await resp.json();

        return PhSongSession.fromJson(data);
    }

    public async cancelSession(sessionId: number) {
        await this.sdk.delete(`/api/webapp/session/${sessionId}`);
    }

    public async moveInQueue(sessionId: number, direction: 'up' | 'down'): Promise<PhSongSession[]> {
        const resp = await this.sdk.post(`/api/webapp/session/${sessionId}/move/${direction}`);
        const data = await resp.json();

        return data.map((x: any) => PhSongSession.fromJson(x)).filter((x: any) => !!x);
    }

    public async directPlay(sessionId: number) {
        await this.sdk.post(`/api/webapp/session/${sessionId}/start`);
    }

    public async songEnded(sessionId: number) {
        await this.sdk.post(`/api/webapp/session/${sessionId}/ended`);
    }

    public async songProgress(timecode: number) {
        await this.sdk.post(`/api/webapp/karaoke/timecode`, {
            timecode: timecode,
        });
    }

    public async togglePlay(playing: boolean) {
        await this.sdk.post(`/api/webapp/karaoke/playing-status/${playing}`);
    }

    public async setVolume(type: VolumeType, volume: number): Promise<PhKaraoke> {
        const resp = await this.sdk.post(`/api/webapp/karaoke/set-volume/${type}/${volume}`);
        const data = await resp.json();

        return PhKaraoke.fromJson(data);
    }
}
