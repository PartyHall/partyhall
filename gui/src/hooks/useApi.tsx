import { createContext, ReactNode, useContext, useState } from "react";
import { useSnackbar } from "./snackbar";
import { EditedEvent } from "../types/appstate";
import AdminSocketProvider from "./adminSocket";
import BoothSocketProvider from "./boothSocket";
import { EventExport } from "../types/event_export";

const SOCKET_MODE_DEBUG = true;

const KNOWN_SOCKET_MODE = ['booth', 'admin'];

type ApiProps = {
    socketMode: string;
    connecting: boolean;
    password: string | null;
};

type ApiContextProps = ApiProps & {
    login: (password: string) => Promise<void>;
    logout: () => void;

    saveEvent: (event: EditedEvent) => Promise<void>;
    getLastExports: (eventId: number) => Promise<EventExport | null>;
};

const defaultState: ApiProps = {
    socketMode: 'booth',
    connecting: false,
    password: localStorage.getItem('password'),
};

const ApiContext = createContext<ApiContextProps>({
    ...defaultState,
    login: async (password: string) => { },
    logout: () => { },
    saveEvent: async (event: EditedEvent) => { },
    getLastExports: async (eventId: number) => null,
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

    const login = async (password: string) => {
        setContext({ ...context, connecting: true });

        const resp = await fetch(`/api/admin/login`, {
            'method': 'POST',
            'headers': {
                'Authorization': password,
            }
        });

        let message = '';

        if (resp.status === 401) {
            message = 'Invalid password'
        } else {
            const data = await resp.text();

            if (data === 'yes') {
                localStorage.setItem('password', password);
                setContext({ ...context, connecting: false, password });
                return;
            }

            message = 'Wrong response from API';
        }

        showSnackbar(message, 'error');
        setContext({ ...context, connecting: false });
    };

    const saveEvent = async (event: EditedEvent) => {
        const query = {
            name: event.name,
            author: event.author,
            location: event.location,
        };

        if (!!event.date) {
            //@ts-ignore
            query.date = Math.floor(event.date.toSeconds());
        }

        try {
            const resp = await fetch(
                '/api/admin/event' + (!!event.id ? `/${event.id}` : ''),
                {
                    method: !!event.id ? 'PUT' : 'POST',
                    headers: {
                        'Authorization': context.password ?? '',
                        'Content-Type': 'application/json',
                    },
                    //@ts-ignore
                    body: JSON.stringify(query),
                }
            );

            if (resp.status !== 200) {
                throw await resp.json();
            }

            showSnackbar('Event saved !', 'success');
        } catch (e) {
            showSnackbar('Failed to save event: ' + e, 'error');
        }
    }

    const getLastExports = async (eventId: number) => {
        const resp = await fetch(`/api/admin/exports/${eventId}`, {
            'method': 'GET',
            'headers': { 'Authorization': context.password ?? '' }
        });

        if (resp.status === 401) {
            showSnackbar('Session expired', 'error');
        } else {
            return await resp.json();
        }

        return [];
    };

    const logout = () => {
        localStorage.removeItem('password');
        setContext({ ...context, password: null });
    }

    return <ApiContext.Provider value={{
        ...context,
        login,
        logout,
        saveEvent,
        getLastExports,
    }}>
        <>
            {
                // When in booth mode, we use the booth socket provider
                context.socketMode === 'booth' &&
                <BoothSocketProvider><>{children}</></BoothSocketProvider>
            }

            {
                // When in admin mode not logged in, don't use any socketprovider
                (context.socketMode === 'admin' && context.password === null) &&
                <>{children}</>
            }

            {
                // When in admin mode logged in, we use the AdminSocketProvider
                (context.socketMode === 'admin' && context.password !== null) &&
                <AdminSocketProvider><>{children}</></AdminSocketProvider>
            }

            {
                SOCKET_MODE_DEBUG &&
                <div className="debug absbr" onClick={changeMode}>
                    Socket mode: {context.socketMode}
                </div>
            }
        </>
    </ApiContext.Provider>
}

export const useApi = () => {
    return useContext<ApiContextProps>(ApiContext);
};