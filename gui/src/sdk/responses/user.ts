import { DateTime } from "luxon";

export type TokenUser = {
    exp: DateTime;
    iat: DateTime;
    subject: string;
    username: string;
    name: string;
    roles: string[];
};