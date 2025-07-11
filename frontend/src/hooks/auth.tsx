import { BackdropAlbum, PhEvent, PhKaraoke, PhSongSession, PhUserSettings, SDK } from '@partyhall/sdk';
import { ReactNode, createContext, useContext, useEffect, useState } from 'react';

import { DateTime } from 'luxon';

import useAsyncEffect from 'use-async-effect';
import { useSnackbar } from 'notistack';

type AuthProps = {
    loaded: boolean;
    api: SDK;

    time: DateTime | null;
    currentMode: string;
    currentEvent: PhEvent | null;

    debug: {
        ip_addresses: {
            [key: string]: string[];
        };
    } | null;

    backdropAlbum: BackdropAlbum | null;
    selectedBackdrop: number;

    karaoke: PhKaraoke;
    karaokeQueue: PhSongSession[];

    userSettings: PhUserSettings | null;
    lastBtnPressed: number | null;

    version: string | null;
    commit: string | null;

    shouldTakePicture: 'unattended' | 'normal' | false;
};

type AuthContextProps = AuthProps & {
    setPictureTaken: () => void;
};

const defaultProps: AuthProps = {
    api: new SDK(null, null),
    loaded: false,

    time: null,
    currentMode: 'photobooth',
    currentEvent: null,
    debug: null,

    backdropAlbum: null,
    selectedBackdrop: 0,

    karaoke: new PhKaraoke({
        current: null,
        is_playing: false,
        timecode: 0,
        countdown: 0,
        volume: 0,
        volume_vocals: 0,
    }),

    karaokeQueue: [],

    userSettings: null,
    lastBtnPressed: null,

    version: null,
    commit: null,

    shouldTakePicture: false,
};

const AuthContext = createContext<AuthContextProps>({
    ...defaultProps,
    setPictureTaken: () => {},
});

export default function AuthProvider({ children, token }: { children: ReactNode; token: string }) {
    const { enqueueSnackbar } = useSnackbar();
    const [ctx, setCtx] = useState<AuthProps>(defaultProps);

    const createEventSource = () => {
        const url = new URL(`${window.location.protocol}//${window.location.host}/.well-known/mercure`);
        [
            '/time',
            '/mode',
            '/event',
            '/ip-addresses',
            '/debug',
            '/take-picture',
            '/snackbar',
            '/karaoke',
            '/karaoke-queue',
            '/backdrop-state',
            '/btn-press',
        ].forEach((x) => url.searchParams.append('topic', x));

        const es = new EventSource(url, { withCredentials: true });
        es.onmessage = (x) => console.log(x);

        es.onopen = () => console.log('Mercure connection opened');
        es.onerror = (e: Event) => {
            const target = e.target as EventSource;

            if (target.readyState === EventSource.CLOSED) {
                console.log('Mercure connection closed, reconnecting...');
                es.close();

                setTimeout(createEventSource, 500);

                return;
            }

            enqueueSnackbar('Mercure connection lost!', { variant: 'error' });
            setTimeout(() => window.location.reload(), 3000);
        };

        es.addEventListener('/time', (x) =>
            setCtx((oldCtx) => ({
                ...oldCtx,
                time: DateTime.fromISO(JSON.parse(x.data).time),
            }))
        );

        es.addEventListener('/event', (x) => {
            setCtx((oldCtx) => ({
                ...oldCtx,
                currentEvent: JSON.parse(x.data),
            }));
        });

        es.addEventListener('/mode', (x) =>
            setCtx((oldCtx) => ({
                ...oldCtx,
                currentMode: JSON.parse(x.data).mode,
            }))
        );

        es.addEventListener('/ip-addresses', (x) =>
            setCtx((oldCtx) => ({
                ...oldCtx,
                ip_addresses: JSON.parse(x.data),
            }))
        );

        es.addEventListener('/debug', (x) => {
            setCtx((oldCtx) => ({ ...oldCtx, debug: JSON.parse(x.data) }));
            setTimeout(() => setCtx((oldCtx) => ({ ...oldCtx, debug: null })), 30000);
        });

        es.addEventListener('/snackbar', (x) => {
            const data = JSON.parse(x.data);
            enqueueSnackbar(data['msg'], { variant: data['type'] });
        });

        es.addEventListener('/take-picture', (x) =>
            setCtx((oldCtx) => {
                const pictureType = JSON.parse(x.data).unattended ? 'unattended' : 'normal';

                // If we are manually taking a picture at the same moment
                // Lets just silently ignore the unattended one
                // Maybe this should be fixed later (?)
                // Because its fun to see in the timelapse people being in
                // front of the appliance and which will not happen
                // if we skip the unattended ones
                // if (pictureType === 'normal') {
                //     return oldCtx;
                // }

                return {
                    ...oldCtx,
                    shouldTakePicture: pictureType,
                };
            })
        );

        es.addEventListener('/karaoke-queue', (x) => {
            setCtx((oldCtx) => ({
                ...oldCtx,
                karaokeQueue: JSON.parse(x.data)
                    .map((y: any) => PhSongSession.fromJson(y))
                    .filter((y: any) => !!y),
            }));
        });

        es.addEventListener('/karaoke', (x) => {
            setCtx((oldCtx) => ({
                ...oldCtx,
                karaoke: PhKaraoke.fromJson(JSON.parse(x.data)),
            }));
        });

        es.addEventListener('/backdrop-state', (x) => {
            const data = JSON.parse(x.data);

            setCtx((oldCtx) => ({
                ...oldCtx,
                backdropAlbum: BackdropAlbum.fromJson(data['backdrop_album']),
                selectedBackdrop: data['selected_backdrop'],
            }));
        });

        es.addEventListener('/btn-press', (x) => {
            setCtx((oldCtx) => ({
                ...oldCtx,
                lastBtnPressed: JSON.parse(x.data),
            }));
        });

        return es;
    };

    useEffect(() => {
        setCtx((oldCtx) => ({ ...oldCtx, api: new SDK(token, null) }));

        const es = createEventSource();

        return () => {
            es.close();
        };
    }, [token]);

    useAsyncEffect(async () => {
        if (ctx.loaded || !ctx.api.token) {
            return;
        }

        setCtx((oldCtx) => ({ ...oldCtx, loaded: false }));

        const state = await ctx.api.state.get();

        setCtx((oldCtx) => ({
            ...oldCtx,
            loaded: true,
            currentEvent: state.currentEvent,
            currentMode: state.currentMode,
            backdropAlbum: state.backdropAlbum,
            selectedBackdrop: state.selectedBackdrop,
            karaoke: state.karaoke,
            karaokeQueue: state.karaokeQueue,
            userSettings: state.userSettings,
            version: state.version,
            commit: state.commit,
        }));
    }, [ctx.api]);

    return (
        <AuthContext.Provider
            value={{
                ...ctx,
                setPictureTaken: () =>
                    setCtx((oldCtx) => ({
                        ...oldCtx,
                        shouldTakePicture: false,
                    })),
            }}
        >
            {ctx.loaded && children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    return useContext<AuthContextProps>(AuthContext);
}
