import { Button, Card, CardActions, CardContent, Typography } from "@mui/material";
import ExportListing from "../../../components/admin/export_listing";
import { useAdminSocket } from "../../../hooks/adminSocket";
import { useTranslation } from "react-i18next";

export default function AdminPhotobooth() {
    const {t} = useTranslation();
    const {sendMessage, appState} = useAdminSocket();

    const hasEvent = !!appState.app_state.current_event;

    return <>
        <Card>
            <CardContent>
                {
                    hasEvent &&
                    <>
                        <Typography variant="h2" fontSize={18}>
                            {t('admin_main.current_event')}:  {appState?.app_state.current_event?.name}
                        </Typography>
                        <ul>
                            <li>{t('photobooth.amt_hand_taken')}: {appState?.app_state?.current_event?.amt_images_handtaken}</li>
                            <li>{t('photobooth.amt_unattended')}: {appState?.app_state?.current_event?.amt_images_unattended}</li>
                        </ul>
                    </>
                }
                {
                    !hasEvent && <p>{t('osd.no_event')}</p>
                }
            </CardContent>
        </Card>
        <Card>
            <CardActions>
                <Button style={{ width: '100%' }} onClick={() => sendMessage('photobooth/REMOTE_TAKE_PICTURE', null)}>{t('photobooth.remote_take_picture')}</Button>
            </CardActions>
        </Card>

        { hasEvent && <ExportListing /> }
    </>
}