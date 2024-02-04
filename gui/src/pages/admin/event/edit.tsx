import { useEffect, useMemo, useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { Navigate, useNavigate, useParams } from "react-router-dom"
import { DateTime } from "luxon";

import { Button, Card, CardContent, Grid, Input, TextField, Typography } from "@mui/material";
import { DateTimePicker, LocalizationProvider } from "@mui/x-date-pickers";
import { AdapterLuxon } from "@mui/x-date-pickers/AdapterLuxon";

import { Event, EditedEvent } from '../../../types/appstate';
import { useAdminSocket } from "../../../hooks/adminSocket";
import { useApi } from "../../../hooks/useApi";

const getEmptyEvent = (): EditedEvent => ({
    id: '',
    name: '',
    author: '',
    date: DateTime.now(),
    location: '',
});

function getEditedEventFromEvent(event: Event): EditedEvent {
    return {
        id: event.id,
        name: event.name ?? '',
        author: event.author ?? '',
        date: DateTime.fromSeconds(event.date ?? 0),
        location: event.location ?? '',
    };
}

export default function EditEvent() {
    const { appState } = useAdminSocket();
    const { saveEvent } = useApi();
    const { id: eventId } = useParams();
    const [editedEvent, setEditedEvent] = useState<EditedEvent | null>(null);
    const navigate = useNavigate();

    const { handleSubmit, control, reset } = useForm({
        defaultValues: useMemo(() => {
            if (!editedEvent) {
                return getEmptyEvent();
            }

            return editedEvent;
        }, [editedEvent]),
    });

    useEffect(() => {
        if (!eventId) {
            setEditedEvent(null);
            reset(getEmptyEvent());
            return;
        }

        const events = appState.known_events.filter(x => ('' + x.id) === eventId);
        if (events.length > 0) {
            const evt = getEditedEventFromEvent(events[0]);
            setEditedEvent(evt);
            reset(evt);
        }
    }, [eventId]);

    const title = !eventId ? "Creating event" : ("Editing event " + (editedEvent?.name ?? ''));
    const save = async (data: EditedEvent) => {
        const saved = await saveEvent(data);
        if (saved) {
            navigate('/admin');
        }
    };

    return <Grid container spacing={0} direction="column" alignItems="center" justifyContent="center" minHeight="100%">
        <form onSubmit={handleSubmit(save)}>
            <Card variant="outlined" style={{ maxWidth: '20em' }}>
                <CardContent style={{ display: 'flex', flexDirection: 'column', alignItems: 'stretch', gap: '1em' }}>
                    <Typography sx={{ fontSize: 20 }} variant="h1" color="text.secondary" gutterBottom>
                        {title}
                    </Typography>

                    <Controller
                        name="id"
                        control={control}
                        render={({ field }) => <input type="hidden" {...field} />}
                    />

                    <Controller
                        name="name"
                        control={control}
                        render={({ field }) => <TextField
                            label="Name"
                            InputLabelProps={{ required: false }}
                            required
                            {...field}
                        />}
                    />

                    <Controller
                        name="author"
                        control={control}
                        render={({ field }) => <TextField
                            label="Author"
                            InputLabelProps={{ required: false }}
                            required
                            {...field}
                        />}
                    />

                    <Controller
                        name="date"
                        control={control}
                        render={({ field }) => <LocalizationProvider dateAdapter={AdapterLuxon}>
                            <DateTimePicker
                                label="Date"
                                //@ts-ignore
                                renderInput={(props: any) => <TextField {...props} />}
                                InputProps={{ required: true }}
                                disableMaskedInput
                                {...field}
                            />
                        </LocalizationProvider>
                        }
                    />

                    <Controller
                        name="location"
                        control={control}
                        render={({ field }) => <TextField
                            label="Location"
                            {...field}
                        />}
                    />

                    <Button type="submit">Save</Button>
                </CardContent>
            </Card>
        </form>
    </Grid>;
}