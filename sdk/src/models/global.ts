import { BackdropAlbum } from './backdrop';
import { PhEvent } from './event';
import { PhSongSession } from './karaoke';

export type VolumeType = 'instrumental' | 'vocals';

export type WebcamResolution = {
    width: number;
    height: number;
};

export type UnattendedSettings = {
    enabled: boolean;
    interval: number;
};

export class PhUserSettingsPhotobooth {
    countdown: number;
    flashBrightness: number;
    resolution: WebcamResolution;
    unattended: UnattendedSettings;

    constructor(data: Record<string, any>) {
        this.countdown = data['countdown'];
        this.flashBrightness = data['flash_brightness'];
        this.resolution = data['resolution'];
        this.unattended = data['unattended'];
    }
}

export class PhWirelessAp {
    enabled: boolean;
    ssid: string;
    password: string;

    constructor(data: Record<string, any>) {
        this.enabled = data['enabled'];
        this.ssid = data['ssid'];
        this.password = data['password'];
    }
}

export class PhUserSettings {
    onboarded: boolean;
    hardwareId: string;
    photobooth: PhUserSettingsPhotobooth;
    wirelessAp: PhWirelessAp;

    constructor(data: Record<string, any>) {
        this.onboarded = data['onboarded'];
        this.hardwareId = data['hardware_id'];
        this.photobooth = new PhUserSettingsPhotobooth(data['photobooth']);
        this.wirelessAp = new PhWirelessAp(data['wireless_ap']);
    }
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

export class PhState {
    currentMode: string;
    currentEvent: PhEvent | null;
    userSettings: PhUserSettings;
    hardwareFlashPowered: boolean;
    guestsAllowed: boolean;
    adminCreated: boolean;

    backdropAlbum: BackdropAlbum | null;
    selectedBackdrop: number;

    karaoke: PhKaraoke;
    karaokeQueue: PhSongSession[];

    syncInProgress: boolean;
    version: string;
    commit: string;

    constructor(data: Record<string, any>) {
        this.currentMode = data['current_mode'];
        this.currentEvent = PhEvent.fromJson(data['current_event']);
        this.userSettings = new PhUserSettings(data['user_settings']);
        this.hardwareFlashPowered = data['hardware_flash_powered'];
        this.guestsAllowed = data['guests_allowed'];
        this.adminCreated = data['admin_created'];

        this.backdropAlbum = BackdropAlbum.fromJson(data['backdrop_album']);
        this.selectedBackdrop = data['selected_backdrop'];

        this.karaoke = PhKaraoke.fromJson(data['karaoke']);

        this.karaokeQueue = [];
        if (data['karaoke_queue']) {
            this.karaokeQueue = data['karaoke_queue'].map((x: Record<string, any>) => PhSongSession.fromJson(x));
        }

        this.syncInProgress = data['sync_in_progress'];
        this.version = data['version'];
        this.commit = data['commit'];
    }

    static fromJson(data: Record<string, any>) {
        return new PhState(data);
    }
}
