import { Collection, PhSong, SDK } from './index';

export default class Songs {
    private sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async getCollection(
        page: number | null,
        search?: string | null,
        formats: string[] = [],
        hasVocals: boolean | null = null
    ): Promise<Collection<PhSong>> {
        if (!page) {
            page = 1;
        }

        const query = new URLSearchParams();
        query.set('page', `${page}`);
        if (search) {
            query.set('search', search);
        }

        if (formats.length > 0) {
            query.set('formats', formats.join(','));
        }

        if (hasVocals !== null) {
            query.set('has_vocals', hasVocals ? 'true' : 'false');
        }

        const resp = await this.sdk.get(`/api/songs?${query.toString()}`);
        const data = await resp.json();

        const events = Collection.fromJson(data, (x) => {
            const event = PhSong.fromJson(x);
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

    public getCoverUrl(id: string): string {
        return `/api/songs/${id}/cover`;
    }
}
