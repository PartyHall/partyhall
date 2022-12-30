import { Button, Card, CardActions, CardContent, Typography } from "@mui/material";
import ExportListing from "../../../components/admin/export_listing";
import { useAdminSocket } from "../../../hooks/adminSocket";

export default function AdminPhotobooth() {
    const {sendMessage, appState} = useAdminSocket();

    const hasEvent = !!appState.app_state.current_event;

    return <>
        <Card>
            <CardContent>
                {
                    hasEvent &&
                    <>
                        <Typography variant="h2" fontSize={18}>
                            Current event:  {appState?.app_state.current_event?.name}
                        </Typography>
                        <ul>
                            <li>Amount of picture handtaken: {appState?.app_state?.current_event?.amt_images_handtaken}</li>
                            <li>Amount of picture unattended: {appState?.app_state?.current_event?.amt_images_unattended}</li>
                        </ul>
                    </>
                }
                {
                    !hasEvent && <p>No event selected</p>
                }
            </CardContent>
        </Card>
        <Card>
            <CardActions>
                <Button style={{ width: '100%' }} onClick={() => sendMessage('REMOTE_TAKE_PICTURE', null)}>Remote take a picture</Button>
            </CardActions>
        </Card>

        {
            hasEvent && <ExportListing />
        }
    </>
}