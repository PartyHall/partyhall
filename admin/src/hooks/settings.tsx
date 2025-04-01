import { ReactNode, createContext, useContext, useState } from 'react';
import Loader from '../components/loader';
import { PhUserSettings } from '@partyhall/sdk';
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
];

/** @TODO: Implement APIs */
type SettingsProps = {
    loaded: boolean;
    pageName: string;

    user_settings: PhUserSettings | null;

    guests_allowed: boolean;
    enabled_modules: string[];

    hwflash_powered: boolean;

    version: string;
    commit: string;

    topics: string[];
};

type SettingsContextProps = SettingsProps & {
    fetch: () => Promise<void>;
    setPageName: (name: string, mercureTopics?: string[]) => void;
    setUserSettings: (user_settings: PhUserSettings) => void;
};

const defaultProps: SettingsProps = {
    loaded: false,
    pageName: 'home',
    topics: DEFAULT_TOPICS,

    user_settings: null,

    guests_allowed: false,
    enabled_modules: [],

    hwflash_powered: false,

    version: 'INDEV',
    commit: 'XXXXXX',
};

const SettingsContext = createContext<SettingsContextProps>({
    ...defaultProps,
    fetch: async () => {},
    setPageName: () => {},
    setUserSettings: () => {},
});

export default function SettingsProvider({ children }: { children: ReactNode }) {
    const [ctx, setCtx] = useState<SettingsProps>(defaultProps);

    const fetchStatus = async () => {
        setCtx((oldCtx) => ({ ...oldCtx, loaded: false }));
        const resp = await fetch('/api/state');
        const data = await resp.json();

        setCtx((oldCtx) => ({
            ...oldCtx,
            loaded: true,
            ...data,
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

    const setUserSettings = (user_settings: PhUserSettings) => setCtx((oldCtx) => ({ ...oldCtx, user_settings }));

    useAsyncEffect(fetchStatus, []);

    return (
        <SettingsContext.Provider value={{ ...ctx, fetch: fetchStatus, setPageName, setUserSettings }}>
            <Loader loading={!ctx.loaded}>{children}</Loader>
        </SettingsContext.Provider>
    );
}

export function useSettings() {
    return useContext<SettingsContextProps>(SettingsContext);
}
