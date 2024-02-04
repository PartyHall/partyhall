import { DateTime } from "luxon";

export type TokenUser = {
    exp: DateTime;
    iat: DateTime;
    subject: string;
    username: string;
    roles: string[];
};