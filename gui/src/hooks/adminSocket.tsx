import { createContext, ReactNode, useContext, useEffect, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { TextLoader } from "../components/loader";
import { AppState } from "../types/appstate";
import { WsMessage } from "../types/ws_message";
import { useSnackbar } from "./snackbar";
import { SOCKET_MODE_DEBUG, useApi } from "./useApi";
import { useTranslation } from "react-i18next";

type WebsocketProps = {
    lastMessage: WsMessage | null;
    appState: AppState;
    currentTime: string | null;
};

type WebsocketContextProps = WebsocketProps & {
    sendMessage: (msgType: string, data?: any | null) => void;
};

const defaultState: WebsocketProps = {
    lastMessage: null,
    //@ts-ignore
    appState: null,
    currentTime: null,
};

const WebsocketContext = createContext<WebsocketContextProps>({
    ...defaultState,
    sendMessage: () => { },
});

export default function AdminSocketProvider({ children }: { children: ReactNode }) {
    const { t, i18n } = useTranslation();
    const [currentLanguage, setCurrentLanguage] = useState<string>('en');
    const { api, logout } = useApi();
    const { sendJsonMessage, lastMessage, readyState } = useWebSocket(
        `ws://${window.location.host}/api/socket/admin`,
        {
            shouldReconnect: () => true,
            reconnectAttempts: 10,
            reconnectInterval: 3000,
            queryParams: { token: api.getToken() ?? '' },
            onError: () => logout(), // @TODO: find a way to check if 401 and only in this case logout
        }
    );

    const { showSnackbar } = useSnackbar();
    const [ctx, setContext] = useState<WebsocketProps>(defaultState);

    useEffect(() => {
        if (!lastMessage) {
            return;
        }

        const data = JSON.parse(lastMessage.data);

        switch (data.type) {
            case "PING":
                sendJsonMessage({ 'type': 'PONG' })
                setContext({ ...ctx, lastMessage: data, currentTime: data.payload })
                break
            case "APP_STATE":
                setContext({ ...ctx, lastMessage: data, appState: data.payload })
                break
            case "ERR_MODAL":
                showSnackbar(data.payload, 'error');
                setContext({ ...ctx, lastMessage: data });
                break
            case "EXPORT_STARTED":
                showSnackbar(t('exports.started'), 'info');
                setContext({ ...ctx, lastMessage: data });
                break
            case "EXPORT_COMPLETED":
                showSnackbar(t('exports.completed'), 'success');
                setContext({ ...ctx, lastMessage: data });
                break
            default:
                setContext({ ...ctx, lastMessage: data });
        }
    }, [lastMessage]);

    useEffect(() => {
        if (!ctx.appState || !ctx.appState.language) {
            return;
        }

        if (ctx.appState.language === currentLanguage) {
            return;
        }

        i18n.changeLanguage(ctx.appState.language);
        setCurrentLanguage(ctx.appState.language);
    }, [ctx]);

    return <WebsocketContext.Provider value={{
        ...ctx,
        sendMessage: (msgType: string, data?: any) => sendJsonMessage({ type: msgType, payload: data }),
    }}>
        <>
            <TextLoader loading={readyState != ReadyState.OPEN || !ctx.appState} text="Connecting...">
                {children}
            </TextLoader>

            {
                SOCKET_MODE_DEBUG &&
                <div className="debug absbl">
                    <p>Last message: {ctx.lastMessage?.type}</p>
                </div>
            }
        </>
    </WebsocketContext.Provider>
}

export function useAdminSocket(): WebsocketContextProps {
    return useContext<WebsocketContextProps>(WebsocketContext);
}
