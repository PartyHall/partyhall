import { Alert, Snackbar } from "@mui/material";
import { createContext, ReactNode, useContext, useEffect, useRef, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { TextLoader } from "../components/loader";
import { AppState } from "../types/appstate";
import { WsMessage } from "../types/ws_message";
import { useSnackbar } from "./snackbar";
import { SOCKET_MODE_DEBUG } from "./useApi";
import useRefresher from "./useRefresher";

type WebsocketProps = {
    lastMessage: WsMessage | null;
    appState: AppState;
    currentTime: string | null;

    currentTimer: number;
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

    currentTimer: -1,
};

const WebsocketContext = createContext<WebsocketContextProps>({
    ...defaultState,
    sendMessage: (msgType: string, data?: any) => { },
    showDebug: () => { },
});

export default function BoothSocketProvider({ children }: { children: ReactNode }) {
    const { sendJsonMessage, lastMessage, readyState } = useWebSocket(
        `ws://${window.location.host}/api/socket/booth`,
        {
            shouldReconnect: () => true,
            reconnectAttempts: 10,
            reconnectInterval: 3000,
        }
    );
    const [ctx, setContext] = useState<WebsocketProps>(defaultState);
    const { showSnackbar } = useSnackbar();

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
        const newContext = { ...ctx, lastMessage: data };

        switch (data.type) {
            case "PING":
                sendJsonMessage({ 'type': 'PONG' })
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
        sendMessage: (msgType: string, data?: any) => sendJsonMessage({ type: msgType, payload: data }),
        showDebug,
    }}>
        <>
            <TextLoader loading={readyState != ReadyState.OPEN || !ctx.appState} text={connectionStatus}>
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

export function useBoothSocket(): WebsocketContextProps {
    return useContext<WebsocketContextProps>(WebsocketContext);
}