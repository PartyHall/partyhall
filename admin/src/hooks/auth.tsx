import { PhEvent, PhKaraoke, PhSongSession, SDK } from '@partyhall/sdk';
import { ReactNode, createContext, useCallback, useContext, useState } from 'react';
import Cookies from 'js-cookie';
import { DateTime } from 'luxon';
import MercureProvider from './mercure';
import useAsyncEffect from 'use-async-effect';
import { useSettings } from './settings';

const TOKEN = localStorage.getItem('token');
const REFRESH_TOKEN = localStorage.getItem('refresh_token');

const storeToken = (token: string | null, refresh: string | null) => {
    if (!token || !refresh) {
        localStorage.removeItem('token');
        localStorage.removeItem('refresh_token');

        Cookies.remove('mercureAuthorization');

        return;
    }

    localStorage.setItem('token', token);
    localStorage.setItem('refresh_token', refresh);

    Cookies.set('mercureAuthorization', token);
};

type AuthProps = {
    loaded: boolean;
    api: SDK;

    displayName: string | null;

    mode: string | null;
    event: PhEvent | null;
    time: DateTime | null;

    hardwareFlash: {
        powered: boolean;
        brightness: number;
    };

    karaoke: PhKaraoke | null;
    karaokeQueue: PhSongSession[];

    syncInProgress: boolean;
    hwid: string | null;
    version: string | null;
    commit: string | null;
};

type AuthContextProps = AuthProps & {
    login: (username: string, password: string) => Promise<void>;
    loginGuest: (username: string) => Promise<void>;
    setToken: (token: string, refresh: string) => void;
    isLoggedIn: () => boolean;
    logout: () => void;

    setMode: (mode: string) => void;
    setEvent: (evt: PhEvent) => void;
    setHardwareFlash: (powered: boolean, brightness: number) => void;
    setKaraoke: (karaoke: PhKaraoke) => void;
    setTimecode: (timecode: number) => void;
    setKaraokeQueue: (queue: PhSongSession[]) => void;
    setSyncInProgress: (syncInProgress: boolean) => void;

    setDisplayName: (displayName: string) => void;
};

const defaultProps: AuthProps = {
    loaded: false,
    api: new SDK(TOKEN, REFRESH_TOKEN, storeToken),
    displayName: localStorage.getItem('PREVIOUS_DISPLAY_NAME'),

    mode: 'photobooth',
    event: null,
    time: null,

    hardwareFlash: {
        powered: false,
        brightness: 100,
    },

    karaoke: null,
    karaokeQueue: [],

    syncInProgress: false,
    hwid: null,
    version: null,
    commit: null,
};

const AuthContext = createContext<AuthContextProps>({
    ...defaultProps,
    login: async () => {},
    loginGuest: async () => {},
    setToken: () => {},
    isLoggedIn: () => false,
    logout: () => {},

    setMode: () => {},
    setEvent: () => {},
    setHardwareFlash: () => {},
    setKaraoke: () => {},
    setTimecode: () => {},
    setKaraokeQueue: () => {},
    setSyncInProgress: () => {},

    setDisplayName: () => {},
});

export default function AuthProvider({ children }: { children: ReactNode }) {
    const { topics, hwflash_powered, modules_settings } = useSettings();
    const [context, setContext] = useState<AuthProps>({
        ...defaultProps,
        hardwareFlash: {
            powered: hwflash_powered,
            brightness: modules_settings.photobooth.flash_brightness,
        },
    });

    const login = async (username: string, password: string) => {
        const data = await context.api.auth.login(username, password);
        setToken(data.token, data.refresh_token);
    };

    const loginGuest = async (username: string) => {
        const data = await context.api.auth.loginGuest(username);
        setToken(data.token, data.refresh_token);
    };

    const setToken = useCallback(
        (token?: string, refresh?: string) => {
            if (!token || !refresh) {
                localStorage.removeItem('token');
                localStorage.removeItem('refresh_token');
                Cookies.remove('mercureAuthorization');

                setContext((oldCtx) => ({
                    ...oldCtx,
                    loaded: true,
                    api: new SDK(null, null, storeToken),
                }));

                return;
            }

            const api = new SDK(token, refresh, storeToken);

            setContext((oldCtx) => ({
                ...oldCtx,
                loaded: true,
                api,
                displayName: localStorage.getItem('PREVIOUS_DISPLAY_NAME') || api.tokenUser?.name || null,
            }));

            localStorage.setItem('token', token);
            localStorage.setItem('refresh_token', refresh);
            Cookies.set('mercureAuthorization', token);
        },
        [context]
    );

    const isLoggedIn = () => !!context.api.tokenUser;

    const logout = () => setToken();

    const setEvent = (event: PhEvent) => setContext((oldCtx) => ({ ...oldCtx, event }));

    const setMode = (mode: string) => setContext((oldCtx) => ({ ...oldCtx, mode }));

    const setKaraoke = (karaoke: PhKaraoke) => setContext((oldCtx) => ({ ...oldCtx, karaoke }));

    const setKaraokeQueue = (queue: PhSongSession[]) => setContext((oldCtx) => ({ ...oldCtx, karaokeQueue: queue }));

    const setTimecode = (timecode: number) =>
        setContext((oldCtx) => {
            if (!oldCtx.karaoke) {
                return oldCtx;
            }

            return { ...oldCtx, karaoke: oldCtx.karaoke.setTimecode(timecode) };
        });

    const setSyncInProgress = (syncInProgress: boolean) => setContext((oldCtx) => ({ ...oldCtx, syncInProgress }));

    const setDisplayName = (displayName: string) => {
        localStorage.setItem('PREVIOUS_DISPLAY_NAME', displayName);
        setContext((oldCtx) => ({ ...oldCtx, displayName: displayName }));
    };

    const setHardwareFlash = (powered: boolean, brightness: number) =>
        setContext((oldCtx) => ({
            ...oldCtx,
            hardwareFlash: { powered, brightness },
        }));

    useAsyncEffect(async () => {
        context.api.setOnExpired(() => setToken());

        if (!isLoggedIn()) {
            return;
        }

        const status = await context.api.global.getStatus();

        setContext((oldCtx) => ({
            ...oldCtx,
            mode: status.currentMode,
            event: status.currentEvent,
            karaoke: status.karaoke,
            karaokeQueue: status.karaokeQueue,
            syncInProgress: status.syncInProgress,
        }));
    }, [context.api]);

    return (
        <AuthContext.Provider
            value={{
                ...context,
                login,
                loginGuest,
                setToken,
                isLoggedIn,
                logout,
                setEvent,
                setHardwareFlash,
                setMode,
                setKaraoke,
                setTimecode,
                setKaraokeQueue,
                setDisplayName,
                setSyncInProgress,
            }}
        >
            <MercureProvider
                sdk={context.api}
                url={`${window.location.protocol}//${window.location.host}/.well-known/mercure`}
                topics={topics}
            >
                {children}
            </MercureProvider>
        </AuthContext.Provider>
    );
}

export function useAuth() {
    return useContext<AuthContextProps>(AuthContext);
}
