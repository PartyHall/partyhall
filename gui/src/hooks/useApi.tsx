import { createContext, ReactNode, useContext, useState } from "react";
import { useSnackbar } from "./snackbar";
import { EditedEvent, KaraokeSong } from "../types/appstate";
import AdminSocketProvider from "./adminSocket";
import BoothSocketProvider from "./boothSocket";
import { EventExport } from "../types/event_export";
import getSocketMode from "../utils/socket_mode";
import { Meta } from "../types/contextualized_response";

//@ts-ignore
export const SOCKET_MODE_DEBUG = import.meta.env.MODE === 'development';

const KNOWN_SOCKET_MODE = ['booth', 'admin'];

type ApiProps = {
    socketMode: string;
    connecting: boolean;
    username: string;
    token: string|null;
    refreshToken: string|null;
};

type ApiContextProps = ApiProps & {
    login: (username: string, password: string) => Promise<void>;
    logout: () => void;

    saveEvent: (event: EditedEvent) => Promise<boolean>;
    getLastExports: (eventId: number) => Promise<EventExport[]>;
    karaokeSongSearch: (currentPage: number, search: string) => Promise<{
        results: KaraokeSong[],
        meta: Meta,
    }>;
};

const defaultState: ApiProps = {
    socketMode: getSocketMode(),
    connecting: false,
    username: localStorage.getItem('username') || '',
    token: localStorage.getItem('token'),
    refreshToken: localStorage.getItem('refreshToken'),
};

const ApiContext = createContext<ApiContextProps>({
    ...defaultState,
    login: async (username: string, password: string) => { },
    logout: () => { },
    saveEvent: async (event: EditedEvent) => false,
    getLastExports: async (eventId: number) => [],
    //@ts-ignore fuck off
    karaokeSongSearch: async () => {},
});

export default function ApiProvider({ children }: { children: ReactNode }) {
    const [context, setContext] = useState<ApiProps>(defaultState)
    const { showSnackbar } = useSnackbar();

    const changeMode = () => {
        if (!SOCKET_MODE_DEBUG) {
            return;
        }

        let idx = KNOWN_SOCKET_MODE.indexOf(context.socketMode) + 1;
        if (idx >= KNOWN_SOCKET_MODE.length) {
            idx = 0;
        }

        setContext({ ...context, socketMode: KNOWN_SOCKET_MODE[idx] })
    };

    const login = async (username: string, password: string) => {
        setContext({ ...context, connecting: true });

        const resp = await fetch(`/api/login`, {
            'method': 'POST',
            'headers': {'Content-Type': 'application/json'},
            'body': JSON.stringify({ username, password }),
        });

        let message = '';

        if (resp.status === 200) {
            const data = await resp.json();
            const token = data['token'];

            if (token) {
                localStorage.setItem('username', username);
                localStorage.setItem('token', token);
                setContext({ ...context, connecting: false, username, token });

                return;
            }
        } else if (resp.status === 401) {
            message = 'Invalid username or password'
        } else {
            message = 'Wrong response from API';
        }

        showSnackbar(message, 'error');
        setContext({ ...context, connecting: false });
    };

    const saveEvent = async (event: EditedEvent) => {
        const query: any = {
            name: event.name,
            author: event.author,
            location: event.location,
        };

        if (!!event.date) {
            query.date = Math.floor(event.date.toSeconds());
        }

        try {
            const resp = await fetch(
                '/api/admin/event' + (!!event.id ? `/${event.id}` : ''),
                {
                    method: !!event.id ? 'PUT' : 'POST',
                    headers: {
                        'Authorization': context.token ? ('Bearer ' + context.token) : '',
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(query),
                }
            );

            if (resp.status !== 201 && resp.status !== 200) {
                throw await resp.json();
            }

            showSnackbar('Event saved !', 'success');
            return true;
        } catch (e) {
            console.log(e);
            showSnackbar('Failed to save event: ' + e, 'error');
        }

        return false;
    }

    const getLastExports = async (eventId: number) => {
        const resp = await fetch(`/api/admin/event/${eventId}/export`, {
            'method': 'GET',
            'headers': { 'Authorization': context.token ? ('Bearer ' + context.token) : '' }
        });

        if (resp.status === 401) {
            showSnackbar('Session expired', 'error');
        } else {
            return await resp.json();
        }

        return [];
    };

    const karaokeSongSearch = async (currentPage: number, search: string) => {
        const resp = await fetch(
            `/api/modules/karaoke/song?page=${currentPage}` + (search.length > 0 ? `&q=${encodeURI(search)}` : ''),
            {headers: {'Authorization': context.token ? ('Bearer ' + context.token) : ''}}
        )
        return await resp.json();
    };

    const logout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('refreshToken');
        setContext({ ...context, token: null, refreshToken: null });
    }

    return <ApiContext.Provider value={{
        ...context,
        login,
        logout,
        saveEvent,
        getLastExports,
        karaokeSongSearch,
    }}>
        <>
            {
                // When in booth mode, we use the booth socket provider
                context.socketMode === 'booth' &&
                <BoothSocketProvider><>{children}</></BoothSocketProvider>
            }

            {
                // When in admin mode not logged in, don't use any socketprovider
                (context.socketMode === 'admin' && context.token === null) &&
                <>{children}</>
            }

            {
                // When in admin mode logged in, we use the AdminSocketProvider
                (context.socketMode === 'admin' && context.token !== null) &&
                <AdminSocketProvider><>{children}</></AdminSocketProvider>
            }

            {
                SOCKET_MODE_DEBUG &&
                <div className="debug absbr" onClick={changeMode}>
                    <p>Socket mode: {context.socketMode}</p>
                </div>
            }
        </>
    </ApiContext.Provider>
}

export const useApi = () => {
    return useContext<ApiContextProps>(ApiContext);
};