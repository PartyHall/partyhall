import { BackdropAlbum, PhEvent, PhKaraoke, PhSongSession, SDK } from '@partyhall/sdk';
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

    hardwareFlashPowered: boolean;

    backdropAlbum: BackdropAlbum | null;
    selectedBackdrop: number;

    karaoke: PhKaraoke | null;
    karaokeQueue: PhSongSession[];

    syncInProgress: boolean;
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
    setKaraoke: (karaoke: PhKaraoke) => void;
    setTimecode: (timecode: number) => void;
    setKaraokeQueue: (queue: PhSongSession[]) => void;
    setSyncInProgress: (syncInProgress: boolean) => void;
    setBackdrops: (backdropAlbum: BackdropAlbum | null, selectedBackdrop: number) => void;
    setHardwareFlashPowered: (powered: boolean) => void;

    setDisplayName: (displayName: string) => void;
};

const defaultProps: AuthProps = {
    loaded: false,
    api: new SDK(TOKEN, REFRESH_TOKEN, storeToken),
    displayName: localStorage.getItem('PREVIOUS_DISPLAY_NAME'),

    mode: 'photobooth',
    event: null,
    time: null,

    hardwareFlashPowered: false,

    backdropAlbum: null,
    selectedBackdrop: 0,

    karaoke: null,
    karaokeQueue: [],

    syncInProgress: false,
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
    setKaraoke: () => {},
    setTimecode: () => {},
    setKaraokeQueue: () => {},
    setSyncInProgress: () => {},
    setBackdrops: () => {},
    setHardwareFlashPowered: () => {},

    setDisplayName: () => {},
});

export default function AuthProvider({ children }: { children: ReactNode }) {
    const { topics } = useSettings();
    const [context, setContext] = useState<AuthProps>({
        ...defaultProps,
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

    const setBackdrops = (backdropAlbum: BackdropAlbum | null, selectedBackdrop: number) =>
        setContext((oldCtx) => ({
            ...oldCtx,
            backdropAlbum,
            selectedBackdrop,
        }));

    const setHardwareFlashPowered = (powered: boolean) =>
        setContext((oldCtx) => ({
            ...oldCtx,
            hardwareFlashPowered: powered,
        }));

    useAsyncEffect(async () => {
        context.api.setOnExpired(() => setToken());

        if (!isLoggedIn()) {
            return;
        }

        const status = await context.api.state.get();

        setContext((oldCtx) => ({
            ...oldCtx,
            mode: status.currentMode,
            event: status.currentEvent,
            backdropAlbum: status.backdropAlbum,
            selectedBackdrop: status.selectedBackdrop,
            karaoke: status.karaoke,
            karaokeQueue: status.karaokeQueue,
            syncInProgress: status.syncInProgress,
            hardwareFlashPowered: status.hardwareFlashPowered,
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
                setMode,
                setKaraoke,
                setTimecode,
                setKaraokeQueue,
                setDisplayName,
                setSyncInProgress,
                setBackdrops,
                setHardwareFlashPowered,
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
