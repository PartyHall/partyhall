import { BackdropAlbum, PhEvent, SDK } from '@partyhall/sdk';
import { ReactNode, createContext, useContext, useEffect, useRef, useState } from 'react';
import { DateTime } from 'luxon';
import { notification } from 'antd';
import { useAuth } from './auth';

type EventListenerCallback = (data: any) => void;
type EventListenerFn = (eventType: string, callback: EventListenerCallback) => void;

type MercureProps = {
    isInitialized: boolean;

    time: DateTime | null;
};

type MercureContextProps = MercureProps & {
    addMercureListener: EventListenerFn;
    removeMercureListener: EventListenerFn;
};

const defaultProps: MercureProps = {
    time: null,
    isInitialized: false,
};

const MercureContext = createContext<MercureContextProps>({
    ...defaultProps,
    addMercureListener: () => {},
    removeMercureListener: () => {},
});

export default function MercureProvider({
    url,
    sdk,
    topics,
    children,
}: {
    url: string;
    sdk?: SDK | null;
    topics: string[];
    children: ReactNode;
}) {
    const [ctx, setCtx] = useState<MercureProps>(defaultProps);
    const { setEvent, setMode, isLoggedIn, setSyncInProgress, setHardwareFlash, setBackdrops } = useAuth();
    const [notif, ctxHolder] = notification.useNotification();

    const eventSource = useRef<EventSource>();
    const listenersRef = useRef<{ event: string; callback: EventListener }[]>([]);

    const createEventSource = () => {
        const urlObj = new URL(url);
        topics.forEach((x) => {
            urlObj.searchParams.append('topic', x);
        });

        const es = new EventSource(urlObj, { withCredentials: true });
        es.onerror = (e) => {
            const target = e.target as EventSource;

            if (target.readyState === EventSource.CLOSED) {
                es.close();

                setTimeout(() => {
                    const newEs = createEventSource();
                    eventSource.current = newEs;
                }, 500);

                return;
            }

            notif.error({
                message: 'Connection lost!',
                description: 'The connection to Mercure was lost. You might miss some updates.',
            });
        };

        eventSource.current = es;
        setCtx((oldCtx) => ({ ...oldCtx, isInitialized: true }));

        es.addEventListener('/time', (x) => {
            setCtx((oldCtx) => ({
                ...oldCtx,
                time: DateTime.fromISO(JSON.parse(x.data).time),
            }));
        });

        es.addEventListener('/event', (x) => {
            const event = PhEvent.fromJson(JSON.parse(x.data));
            if (!event) {
                return;
            }

            setEvent(event);
        });

        es.addEventListener('/flash', (x) => {
            const data = JSON.parse(x.data);
            setHardwareFlash(data.powered, data.brightness);
        });

        es.addEventListener('/mode', (x) => setMode(JSON.parse(x.data).mode));

        es.addEventListener('/sync-progress', (x) => setSyncInProgress(JSON.parse(x.data).syncInProgress));

        es.addEventListener('/snackbar', (x) => {
            // @TODO: Handle better
            const data = JSON.parse(x.data);
            notif.open({
                message: 'Incoming transmission',
                description: data['msg'],
            });
        });

        es.addEventListener('/backdrop-state', (x) => {
            const data = JSON.parse(x.data);

            setBackdrops(BackdropAlbum.fromJson(data['backdrop_album']), data['selected_backdrop']);
        });

        listenersRef.current.forEach((l) => {
            es.addEventListener(l.event, l.callback);
        });

        return es;
    };

    useEffect(() => {
        if (!sdk || !sdk.token || topics.length === 0 || !isLoggedIn()) {
            return;
        }
        setCtx((oldCtx) => ({ ...oldCtx, isInitialized: false }));

        const es = createEventSource();

        return () => {
            eventSource.current = undefined;
            es.close();
        };
    }, [sdk, topics]);

    const addMercureListener = (event: string, callback: EventListener) => {
        if (eventSource.current) {
            eventSource.current.addEventListener(event, callback);
        }

        listenersRef.current.push({ event, callback });
    };

    const removeMercureListener = (event: string, callback: EventListener) => {
        if (eventSource.current) {
            eventSource.current.removeEventListener(event, callback);
        }

        listenersRef.current = listenersRef.current.filter(
            (listener) => !(listener.event === event && listener.callback === callback)
        );
    };

    return (
        <MercureContext.Provider
            value={{
                ...ctx,
                addMercureListener,
                removeMercureListener,
            }}
        >
            {children}
            {ctxHolder}
        </MercureContext.Provider>
    );
}

export function useMercure() {
    return useContext<MercureContextProps>(MercureContext);
}

export function useMercureTopic<T>(event: string, listener: (data: T) => void) {
    const { addMercureListener, removeMercureListener } = useMercure();

    useEffect(() => {
        const method = (el: MessageEvent) => {
            return listener(JSON.parse(el.data));
        };

        addMercureListener(event, method);

        return () => {
            removeMercureListener(event, method);
        };
    }, []);
}
