import { PhEvent, SDK } from './index';

export default class Nexus {
    private sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    async createEvent(eventId: number) {
        const resp = await this.sdk.post(`/api/webapp/nexus/events/${eventId}`, {});
        const data = await resp.json();

        return PhEvent.fromJson(data);
    }
}
