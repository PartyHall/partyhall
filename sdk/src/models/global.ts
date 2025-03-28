import { BackdropAlbum } from './backdrop';
import { PhEvent } from './event';
import { PhSongSession } from './karaoke';

export type VolumeType = 'instrumental' | 'vocals';

export interface ModuleSettings {
    photobooth: {
        countdown: number;
        flash_brightness: number;
        resolution: {
            width: number;
            height: number;
        };
    };
}

export class PhKaraoke {
    current: PhSongSession | null;

    isPlaying: boolean;
    countdown: number;
    timecode: number;
    volume: number;
    volumeVocals: number;

    constructor(data: Record<string, any>) {
        this.current = data['current'] ? PhSongSession.fromJson(data['current']) : null;
        this.isPlaying = data['is_playing'];
        this.timecode = data['timecode'];
        this.countdown = data['countdown'];
        this.volume = data['volume'];
        this.volumeVocals = data['volume_vocals'];
    }

    public setTimecode(timecode: number) {
        this.timecode = timecode;

        return this;
    }

    static fromJson(data: Record<string, any>) {
        return new PhKaraoke(data);
    }
}

export class PhStatus {
    currentMode: string;
    currentEvent: PhEvent | null;
    modulesSettings: ModuleSettings;
    hardwareFlashPowered: boolean;
    guestsAllowed: boolean;

    backdropAlbum: BackdropAlbum | null;
    selectedBackdrop: number;

    karaoke: PhKaraoke;
    karaokeQueue: PhSongSession[];

    syncInProgress: boolean;
    hardwareId: string | null;
    version: string | null;
    commit: string | null;

    constructor(data: Record<string, any>) {
        this.currentMode = data['current_mode'];
        this.currentEvent = PhEvent.fromJson(data['current_event']);
        this.modulesSettings = data['modules_settings'];
        this.hardwareFlashPowered = data['hardware_flash_powered'];
        this.guestsAllowed = data['guests_allowed'];

        this.backdropAlbum = BackdropAlbum.fromJson(data['backdrop_album']);
        this.selectedBackdrop = data['selected_backdrop'];

        this.karaoke = PhKaraoke.fromJson(data['karaoke']);

        this.karaokeQueue = [];
        if (data['karaoke_queue']) {
            this.karaokeQueue = data['karaoke_queue'].map((x: Record<string, any>) => PhSongSession.fromJson(x));
        }

        this.syncInProgress = data['sync_in_progress'];
        this.hardwareId = data['hwid'];
        this.version = data['version'];
        this.commit = data['commit'];
    }

    static fromJson(data: Record<string, any>) {
        return new PhStatus(data);
    }
}
