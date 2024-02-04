import { createContext, ReactNode, useCallback, useContext, useState } from "react";
import { useSnackbar } from "./snackbar";
import { EditedEvent, KaraokeSong } from "../types/appstate";
import AdminSocketProvider from "./adminSocket";
import BoothSocketProvider from "./boothSocket";
import { EventExport } from "../types/event_export";
import getSocketMode from "../utils/socket_mode";
import { Meta } from "../types/contextualized_response";
import { SDK } from "../sdk/sdk";

//@ts-ignore
export const SOCKET_MODE_DEBUG = import.meta.env.MODE === 'development';

const KNOWN_SOCKET_MODE = ['booth', 'admin'];

const TOKEN = localStorage.getItem('token');
const REFRESH_TOKEN = localStorage.getItem('refresh_token');

const storeToken = (token: string|null, refresh: string|null) => {
    if (!token || !refresh) {
        localStorage.removeItem('token');
        localStorage.removeItem('refresh_token');

        return;
    }

    localStorage.setItem('token', token);
    localStorage.setItem('refresh_token', refresh);
}

type ApiProps = {
    api: SDK;
    socketMode: string;
};

type ApiContextProps = ApiProps & {
    login: (username: string, password: string) => Promise<void>;
    logout: () => void;
    isLoggedIn: () => boolean;
};

const defaultState: ApiProps = {
    api: new SDK(TOKEN, REFRESH_TOKEN, storeToken),
    socketMode: getSocketMode(),
};

const ApiContext = createContext<ApiContextProps>({
    ...defaultState,
    login: async (username: string, password: string) => { },
    logout: () => { },
    isLoggedIn: () => false,
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

    const setToken = useCallback((token?: string, refresh?: string) => {
        if (!token || !refresh) {
            localStorage.removeItem('token');
            localStorage.removeItem('refresh_token');

            setContext({...context, api: new SDK(null, null, storeToken)});

            return;
        }

        setContext({...context, api: new SDK(token, refresh, storeToken)});

        localStorage.setItem('token', token);
        localStorage.setItem('refresh_token', refresh);
    }, [context]);

    const login = async (username: string, password: string) => {
        const data = await context.api.auth.login(username, password);
        setToken(data.token, data.refresh_token);
    };

    const isLoggedIn = () => {
        return !!context.api && !!context.api.tokenUser;
    };

    const logout = () => {
        setToken();
    };

    return <ApiContext.Provider value={{
        ...context,
        login,
        logout,
        isLoggedIn,
    }}>
        <>
            {
                // When in booth mode, we use the booth socket provider
                context.socketMode === 'booth' &&
                <BoothSocketProvider><>{children}</></BoothSocketProvider>
            }

            {
                // When in admin mode not logged in, don't use any socketprovider
                (context.socketMode === 'admin' && !isLoggedIn()) &&
                <>{children}</>
            }

            {
                // When in admin mode logged in, we use the AdminSocketProvider
                (context.socketMode === 'admin' && isLoggedIn()) &&
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