import { AuthResponse } from './responses/auth';
import { SDK } from './sdk';

export default class Auth {
    sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    async login(username: string, password: string): Promise<AuthResponse> {
        const resp = await this.sdk.post('/api/login', {
            username: username,
            password: password,
        });

        return await resp.json();
    }

    async refresh(refresh_token: string): Promise<AuthResponse> {
        const resp = await this.sdk.request('/api/refresh', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ refresh_token }),
        });
        return await resp.json();
    }
}