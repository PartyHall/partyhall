import { SDK } from "./sdk";
import { EditedEvent } from "../types/appstate";

export class Events {
    sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    async save(event: EditedEvent) {
        const query: any = {
            name: event.name,
            author: event.author,
            location: event.location,
        };

        if (!!event.date) {
            query.date = Math.floor(event.date.toSeconds());
        }

        try {
            const resp = await this.sdk.request(
                '/api/admin/event' + (!!event.id ? `/${event.id}` : ""),
                {
                    'method': !!event.id ? 'PUT' : 'POST',
                    'headers': {'Content-Type': 'application/json'},
                    'body': JSON.stringify(query),
                }
            );

            if (resp.status !== 201 && resp.status !== 200) {
                throw await resp.json();
            }
        } catch (e) {
            console.log(e);
            throw e;
        }
    }

    async getLastExports(eventId: number) {
        return await (await this.sdk.get(`/api/admin/event/${eventId}/export`)).json();
    };

    async downloadExport(exportId: number) {
        return this.sdk.request(
            `/api/admin/event/${exportId}/export/download`,
            { method: 'GET' },
        );
    }
}