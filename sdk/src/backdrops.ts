import { Collection, SDK } from './index';
import { BackdropAlbum } from './models/backdrop';

export default class Backdrop {
    private sdk: SDK;

    public constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async getAlbumCollection(page: number | null, search?: string | null): Promise<Collection<BackdropAlbum>> {
        if (!page) {
            page = 1;
        }

        const query = new URLSearchParams();
        query.set('page', `${page}`);
        if (search) {
            query.set('search', search);
        }

        const resp = await this.sdk.get(`/api/webapp/backdrops?${query.toString()}`);
        const data = await resp.json();

        const events = Collection.fromJson(data, (x) => {
            const event = BackdropAlbum.fromJson(x);
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
}
