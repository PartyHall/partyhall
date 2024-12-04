import { DateTime } from 'luxon';

export class PhTokenUser {
    issuedAt: DateTime;
    expiresAt: DateTime;

    id: string;
    username: string;
    name: string | null;
    roles: string[];

    constructor(
        issuedAt: DateTime,
        expiresAt: DateTime,
        id: string,
        username: string,
        name: string,
        roles: string[]
    ) {
        this.issuedAt = issuedAt;
        this.expiresAt = expiresAt;
        this.id = id;
        this.username = username;
        this.name = name;
        this.roles = roles;
    }

    public hasRole(role: string): boolean {
        return this.roles.includes(role);
    }

    static fromJson(json: Record<string, any> | null) {
        if (!json) {
            return null;
        }

        return new PhTokenUser(
            DateTime.fromMillis(json.iat * 1000),
            DateTime.fromMillis(json.exp * 1000),
            json.id,
            json.username,
            json.name,
            json.roles
        );
    }

    static fromToken(token: string | null) {
        if (!token) {
            return null;
        }

        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const data = JSON.parse(
            decodeURIComponent(
                window
                    .atob(base64)
                    .split('')
                    .map(function (c) {
                        return (
                            '%' +
                            ('00' + c.charCodeAt(0).toString(16)).slice(-2)
                        );
                    })
                    .join('')
            )
        );

        return PhTokenUser.fromJson(data);
    }
}
