import Auth from './auth';
import { DateTime } from 'luxon';
import Events from './events';
import Global from './global';
import Karaoke from './karaoke';
import Nexus from './nexus';
import { PhTokenUser } from './models/user';
import Photobooth from './photobooth';
import { SdkError } from './models/sdk_error';
import Settings from './settings';
import Backdrop from './backdrops';

export type StoreToken = (token: string | null, refreshToken: string | null) => void;
export type OnExpired = () => void;

export class SDK {
    public token: string | null;
    public refreshToken: string | null;

    private autorefreshTimeout: number | null = null;
    private storeToken: StoreToken = () => {};
    private onExpired: OnExpired = () => {};

    public tokenUser: PhTokenUser | null = null;

    public auth: Auth;
    public backdrops: Backdrop;
    public global: Global;
    public events: Events;
    public photobooth: Photobooth;
    public karaoke: Karaoke;
    public nexus: Nexus;
    public settings: Settings;

    public constructor(token: string | null, refreshToken: string | null, storeToken?: StoreToken) {
        this.token = token;
        this.refreshToken = refreshToken;

        this.auth = new Auth(this);
        this.backdrops = new Backdrop(this);
        this.global = new Global(this);
        this.events = new Events(this);
        this.photobooth = new Photobooth(this);
        this.karaoke = new Karaoke(this);
        this.nexus = new Nexus(this);
        this.settings = new Settings(this);

        this.storeToken = storeToken || (() => {});

        if (token) {
            this.setToken(token, refreshToken);
        }
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
                Authorization: 'Bearer ' + this.token,
            };
        }

        try {
            const resp = await fetch(url, init);

            if (this.isHttpError(resp)) {
                let body = await resp.text();
                try {
                    body = JSON.parse(body);
                } catch {
                    /* empty */
                }

                throw new SdkError(resp.status, body);
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
        return await this.request(url, { method: 'DELETE' });
    }

    isHttpError(response: Response): boolean {
        return !(response.status >= 200 && response.status <= 299);
    }

    setToken(token: string | null, refreshToken: string | null) {
        this.clearRefresh();

        this.token = token;
        this.refreshToken = refreshToken;
        this.storeToken(token, refreshToken);

        this.tokenUser = PhTokenUser.fromToken(token);

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
        if (!this.refreshToken || !this.tokenUser) {
            return;
        }

        const now = DateTime.now();

        const diffSeconds = this.tokenUser.expiresAt.diff(now, 'seconds').seconds;

        if (diffSeconds < 30) {
            try {
                const data = await this.auth.refresh(this.refreshToken);
                this.setToken(data.token, data.refresh_token);
            } catch {
                this.setToken(null, null);
            }
        } else {
            setTimeout(() => this.autoRefresh(), (diffSeconds - 30) * 1000);
        }
    }
}

export * from './models/index';
