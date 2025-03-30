import { PhSongSession, SDK } from "./index";

export default class SongSessions {
    private sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async addToQueue(
        songId: string,
        singerName: string,
        directPlay: boolean = false
    ): Promise<PhSongSession | null> {
        const resp = await this.sdk.post(`/api/song_sessions`, {
            song_id: songId,
            display_name: singerName,
            direct_play: directPlay,
        });

        const data = await resp.json();

        return PhSongSession.fromJson(data);
    }

    public async cancelSession(sessionId: number) {
        await this.sdk.delete(`/api/song_sessions/${sessionId}`);
    }

    public async moveInQueue(sessionId: number, direction: 'up' | 'down'): Promise<PhSongSession[]> {
        const resp = await this.sdk.post(`/api/song_sessions/${sessionId}/move/${direction}`);
        const data = await resp.json();

        return data.map((x: any) => PhSongSession.fromJson(x)).filter((x: any) => !!x);
    }

    public async directPlay(sessionId: number) {
        await this.sdk.post(`/api/song_sessions/${sessionId}/start`);
    }

    public async songEnded(sessionId: number) {
        await this.sdk.post(`/api/song_sessions/${sessionId}/ended`);
    }
}