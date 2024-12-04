import { PhEvent, SDK } from './index';

export default class Settings {
    sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async setMode(mode: string) {
        await this.sdk.post(`/api/webapp/settings/mode/${mode}`);
    }

    public async setEvent(id: number): Promise<PhEvent | null> {
        const resp = await this.sdk.post(`/api/webapp/settings/event/${id}`);
        const data = await resp.json();

        return PhEvent.fromJson(data);
    }

    public async showDebug() {
        await this.sdk.post('/api/webapp/settings/debug');
    }
}
