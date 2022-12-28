import { useEffect, useState } from "react";

import { Button, Card, CardActions, CardContent, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, MenuItem, Select, SelectChangeEvent, Typography, IconButton } from "@mui/material";
import { DateTime } from 'luxon';

import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';

import { useWebsocket } from "../hooks/ws";
import { EditedEvent, Event } from '../types/appstate';
import App from "./App";
import EventEditor from "../components/event_editor";

export default function PageIndex() {
    const { sendMessage, appState, currentTime } = useWebsocket();

    const [knownEvents, setKnownEvents] = useState<Event[]>([]);
    const [currentEventInst, setCurrentEventInst] = useState<Event | null>(null);
    const [currentEvent, setCurrentEvent] = useState<string>('');
    const [modes, setModes] = useState<string[]>([]);

    const [newEvent, setNewEvent] = useState<Event | null>(null);
    const [shutdown, setShutdown] = useState<boolean>(false);

    const [editedEvent, setEditedEvent] = useState<EditedEvent | null>(null);

    useEffect(() => {
        if (!appState) {
            setModes([]);
            setCurrentEvent('');
            setCurrentEventInst(null);
            setKnownEvents([]);
            return;
        }

        setKnownEvents(appState.known_events);
        setModes(appState.known_modes);
        if (!!appState.app_state?.current_event) {
            setCurrentEventInst(appState.app_state.current_event);
            setCurrentEvent('' + appState.app_state.current_event.id);
        }
    }, [appState]);

    const setMode = (evt: SelectChangeEvent) => sendMessage('SET_MODE', evt.target.value);

    const setEvent = (evt: SelectChangeEvent) => {
        const newId = (evt.target.value as unknown as number);
        if (newId === appState?.app_state.current_event?.id) {
            return;
        }

        const events = knownEvents.filter(x => x.id === newId); // wow such typescript

        if (events.length > 0) {
            if ((currentEvent ?? '') === '') {
                sendMessage('SET_EVENT', events[0].id);
            } else {
                setNewEvent(events[0]);
            }
        }
    };

    const submitDatetime = () => sendMessage('SET_DATETIME', DateTime.now().toFormat('yyyy-MM-dd HH:mm:ss'));

    return <App>
        <Card>
            <CardContent>
                <Typography variant="h2" fontSize={18}>Current event</Typography>
                {
                    knownEvents.length > 0
                    &&
                    <Select value={currentEvent} label="Event" onChange={setEvent} style={{ marginTop: '1em' }}>
                        {
                            knownEvents.map(x => <MenuItem key={x.id} value={x.id}>{x.name}</MenuItem>)
                        }
                    </Select>
                }
            </CardContent>
            <CardActions style={{ justifyContent: 'center' }}>
                <IconButton color="primary" onClick={() => setEditedEvent({})}><AddIcon /></IconButton>
                {
                    currentEvent != ''
                    && <IconButton color="warning" onClick={() => setEditedEvent(currentEventInst)}><EditIcon /></IconButton>
                }
            </CardActions>
        </Card>
        {
            !!modes && modes.length > 0
            && <Card>
                <CardContent>
                    <Typography variant="h2" fontSize={18}>Mode</Typography>
                    {
                        appState?.current_mode &&
                        <Select value={appState.current_mode} label="Mode" onChange={setMode} style={{ marginTop: '1em' }}>
                            {
                                modes.map(x => <MenuItem key={x} value={x}>{x}</MenuItem>)
                            }
                        </Select>
                    }
                </CardContent>
            </Card>
        }
        <Card>
            <CardContent>
                <Typography variant="h2" fontSize={18}>System time</Typography>
                <Typography variant="body1" style={{ textAlign: "center", marginTop: '2em' }}>{currentTime}</Typography>
            </CardContent>
            <CardActions>
                <Button style={{ width: '100%' }} onClick={submitDatetime}>Set to my device's time</Button>
            </CardActions>
        </Card>
        <Card>
            <CardActions>
                <Button style={{ width: '100%' }} onClick={() => sendMessage('SHOW_DEBUG', null)}>Show debug info (30 sec)</Button>
            </CardActions>
        </Card>
        <Card>
            <CardActions>
                <Button style={{ width: '100%' }} color="error" onClick={() => setShutdown(true)}>Shutdown</Button>
            </CardActions>
        </Card>

        <Dialog open={!!newEvent} onClose={() => setNewEvent(null)}>
            <DialogTitle>Change event</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    You are updating the current event to "{newEvent?.name} (by {newEvent?.author})".
                    Doing so will make that all new pictures are sent to this event instead of the current one.
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => setNewEvent(null)}>Cancel</Button>
                <Button onClick={() => {
                    sendMessage('SET_EVENT', newEvent?.id);
                    setNewEvent(null);
                }} color="warning" autoFocus>Change event</Button>
            </DialogActions>
        </Dialog>

        <Dialog open={shutdown} onClose={() => setShutdown(false)}>
            <DialogTitle>Shutting down</DialogTitle>
            <DialogContent>
                <DialogContentText>You are trying to shutdown the partyhall. Are you sure ?</DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => setShutdown(false)}>Cancel</Button>
                <Button onClick={() => {
                    sendMessage('SHUTDOWN', null);
                    setShutdown(false);
                }} color="error" autoFocus>Shutdown</Button>
            </DialogActions>
        </Dialog>

        {
            editedEvent &&
            <EventEditor event={editedEvent} hide={() => setEditedEvent(null)} />
        }
    </App>;
}