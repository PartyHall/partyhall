import { PhEvent, PhState, SDK } from "./index";

export default class State {
    private sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async get(): Promise<PhState> {
        const resp = await this.sdk.get('/api/state');
        const data = await resp.json();

        return PhState.fromJson(data);
    }

    /**
     * @TODO: This only handles ON/OFF. Brightness will be in the settings
     */
    public async setFlash(powered: boolean) {
        await this.sdk.put(`/api/state/flash`, { powered });
    }

    public async showDebug(): Promise<void> {
        await this.sdk.post('/api/state/debug');
    }

    public async setMode(mode: string): Promise<void> {
        await this.sdk.put(`/api/state/mode`, { mode });
    }

    public async setEvent(id: number): Promise<PhEvent | null> {
        const resp = await this.sdk.put(`/api/state/event`, {
            'event': id,
        });
        const data = await resp.json();

        return PhEvent.fromJson(data);
    }

    public async setBackdrops(albumId: number | null, selectedBackdrop: number | null) {
        const resp = await this.sdk.put(`/api/state/backdrops`, {
            backdrop_album: albumId,
            selected_backdrop: selectedBackdrop,
        });

        return await resp.json();
    }
}