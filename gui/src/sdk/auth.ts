import { AuthResponse } from './responses/auth';
import { SDK } from './sdk';

export default class Auth {
    sdk: SDK;

    constructor(sdk: SDK) {
        this.sdk = sdk;
    }

    async login(username: string, password: string): Promise<AuthResponse> {
        try {
            const resp = await this.sdk.post('/api/login', {
                username: username,
                password: password,
            });

            if (resp.status !== 200) {
                throw 'Invalid username/password';
            }

            return await resp.json();
        } catch {
            throw 'Invalid username/password';
        }
    }

    async loginAsGuest(username: string): Promise<AuthResponse> {
        const resp = await this.sdk.post('/api/login-guest', {
            username: username,
        });

        return await resp.json();
    }

    async refresh(refresh_token: string|null): Promise<AuthResponse> {
        // If we are a guest
        if (this.sdk.tokenUser?.subject == '0' && this.sdk.guestUsername) {
            return await this.loginAsGuest(this.sdk.guestUsername);
        }

        const resp = await this.sdk.request('/api/refresh', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ refresh_token }),
        });
        return await resp.json();
    }
}