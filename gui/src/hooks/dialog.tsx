import { Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from "@mui/material";
import { createContext, ReactNode, useContext, useState } from "react";
import { useTranslation } from "react-i18next";

type ConfirmDialogProps = {
    open: boolean;
    title: string | null;
    message: string | ReactNode | null;
    validateButton: string | null;
    action: () => Promise<void>;
};

type ConfirmDialogContextProps = ConfirmDialogProps & {
    showDialog: (title: string, message: string | ReactNode, validateButton: string, action: () => Promise<void>) => void;
};

const defaultState: ConfirmDialogProps = {
    open: false,
    title: null,
    message: null,
    validateButton: null,
    action: async () => { },
};

const ConfirmDialogContext = createContext<ConfirmDialogContextProps>({
    ...defaultState,
    showDialog: () => { },
});

export default function ConfirmDialogProvider({ children }: { children: ReactNode }) {
    const {t} = useTranslation();
    const [context, setContext] = useState<ConfirmDialogProps>(defaultState);

    const show = (title: string, message: string | ReactNode, validateButton: string, action: () => Promise<void>) => setContext({ ...context, open: true, title, message, validateButton, action });
    const close = () => setContext({ ...context, open: false });

    return <ConfirmDialogContext.Provider value={{
        ...context,
        showDialog: show,
    }}>
        {children}
        <Dialog open={context.open} onClose={close}>
            <DialogTitle>{context.title}</DialogTitle>
            <DialogContent>
                <DialogContentText>{context.message}</DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={close}>{t('general.cancel')}</Button>
                <Button onClick={async () => {
                    await context.action();
                    close();
                }} color="error" autoFocus>{context.validateButton}</Button>
            </DialogActions>
        </Dialog>
    </ConfirmDialogContext.Provider>;
}

export const useConfirmDialog = () => useContext<ConfirmDialogContextProps>(ConfirmDialogContext);