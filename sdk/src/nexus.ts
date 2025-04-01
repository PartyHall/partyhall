import { PhEvent, SDK } from './index';

export default class Nexus {
    private sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    async createEvent(eventId: number) {
        const resp = await this.sdk.post(`/api/nexus/create_event/${eventId}`, {});

        return PhEvent.fromJson(await resp.json());
    }

    async sync() {
        await this.sdk.post('/api/nexus/sync');
    }
}
