import { KaraokeSong } from '../types/appstate';
import { SongPost } from './requests/song';
import { ApiSong } from './responses/karaoke';
import { PaginatedResponse } from './responses/paginated_responses';
import { SDK } from './sdk';

export default class Karaoke {
    sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    async search(query: string, page: number): Promise<PaginatedResponse<KaraokeSong>> {
        const resp = await this.sdk.get(
            `/api/modules/karaoke/song?page=${page}` + (query.length > 0 ? `&q=${encodeURI(query)}` : ''),
        );

        return new PaginatedResponse(await resp.json());
    }

    async post(data: SongPost) {
        const fd = new FormData();
        fd.append('title', data.title);
        fd.append('artist', data.artist);
        fd.append('format', data.format);
        fd.append('cover_type', data.cover_type)

        if (data.cover_type === 'LINK') {
            fd.append('cover_url', data.cover_data);
        } else if (data.cover_type === 'UPLOADED') {

        }

        fd.append('song', data.song[0]);

        if (data.format === 'CDG') {
            fd.append('cdg', data.cdg[0]);
        }

        const resp = await this.sdk.request(
            '/api/modules/karaoke/song',
            {
                'method': 'POST',
                'body': fd,
            },
        );

        if (resp.status != 201) {
            throw resp.text();
        }
    
        return await resp.json();
    }

    async searchSpotify(artist: string, title: string): Promise<ApiSong[]> {
        const resp = await this.sdk.post(`/api/modules/karaoke/spotify-search?q=${encodeURI(artist + ' ' + title)}`)
        return resp.json();
    }

    async rescanSongs() {
        const resp = await this.sdk.post('/api/modules/karaoke/rescan');
        if (resp.status != 200) {
            throw resp.text();
        }
    }
}