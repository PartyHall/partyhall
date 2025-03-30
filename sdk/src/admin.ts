import { DateTime } from "luxon";
import { Log } from "./models/log";
import { SDK } from "./index";

export default class Admin {
    private sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async getLogs(): Promise<Log[]> {
        const resp = await this.sdk.get('/api/admin/logs');
        const data = await resp.json();

        /** @TODO make pagination */

        return data.results.map((x: any) => ({
            ...x,
            timestamp: DateTime.fromISO(x.timestamp).toFormat('yyyy-MM-dd HH:mm:ss'),
        }));
    }

    public async shutdown(): Promise<void> {
        await this.sdk.post(`/api/admin/shutdown`);
    }
}