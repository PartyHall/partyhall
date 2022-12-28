import { Alert, AlertColor, Snackbar } from "@mui/material";
import { createContext, ReactNode, useContext, useState } from "react";

type SnackbarProps = {
    open: boolean;
    message: string | null;
    type: AlertColor;
};

type SnackbarContextProps = SnackbarProps & {
    showSnackbar: (message: string, type: AlertColor) => void;
};

const defaultState: SnackbarProps = {
    open: false,
    message: null,
    type: 'error',
};

const SnackbarContext = createContext<SnackbarContextProps>({
    ...defaultState,
    showSnackbar: () => { },
});

export default function SnackbarProvider({ children }: { children: ReactNode }) {
    const [context, setContext] = useState<SnackbarProps>(defaultState);

    const show = (message: string, type: AlertColor) => setContext({...context, open: true, message, type});
    const close = () => setContext({...context, open: false});

    return <SnackbarContext.Provider value={{
        ...context,
        showSnackbar: show,
    }}>
        {children}
        <Snackbar open={context.open} autoHideDuration={6000} onClose={close} anchorOrigin={{ vertical: "bottom", horizontal: "center" }}>
            <Alert onClose={close} severity={context.type} sx={{ width: '100%' }}>
                {
                    context.message
                }
            </Alert>
        </Snackbar>
    </SnackbarContext.Provider>;
}

export const useSnackbar = () => useContext<SnackbarContextProps>(SnackbarContext);