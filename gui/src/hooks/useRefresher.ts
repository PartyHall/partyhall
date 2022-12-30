import { useEffect, useState } from "react";
import { ReadyState } from "react-use-websocket";

/**
 * This hooks checks whether the websocket is still active
 * If it's not, it waits for some amount of time then refresh the full page to try again
 * It also returns a user-friendly text to be displayed regarding the current status
 */
export default function useRefresher(readyState: ReadyState, reloadTimeout: number) {
    const [killSwitch, setKillswitch] = useState<number>(-1);

    useEffect(() => {
        if ([ReadyState.CLOSING, ReadyState.CLOSED, ReadyState.UNINSTANTIATED,].includes(readyState)) {
            setKillswitch(setTimeout(() => {
                window.location.reload();
            }, reloadTimeout * 1000));
        } else {
            if (killSwitch >= 0) {
                clearTimeout(killSwitch);
                setKillswitch(-1);
            }
        }
    }, [readyState]);

    return {
        [ReadyState.CONNECTING]: "Connecting",
        [ReadyState.OPEN]: "Open - Waiting for state",
        [ReadyState.CLOSING]: "Websocket closing",
        [ReadyState.CLOSED]: "Websocket closed",
        [ReadyState.UNINSTANTIATED]: "Websocket uninstantiated"
    }[readyState];
}