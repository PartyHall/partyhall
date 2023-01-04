import { DateTime } from "luxon";

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

export type AppState = {
    app_state: appstate;
    current_mode: string;
    ip_addresses: { [key: string]: string[]; };
    known_events: Event[];
    known_modes: string[];

    modules: {
        photobooth: PhotoboothModule;
    };

    partyhall_version: string;
    partyhall_commit: string;
};