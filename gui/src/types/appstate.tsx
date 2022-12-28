type Event = {
    id: number;
    name: string;
    author: string;
    date: number;
    location?: string|null;
};

type appstate = {
    hwid: string;
    token: string;
    current_event: Event|null;
};

type Photobooth = {
    hardware_flash: boolean;
    webcam_resolution: {
        width: number;
        height: number;
    };
};

export type AppState = {
    app_state: appstate;
    photobooth: Photobooth;
    debug: boolean;
    current_mode: string|null;

    ip_addresses: {
        [key: string]: string[];
    };

    known_events: Event[];
    known_modes: string[];

    partyhall_version: string;
    partyhall_commit: string;
};