import { Alert, AlertColor, Snackbar } from "@mui/material";
import { createContext, ReactNode, useContext, useEffect, useState } from "react";
import Login from "../pages/Login";
import { EventExport } from "../types/event_export";
import { SnackbarData } from "../types/SnackbarData";

const buildSnackbarMessage = (message: string, type: AlertColor = 'error') => {
    return { message, open: true, type };
};

type ApiProps = {
    connecting: boolean;
    password: string | null;
    snackbarMessage: SnackbarData | null;
};

type ApiContextProps = ApiProps & {
    connect: (password: string) => void;
    getLastExports: (eventId: number) => Promise<EventExport[]>;
    logout: () => void;
    showError: (message: string, type: AlertColor) => void;
};

const defaultState: ApiProps = {
    connecting: false,
    password: localStorage.getItem('password'),
    snackbarMessage: null,
};

const ApiContext = createContext<ApiContextProps>({
    ...defaultState,
    connect: () => { },
    logout: () => { },
    getLastExports: async () => [],
    showError: () => {},
});

export default function ApiProvider({ children }: { children: ReactNode }) {
    const [ctx, setContext] = useState<ApiProps>(defaultState);

    const connect = async (password: string) => {
        setContext({ ...ctx, connecting: true });

        const resp = await fetch(`/api/admin/login`, {
            'method': 'POST',
            'headers': {
                'Authorization': password,
            }
        });

        if (resp.status === 401) {
            setContext({ ...ctx, snackbarMessage: buildSnackbarMessage('Wrong password !') });
            return;
        } else {
            const data = await resp.text();

            if (data !== 'yes') {
                setContext({ ...ctx, snackbarMessage: buildSnackbarMessage('Wrong response from the partyhall') })
                return;
            } else {
                localStorage.setItem('password', password);
                setContext({ ...ctx, snackbarMessage: null, password });
            }
        }
    };

    const logout = () => {
        localStorage.removeItem('password');
        setContext({ ...ctx, password: null });
    };

    const showError = (message: string, severity: AlertColor = 'error') => setContext({...ctx, snackbarMessage: buildSnackbarMessage(message, severity)});

    const getLastExports = async (eventId: number) => {
        const resp = await fetch(`/api/admin/exports/${eventId}`, {
            'method': 'GET',
            'headers': { 'Authorization': ctx.password ?? '' }
        });

        if (resp.status === 401) {
            setContext({ ...ctx, password: null, snackbarMessage: buildSnackbarMessage('Session expired') });
        } else {
            return await resp.json();
        }

        return [];
    };

    const closeSnackbar = () => {
        const newCtx = {...ctx, snackbarMessage: null};
        
        if (!!ctx.snackbarMessage) {
            // fck ts sometimes
            // @ts-ignore
            newCtx.snackbarMessage = {...ctx.snackbarMessage, open: false};
        }

        setContext(newCtx);
    }

    return <ApiContext.Provider value={{
        ...ctx,
        connect,
        logout,
        showError,
        getLastExports,
    }}>
        <>
            {(!!ctx.password) && <>{children}</>}

            {!ctx.password  && <Login />}

            <Snackbar open={!!ctx.snackbarMessage?.open} autoHideDuration={6000} onClose={closeSnackbar} anchorOrigin={{ vertical: "bottom", horizontal: "center" }}>
                <Alert onClose={closeSnackbar} severity={ctx.snackbarMessage?.type} sx={{ width: '100%' }}>
                    {
                        ctx.snackbarMessage?.message
                    }
                </Alert>
            </Snackbar>
        </>
    </ApiContext.Provider>
}

export function useApi(): ApiContextProps {
    return useContext<ApiContextProps>(ApiContext);
}