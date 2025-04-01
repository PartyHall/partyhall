import { DateTime } from 'luxon';

// @TODO: Camel case
export class PhSong {
    nexus_id: string;
    title: string;
    artist: string;
    spotify_id: string | null;
    format: string;
    hotspot: string | null;
    duration: number;

    has_cover: boolean;
    has_vocals: boolean;
    has_combined: boolean;

    constructor(data: Record<string, any>) {
        this.nexus_id = data['nexus_id'];
        this.title = data['title'];
        this.artist = data['artist'];
        this.spotify_id = data['spotify_id'] ?? null;
        this.format = data['format'];
        this.hotspot = data['hotspot'] ?? null;
        this.duration = data['duration'] ?? null;
        this.has_cover = data['has_cover'];
        this.has_vocals = data['has_vocals'];
        this.has_combined = data['has_combined'];
    }

    private getUrl(filename: string) {
        return `/api/songs/${this.nexus_id}/file/${filename}`;
    }

    getInstrumentalUrl() {
        if (this.format === 'cdg') {
            return this.getUrl('instrumental.mp3');
        }

        return this.getUrl('instrumental.webm');
    }

    getLyricsUrl() {
        if (this.format !== 'cdg') {
            throw 'video songs do not have a lyrics URL';
        }

        return this.getUrl('lyrics.cdg');
    }

    getVocalsUrl() {
        return this.getUrl('vocals.mp3');
    }

    getCombinedUrl() {
        return this.getUrl('combined.mp3');
    }

    static fromJson(data: Record<string, any> | null) {
        if (!data) {
            return null;
        }

        return new PhSong(data);
    }
}

export class PhSongSession {
    id: number;
    event_id: string;
    title: string;
    artist: string;
    sung_by: string;
    song: PhSong;
    added_at: DateTime | null;
    started_at: DateTime | null;
    ended_at: DateTime | null;
    cancelled_at: DateTime | null;

    constructor(data: Record<string, any>) {
        this.id = data['id'];
        this.event_id = data['event_id'];
        this.title = data['title'];
        this.artist = data['artist'];
        this.sung_by = data['sung_by'];
        this.added_at = data['added_at'] ? DateTime.fromISO(data['added_at']) : null;
        this.started_at = data['started_at'] ? DateTime.fromISO(data['started_at']) : null;
        this.ended_at = data['ended_at'] ? DateTime.fromISO(data['ended_at']) : null;
        this.cancelled_at = data['cancelled_at'] ? DateTime.fromISO(data['cancelled_at']) : null;

        const song = PhSong.fromJson(data['song']);
        if (!song) {
            throw 'no song parsed';
        }

        this.song = song;
    }

    static fromJson(data: Record<string, any> | null) {
        if (!data) {
            return null;
        }

        return new PhSongSession(data);
    }
}
