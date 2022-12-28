import { Alert, Snackbar } from "@mui/material";
import { createContext, ReactNode, useContext, useEffect, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { TextLoader } from "../components/loader";
import { AppState } from "../types/appstate";
import { WsMessage } from "../types/ws_message";

type WebsocketProps = {
    lastMessage: WsMessage | null;
    appState: AppState | null;
    currentTime: string | null;
    lastError: string | null
};

type WebsocketContextProps = WebsocketProps & {
    sendMessage: (msgType: string, data?: any | null) => void;
    showDebug: () => void;
};

const defaultState: WebsocketProps = {
    lastMessage: null,
    appState: null,
    currentTime: null,
    lastError: null,
};

const WebsocketContext = createContext<WebsocketContextProps>({
    ...defaultState,
    sendMessage: (msgType: string, data?: any) => { },
    showDebug: () => {},
});

export default function WebsocketProvider({ children }: { children: ReactNode }) {
    const HOST = `ws://${window.location.host}/api/socket/booth`
    const [killSwitch, setKillswitch] = useState<number>(-1);
    const { sendMessage, lastMessage, readyState } = useWebSocket(HOST);
    const [ctx, setContext] = useState<WebsocketProps>(defaultState);

    useEffect(() => {
        if ([
            ReadyState.CLOSING,
            ReadyState.CLOSED,
            ReadyState.UNINSTANTIATED,
        ].includes(readyState)) {
            setKillswitch(setTimeout(() => {
                window.location.reload();
            }, 5000));
        } else {
            if (killSwitch >= 0) {
                clearTimeout(killSwitch);
                setKillswitch(-1);
            }
        }
    }, [readyState]);

    const connectionStatus = {
        [ReadyState.CONNECTING]: "Connecting",
        [ReadyState.OPEN]: "Open",
        [ReadyState.CLOSING]: "Websocket closing",
        [ReadyState.CLOSED]: "Websocket closed",
        [ReadyState.UNINSTANTIATED]: "Websocket uninstantiated"
    }[readyState];

    useEffect(() => {
        if (!lastMessage) {
            return;
        }


        const data = JSON.parse(lastMessage.data);
        let newCtx = { ...ctx, lastMessage: data };

        switch (data.type) {
            case "PING":
                sendMessage('{"type": "PONG"}')
                newCtx.currentTime = data.payload;
                break
            case "APP_STATE":
                newCtx.appState = data.payload;
                break
            case "ERR_MODAL":
                newCtx.lastError = data.payload;
                break
        }

        setContext(newCtx)
    }, [lastMessage]);

    // @TODO: Make the timer clientside with config
    // Remove the timer serverside for hiding the debug mode
    const showDebug = () => {
        // @ts-ignore
        setContext({...ctx, appState: {...ctx.appState, debug: true}});
        setTimeout(() => {
            // @ts-ignore
            setContext({...ctx, appState: {...ctx.appState, debug: false}});
        }, 20000);
    };

    return <WebsocketContext.Provider value={{
        ...ctx,
        sendMessage: (msgType: string, data?: any) => sendMessage(JSON.stringify({ type: msgType, payload: data })),
        showDebug,
    }}>
        <>
            <TextLoader loading={readyState != ReadyState.OPEN} text={connectionStatus}>
                {children}
            </TextLoader>

            <Snackbar open={!!ctx.lastError} autoHideDuration={6000} onClose={() => setContext({ ...ctx, lastError: null })} anchorOrigin={{ vertical: "bottom", horizontal: "center" }}>
                <Alert onClose={() => setContext({ ...ctx, lastError: null })} severity="error" sx={{ width: '100%' }}>
                    {ctx.lastError}
                </Alert>
            </Snackbar>
        </>
    </WebsocketContext.Provider>
}

export function useWebsocket(): WebsocketContextProps {
    return useContext<WebsocketContextProps>(WebsocketContext);
}