import { Alert, Snackbar } from "@mui/material";
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
    showDebug: () => void;
};

const defaultState: WebsocketProps = {
    lastMessage: null,
    // Ignoring this error, the appstate will NOT be null once the app is loaded
    // @ts-ignore
    appState: null,
    currentTime: null,
    lastError: null,
};

const WebsocketContext = createContext<WebsocketContextProps>({
    ...defaultState,
    sendMessage: (msgType: string, data?: any) => { },
    showDebug: () => { },
});

export default function BoothSocketProvider({ children }: { children: ReactNode }) {
    const HOST = `ws://${window.location.host}/api/socket/booth`
    const { sendMessage, lastMessage, readyState } = useWebSocket(HOST);
    const [ctx, setContext] = useState<WebsocketProps>(defaultState);
    const {showSnackbar} = useSnackbar();

    const connectionStatus = useRefresher(readyState, 5);

    const showDebug = () => {
        // @ts-ignore
        setContext({ ...ctx, appState: { ...ctx.appState, debug: true } });
        setTimeout(() => {
            // @ts-ignore
            setContext({ ...ctx, appState: { ...ctx.appState, debug: false } });
        }, 20000);
    };

    useEffect(() => {
        if (!lastMessage) {
            return;
        }

        const data = JSON.parse(lastMessage.data);
        const newContext = {...ctx, lastMessage: data};

        switch (data.type) {
            case "PING":
                sendMessage('{"type": "PONG"}')
                newContext.currentTime = data.payload;
                break;
            case "APP_STATE":
                newContext.appState = data.payload;
                break
            case "DISPLAY_DEBUG":
                showDebug();
                return
            case "ERR_MODAL":
                showSnackbar(data.payload, 'error');
                break
        }

        setContext(newContext);
    }, [lastMessage]);

    return <WebsocketContext.Provider value={{
        ...ctx,
        sendMessage: (msgType: string, data?: any) => sendMessage(JSON.stringify({ type: msgType, payload: data })),
        showDebug,
    }}>
        <TextLoader loading={readyState != ReadyState.OPEN || !ctx.appState} text={connectionStatus}>
            {children}
        </TextLoader>
    </WebsocketContext.Provider>
}

export function useWebsocket(): WebsocketContextProps {
    return useContext<WebsocketContextProps>(WebsocketContext);
}