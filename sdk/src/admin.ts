import { SDK, User } from './index';
import { DateTime } from 'luxon';
import { Log } from './models/log';

export default class Admin {
    private sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    public async createAdmin(displayName: string | null, username: string, password: string): Promise<User | null> {
        const resp = await this.sdk.post('/api/admin/create-admin', {
            name: displayName,
            username,
            password,
        });

        return User.fromJson(await resp.json());
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
