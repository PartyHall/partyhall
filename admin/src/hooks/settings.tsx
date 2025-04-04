import { PhState, PhUserSettings } from '@partyhall/sdk';
import { ReactNode, createContext, useContext, useState } from 'react';
import Loader from '../components/loader';
import useAsyncEffect from 'use-async-effect';

const DEFAULT_TOPICS = [
    '/time',
    '/mode',
    '/event',
    '/snackbar',
    '/karaoke',
    '/karaoke_queue',
    '/sync-progress',
    '/flash',
    '/backdrop-state',
    '/user-settings',
    '/btn-press',
];

/** @TODO: Implement APIs */
type SettingsProps = {
    loaded: boolean;
    pageName: string;

    userSettings: PhUserSettings | null;

    guestsAllowed: boolean;

    version: string;
    commit: string;

    topics: string[];
};

type SettingsContextProps = SettingsProps & {
    fetch: () => Promise<void>;
    setPageName: (name: string, mercureTopics?: string[]) => void;
    setUserSettings: (userSettings: PhUserSettings) => void;
    setHardwareFlashBrightness: (b: number) => void;
};

const defaultProps: SettingsProps = {
    loaded: false,
    pageName: 'home',
    topics: DEFAULT_TOPICS,

    userSettings: null,

    guestsAllowed: false,

    version: 'INDEV',
    commit: 'XXXXXX',
};

const SettingsContext = createContext<SettingsContextProps>({
    ...defaultProps,
    fetch: async () => {},
    setPageName: () => {},
    setUserSettings: () => {},
    setHardwareFlashBrightness: () => {},
});

export default function SettingsProvider({ children }: { children: ReactNode }) {
    const [ctx, setCtx] = useState<SettingsProps>(defaultProps);

    const fetchStatus = async () => {
        setCtx((oldCtx) => ({ ...oldCtx, loaded: false }));
        const resp = await fetch('/api/state');
        const data = await resp.json();

        const state = PhState.fromJson(data);

        setCtx((oldCtx) => ({
            ...oldCtx,
            loaded: true,
            userSettings: state.userSettings,
            guestsAllowed: state.guestsAllowed,
            version: state.version,
            commit: state.commit,
        }));
    };

    const setPageName = (name: string, mercureTopics?: string[]) => {
        let topics = [...DEFAULT_TOPICS];

        if (mercureTopics) {
            topics = [...topics, ...mercureTopics];
        }

        setCtx((oldCtx) => ({
            ...oldCtx,
            pageName: name,
            topics,
        }));
    };

    const setUserSettings = (userSettings: PhUserSettings) => setCtx((oldCtx) => ({ ...oldCtx, userSettings }));

    const setHardwareFlashBrightness = (b: number) => {
        const settings = ctx.userSettings;
        if (!settings) {
            return;
        }

        const photoboothSettings = settings.photobooth;
        photoboothSettings.flashBrightness = b;
        settings.photobooth = photoboothSettings;

        setUserSettings(settings);
    };

    useAsyncEffect(fetchStatus, []);

    return (
        <SettingsContext.Provider
            value={{ ...ctx, fetch: fetchStatus, setPageName, setUserSettings, setHardwareFlashBrightness }}
        >
            <Loader loading={!ctx.loaded}>{children}</Loader>
        </SettingsContext.Provider>
    );
}

export function useSettings() {
    return useContext<SettingsContextProps>(SettingsContext);
}
