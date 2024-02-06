import { Box, Button, Card, CardActions, CardContent, IconButton, MenuItem, Select, SelectChangeEvent, Typography } from "@mui/material";
import { DateTime } from 'luxon';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';

import { useConfirmDialog } from "../../hooks/dialog";
import { useAdminSocket } from "../../hooks/adminSocket";
import { useNavigate } from "react-router-dom";
import { useApi } from "../../hooks/useApi";
import { useTranslation } from "react-i18next";

export default function AdminIndex() {
    const { hasRole } = useApi();
    const { t } = useTranslation();
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

        const events = appState.known_events.filter(x => x.id === newId);
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
                <Typography variant="h2" fontSize={18}>{t('admin_main.current_event')}</Typography>
                {
                    appState.known_events.length > 0
                    &&
                    <Select value={currentEvent} label="Event" onChange={setNewEvent} style={{ marginTop: '1em' }} disabled={!hasRole('ADMIN')}>
                        {
                            appState.known_events.map(x => <MenuItem key={x.id} value={x.id}>{x.name}</MenuItem>)
                        }
                    </Select>
                }
            </CardContent>
            <CardActions style={{ justifyContent: 'center' }}>
                {
                    hasRole('ADMIN') && <>
                        <IconButton color="primary" onClick={() => navigate('/admin/event/edit')}>
                            <AddIcon />
                        </IconButton>

                        {
                            currentEvent != ''
                            && <IconButton color="warning" onClick={() => navigate(`/admin/event/edit/${currentEvent}`)}>
                                <EditIcon />
                            </IconButton>
                        }
                    </>
                }
            </CardActions>
        </Card>
        {
            !!appState.known_modes && appState.known_modes.length > 0
            && <Card>
                <CardContent>
                    <Typography variant="h2" fontSize={18}>{t('admin_main.mode')}</Typography>
                    {
                        appState?.current_mode &&
                        <Select value={appState.current_mode} label="Mode" onChange={(evt: SelectChangeEvent) => sendMessage('SET_MODE', evt.target.value)} style={{ marginTop: '1em' }} disabled={!hasRole('ADMIN')}>
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
                <Typography variant="h2" fontSize={18}>{t('admin_main.system_info')}</Typography>
                <Box mt={2}>
                    <Typography variant="body1" color="GrayText">PartyHall {appState.partyhall_version} ({appState.partyhall_commit})</Typography>
                    <Typography variant="body1" color="GrayText">{t('admin_main.current_time')}: {currentTime}</Typography>
                </Box>
            </CardContent>
            <CardActions>
                {
                    hasRole('ADMIN') &&
                    <Button style={{ width: '100%' }} onClick={() => sendMessage('SET_DATETIME', DateTime.now().toFormat('yyyy-MM-dd HH:mm:ss'))}>{t('admin_main.set_to_my_time')}</Button>
                }
            </CardActions>
        </Card>

        {
            hasRole('ADMIN') &&
            <Card>
                <CardActions>
                    <Button style={{ width: '100%' }} onClick={() => sendMessage('DISPLAY_DEBUG', null)}>{t('admin_main.show_debug_info')}</Button>
                </CardActions>
            </Card>
        }

        {
            hasRole('ADMIN') &&
            <Card>
                <CardActions>
                    <Button style={{ width: '100%' }} color="error" onClick={shutdown}>{t('admin_main.shutdown')}</Button>
                </CardActions>
            </Card>
        }
    </>
}