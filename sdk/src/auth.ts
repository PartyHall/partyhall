import { AuthResponse } from './models/auth';
import { SDK } from './index';

export default class Auth {
    private sdk: SDK;

    public constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    async login(username: string, password: string): Promise<AuthResponse> {
        try {
            const resp = await this.sdk.post('/api/login', {
                username,
                password,
            });

            return await resp.json();
        } catch (e: any) {
            if (e.status == 401) {
                throw e.body.message;
            }

            throw e;
        }
    }

    async loginGuest(username: string): Promise<AuthResponse> {
        try {
            const resp = await this.sdk.post('/api/guest-login', { username });

            return await resp.json();
        } catch (e: any) {
            if (e.status == 401) {
                throw e.body.message;
            }

            throw e;
        }
    }

    async refresh(rt: string): Promise<AuthResponse> {
        const resp = await this.sdk.request('/api/refresh', {
            method: 'POST',
            body: JSON.stringify({
                refresh_token: rt,
            }),
            headers: {
                'Content-Type': 'application/json',
            },
        });
        return await resp.json();
    }
}
