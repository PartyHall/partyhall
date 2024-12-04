import { PhEvent, SDK } from './index';
import { Collection } from './models/collection';

export default class Global {
    private sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async get(id: number | string): Promise<PhEvent> {
        const resp = await this.sdk.get(`/api/webapp/events/${id}`);
        const data = await resp.json();

        const event = PhEvent.fromJson(data);
        if (!event) {
            throw 'Failed to parse event';
        }

        return event;
    }

    public async getCollection(
        page: number | null
    ): Promise<Collection<PhEvent>> {
        if (!page) {
            page = 1;
        }

        const resp = await this.sdk.get(`/api/webapp/events?page=${page}`);
        const data = await resp.json();

        const events = Collection.fromJson(data, (x) => {
            const event = PhEvent.fromJson(x);
            if (!event) {
                throw 'Failed to parse event';
            }

            return event;
        });
        if (!events) {
            throw 'Failed to parse event';
        }

        return events;
    }

    public async create(event: PhEvent) {
        const resp = await this.sdk.post('/api/webapp/events', event.asJson());
        const data = await resp.json();

        return PhEvent.fromJson(data);
    }

    public async update(event: PhEvent) {
        const resp = await this.sdk.put(
            `/api/webapp/events/${event.id}`,
            event.asJson()
        );
        const data = await resp.json();

        return PhEvent.fromJson(data);
    }

    public async delete(id: number | string) {
        await this.sdk.delete(`/api/webapp/events/${id}`);
    }
}
