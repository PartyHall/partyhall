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
];

/** @TODO: Implement APIs */
type SettingsProps = {
    loaded: boolean;
    pageName: string;

    modules_settings: Record<string, any>;

    guests_allowed: boolean;
    enabled_modules: string[];

    hwflash_powered: boolean;

    version: string;
    commit: string;
    hwid: string | null;

    topics: string[];
};

type SettingsContextProps = SettingsProps & {
    fetch: () => Promise<void>;
    setPageName: (name: string, mercureTopics?: string[]) => void;
};

const defaultProps: SettingsProps = {
    loaded: false,
    pageName: 'home',
    topics: DEFAULT_TOPICS,

    modules_settings: {},

    guests_allowed: false,
    enabled_modules: [],

    hwflash_powered: false,

    version: 'INDEV',
    commit: 'XXXXXX',
    hwid: null,
};

const SettingsContext = createContext<SettingsContextProps>({
    ...defaultProps,
    fetch: async () => {},
    setPageName: () => {},
});

export default function SettingsProvider({ children }: { children: ReactNode }) {
    const [ctx, setCtx] = useState<SettingsProps>(defaultProps);

    const fetchStatus = async () => {
        setCtx((oldCtx) => ({ ...oldCtx, loaded: false }));
        const resp = await fetch('/api/status');
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

    useAsyncEffect(fetchStatus, []);

    return (
        <SettingsContext.Provider value={{ ...ctx, fetch: fetchStatus, setPageName }}>
            <Loader loading={!ctx.loaded}>{children}</Loader>
        </SettingsContext.Provider>
    );
}

export function useSettings() {
    return useContext<SettingsContextProps>(SettingsContext);
}
