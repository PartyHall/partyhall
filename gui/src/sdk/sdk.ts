import { DateTime } from "luxon";

import Auth from "./auth";
import Karaoke from "./karaoke";
import { TokenUser } from "./responses/user";
import { parseUser } from "./utils";
import { Events } from "./events";

export type StoreToken = (token: string | null, refreshToken: string | null) => void;
export type OnExpired = () => void;

export class SDK {
    private previousTokenRefresher: number | null = null;
    private token: string | null = null;
    private refreshToken: string | null = null;
    private autorefreshTimeout: number | null = null;
    private storeToken: StoreToken = () => { };
    private onExpired: OnExpired = () => { };
    tokenUser: TokenUser | null = null;

    auth: Auth;
    events: Events;
    karaoke: Karaoke;

    constructor(token?: string | null, refreshToken?: string | null, storeToken?: StoreToken) {
        this.auth = new Auth(this);
        this.events = new Events(this);
        this.karaoke = new Karaoke(this);

        this.storeToken = storeToken || (() => { });

        if (!!token && !!refreshToken) {
            this.setToken(token, refreshToken);
        }
    }

    getToken(): string | null {
        return this.token;
    }

    async request(url: string | URL, init?: RequestInit) {
        if (!init) {
            init = {
                headers: {},
            };
        } else if (!init.headers) {
            init.headers = {};
        }

        if (this.token) {
            init.headers = {
                ...init.headers,
                'Authorization': 'Bearer ' + this.token,
            };
        }

        try {
            const resp = await fetch(url, init);

            if (this.isHttpError(resp)) {
                let body = await resp.text();
                try {
                    body = JSON.parse(body);
                } catch { }

                throw {
                    status: resp.status,
                    message: resp.statusText,
                    body,
                    resp,
                };
            }

            return resp;
        } catch (e: any) {
            if (this.onExpired && e.status == 401) {
                this.onExpired();
            }

            throw e;
        }
    }

    async get(url: string, options?: any) {
        return await this.request(url, options);
    }

    async post(url: string, data?: any, options?: any) {
        if (!options) {
            options = {};
        }

        if (data) {
            if (!options.headers) {
                options.headers = {};
            }

            options.headers['Content-Type'] = 'application/json';
            options.body = JSON.stringify(data);
        }

        options['method'] = 'POST';
        return await this.request(url, options);
    }

    async patch(url: string, data?: any, options?: any) {
        if (!options) {
            options = { headers: {} };
        }

        if (data) {
            if (!options.headers) {
                options.headers = {};
            }

            options.headers['Content-Type'] = 'application/json';
            options.body = JSON.stringify(data);
        }

        options['method'] = 'PATCH';
        return await this.request(url, options);
    }

    async put(url: string, data?: any, options?: any) {
        if (!options) {
            options = { headers: {} };
        }

        if (data) {
            if (!options.headers) {
                options.headers = {};
            }

            options.headers['Content-Type'] = 'application/json';
            options.body = JSON.stringify(data);
        }

        options['method'] = 'PUT';
        return await this.request(url, options);
    }

    async delete(url: string) {
        return await this.request(url, { method: 'DELETE' })
    }

    isHttpError(response: Response): boolean {
        return !(response.status >= 200 && response.status <= 299);
    }

    setToken(token: string | null, refreshToken: string | null) {
        this.clearRefresh();

        this.token = token;
        this.refreshToken = refreshToken;
        this.storeToken(token, refreshToken);

        this.tokenUser = parseUser(token);

        this.autoRefresh();
    }

    setStoreToken(storeToken: (token: string | null, refreshToken: string | null) => void) {
        this.storeToken = storeToken;
    }

    setOnExpired(onExpired: OnExpired) {
        this.onExpired = onExpired;
    }

    clearRefresh() {
        if (this.autorefreshTimeout !== null) {
            clearTimeout(this.autorefreshTimeout);
        }
    }

    private async autoRefresh() {
        if (!!this.previousTokenRefresher) {
            clearTimeout(this.previousTokenRefresher);
        }

        if (!this.refreshToken || !this.tokenUser) {
            return;
        }

        const now = DateTime.now();
        const exp = this.tokenUser.exp;

        const diffSeconds = exp.diff(now, 'seconds').seconds;

        if (diffSeconds < 30) {
            try {
                const data = await this.auth.refresh(this.refreshToken);
                this.setToken(data.token, data.refresh_token);
            } catch {
                this.setToken(null, null);
            }
        } else {
            this.previousTokenRefresher = setTimeout(() => this.autoRefresh(), (diffSeconds - 30) * 1000);
        }
    }
}