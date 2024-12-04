import { DateTime } from 'luxon';
import { Log } from './models/log';
import { PhStatus } from './models/global';
import { SDK } from './index';

export default class Global {
    private sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async getStatus(): Promise<PhStatus> {
        const resp = await this.sdk.get('/api/status');
        const data = await resp.json();

        return PhStatus.fromJson(data);
    }

    public async forceSync() {
        await this.sdk.post('/api/webapp/settings/force-sync');
    }

    public async getLogs(): Promise<Log[]> {
        const resp = await this.sdk.get('/api/webapp/logs');
        const data = await resp.json();

        /** @TODO make pagination */

        return data.results.map((x: any) => ({
            ...x,
            timestamp: DateTime.fromISO(x.timestamp).toFormat(
                'yyyy-MM-dd HH:mm:ss'
            ),
        }));
    }
}
