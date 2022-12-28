import { createContext, ReactNode, useContext, useEffect, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { TextLoader } from "../components/loader";
import { AppState } from "../types/appstate";
import { WsMessage } from "../types/ws_message";
import { useSnackbar } from "./snackbar";
import useRefresher from "./useRefresher";

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
    const { sendMessage, lastMessage, readyState } = useWebSocket(`ws://${window.location.host}/api/socket/admin`);

    const { showSnackbar } = useSnackbar();
    const [ctx, setContext] = useState<WebsocketProps>(defaultState);
    const connectionStatus = useRefresher(readyState, 5);

    useEffect(() => {
        if (!lastMessage) {
            return;
        }

        const data = JSON.parse(lastMessage.data);

        switch (data.type) {
            case "PING":
                sendMessage('{"type": "PONG"}')
                setContext({...ctx, lastMessage: data, currentTime: data.payload})
                break
            case "APP_STATE":
                setContext({...ctx, lastMessage: data, appState: data.payload})
                break
            case "ERR_MODAL":
                showSnackbar(data.payload, 'error');
                setContext({...ctx, lastMessage: data});
                break
            case "EXPORT_STARTED":
                showSnackbar('Export started', 'info');
                setContext({...ctx, lastMessage: data});
                break
            case "EXPORT_COMPLETED":
                showSnackbar('Export completed', 'success');
                setContext({...ctx, lastMessage: data});
                break
            default:
                setContext({...ctx, lastMessage: data});
        }
    }, [lastMessage]);

    return <WebsocketContext.Provider value={{
        ...ctx,
        sendMessage: (msgType: string, data?: any) => sendMessage(JSON.stringify({ type: msgType, payload: data })),
    }}>
        <>
            <TextLoader loading={readyState != ReadyState.OPEN || !ctx.appState} text={connectionStatus}>
                {children}
            </TextLoader>
        </>
    </WebsocketContext.Provider>
}

export function useAdminSocket(): WebsocketContextProps {
    return useContext<WebsocketContextProps>(WebsocketContext);
}
