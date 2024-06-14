import { DateTime } from "luxon";

export type SoundDevice = {
    index: number;
    state: string;
    name: string;
    description: string;
    driver: string;
    mute: boolean;
    volume: {
        db: string;
        value: number;
        value_percent: string;
    }
};

export type Event = {
    id: number;
    name: string;
    author: string;
    date: number;
    location?: string|null;

    amt_images_handtaken?: number|null;
    amt_images_unattended?: number|null;
};

export type EditedEvent = {
    id?: number|'';
    name?: string|null;
    author?: string|null;
    date?: DateTime;
    location?: string|null;
}

type appstate = {
    hwid: string;
    token: string;
    current_event: Event|null;
};

type PhotoboothModule = {
    hardware_flash: boolean;

    default_timer: number;
    unattended_interval: number;

    webcam_resolution: {
        width: number;
        height: number;
    };
};

export type KaraokeSong = {
    id: number;
    uuid: string;
    spotify_id?: string;
    artist: string;
    title: string;
    hotspot?: string;
    format: string;

    has_cover: boolean;
    has_vocals: boolean;
    has_full: boolean;
};

export type KaraokeSongSession = {
    id: number;
    song: KaraokeSong;
    sung_by: string;
}

type KaraokeModule = {
    currentSong: KaraokeSongSession;
    queue: KaraokeSongSession[];
    started: boolean;
    preplayTimer: number;
};

export type AppState = {
    debug: boolean;
    app_state: appstate;
    current_mode: string;
    guests_allowed: boolean;

    ip_addresses: { [key: string]: string[]; };
    known_events: Event[];
    known_modes: string[];

    pulseaudio_selected: SoundDevice|null;
    pulseaudio_devices?: SoundDevice[];

    modules: {
        photobooth: PhotoboothModule;
        karaoke: KaraokeModule;
    };

    partyhall_version: string;
    partyhall_commit: string;
    language: string;
};