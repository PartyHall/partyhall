import { DateTime } from "luxon";
import { TokenUser } from "./responses/user";

export const parseUser = (token: string|null): TokenUser|null => {
    if (!token) {
        return null;
    }

    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const data = JSON.parse(decodeURIComponent(window.atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join('')));
    
    for (const [key, value] of Object.entries(data)) {
        if (['iat', 'exp'].includes(key)) {
            //@ts-ignore
            data[key] = DateTime.fromMillis(value * 1000);
        }
    }
    
    return data;
};