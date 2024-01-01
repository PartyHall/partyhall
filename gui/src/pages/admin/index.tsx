import { Box, Button, Card, CardActions, CardContent, IconButton, MenuItem, Select, SelectChangeEvent, Typography } from "@mui/material";
import { DateTime } from 'luxon';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';

import { useConfirmDialog } from "../../hooks/dialog";
import { useAdminSocket } from "../../hooks/adminSocket";
import { useNavigate } from "react-router-dom";

export default function AdminIndex() {
    const { sendMessage, appState, currentTime } = useAdminSocket();
    const { showDialog } = useConfirmDialog();
    const navigate = useNavigate();

    const currentEvent = '' + (appState.app_state.current_event?.id ?? '');

    const shutdown = () => showDialog(
        'Shutting down',
        'You are trying to shutdown the partyhall. Are you sure ?',
        'Shut down',
        async () => sendMessage('SHUTDOWN', null),
    );


    const setNewEvent = (evt: SelectChangeEvent) => {
        const newId = (evt.target.value as unknown as number);
        if (newId === appState.app_state.current_event?.id) {
            return;
        }

        const events = appState.known_events.filter(x => x.id === newId); // wow such typescript
        if (events.length === 0) {
            return;
        }

        if (!appState.app_state.current_event) {
            sendMessage('SET_EVENT', events[0].id);
            return;
        }

        const newEvent = events[0];

        showDialog(
            'Change event',
            `You are updating the current event to "${newEvent.name} (by ${newEvent.author})".
             Doing so will make that all new pictures are sent to this event instead of the current one.`,
            'Change event',
            async () => sendMessage('SET_EVENT', newEvent.id),
        );
    };

    return <>
        <Card>
            <CardContent>
                <Typography variant="h2" fontSize={18}>Current event</Typography>
                {
                    appState.known_events.length > 0
                    &&
                    <Select value={currentEvent} label="Event" onChange={setNewEvent} style={{ marginTop: '1em' }}>
                        {
                            appState.known_events.map(x => <MenuItem key={x.id} value={x.id}>{x.name}</MenuItem>)
                        }
                    </Select>
                }
            </CardContent>
            <CardActions style={{ justifyContent: 'center' }}>
                <IconButton color="primary" onClick={() => navigate('/admin/event/edit')}>
                    <AddIcon />
                </IconButton>

                {
                    currentEvent != ''
                    && <IconButton color="warning" onClick={() => navigate(`/admin/event/edit/${currentEvent}`)}>
                        <EditIcon />
                    </IconButton>
                }
            </CardActions>
        </Card>
        {
            !!appState.known_modes && appState.known_modes.length > 0
            && <Card>
                <CardContent>
                    <Typography variant="h2" fontSize={18}>Mode</Typography>
                    {
                        appState?.current_mode &&
                        <Select value={appState.current_mode} label="Mode" onChange={(evt: SelectChangeEvent) => sendMessage('SET_MODE', evt.target.value)} style={{ marginTop: '1em' }}>
                            {
                                appState.known_modes.map(x => <MenuItem key={x} value={x}>{x}</MenuItem>)
                            }
                        </Select>
                    }
                </CardContent>
            </Card>
        }
        <Card>
            <CardContent>
                <Typography variant="h2" fontSize={18}>System info</Typography>
                <Box mt={2}>
                    <Typography variant="body1" color="GrayText">PartyHall {appState.partyhall_version} ({appState.partyhall_commit})</Typography>
                    <Typography variant="body1" color="GrayText">Current time: {currentTime}</Typography>
                </Box>
            </CardContent>
            <CardActions>
                <Button style={{ width: '100%' }} onClick={() => sendMessage('SET_DATETIME', DateTime.now().toFormat('yyyy-MM-dd HH:mm:ss'))}>Set to my device's time</Button>
            </CardActions>
        </Card>

        <Card>
            <CardActions>
                <Button style={{ width: '100%' }} onClick={() => sendMessage('DISPLAY_DEBUG', null)}>Show debug info (30 sec)</Button>
            </CardActions>
        </Card>

        <Card>
            <CardActions>
                <Button style={{ width: '100%' }} color="error" onClick={shutdown}>Shutdown</Button>
            </CardActions>
        </Card>
    </>
}