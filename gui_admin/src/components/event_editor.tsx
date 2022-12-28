import React, { useEffect, useState } from "react";
import { Button, Dialog, DialogContent, DialogTitle, TextField } from "@mui/material";
import { DateTimePicker, LocalizationProvider } from "@mui/x-date-pickers";
import { AdapterLuxon } from "@mui/x-date-pickers/AdapterLuxon";
import { DateTime } from "luxon";
import { EditedEvent } from "../types/appstate";
import { useApi } from "../hooks/auth";

import '../assets/css/event_editor.scss';

type Props = {
    event: EditedEvent;
    hide: () => void;
};

export default function EventEditor({ event, hide }: Props) {
    const { password, showError } = useApi();
    const [evt, setEvent] = useState<EditedEvent>(event);

    const [date, setDate] = useState<DateTime | null>(null);
    useEffect(() => {
        if (!event.date) {
            setDate(null);
        } else {
            setDate(DateTime.fromSeconds(event.date));
        }
    }, [event]);

    useEffect(() => {
        if (!date) {
            setEvent({ ...evt, date: undefined });
        } else {
            setEvent({ ...evt, date: date.toSeconds() })
        }
    }, [date]);

    const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        const query = {
            name: evt.name,
            author: evt.author,
            location: evt.location,
            date: evt.date,
        };

        (async () => {
            try {
                const resp = await fetch(
                    '/api/admin/event' + (!!event.id ? `/${event.id}` : ''),
                    {
                        method: !!event.id ? 'PUT' : 'POST',
                        headers: {
                            'Authorization': password ?? '',
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(query),
                    }
                );

                if (resp.status !== 200) {
                    throw await resp.json();
                }

                hide();
            } catch (e) {
                showError('An error has occured: ' + e, 'error');
            }
        })();

        return false;
    };

    const title = (!evt.id ? "Creating " : "Editing ") + ((evt.name && evt.name.length > 0) ? evt.name : "an event");

    return <Dialog open={true} onClose={hide}>
        <DialogTitle>{title}</DialogTitle>
        <DialogContent>
            <form className="EventEditorForm" onSubmit={onSubmit}>
                <TextField label="Name" value={evt.name ?? ''} onChange={(e) => setEvent({ ...evt, name: e.target.value })} />
                <TextField label="Author" value={evt.author ?? ''} onChange={(e) => setEvent({ ...evt, author: e.target.value })} />
                <TextField label="Location" value={evt.location ?? ''} onChange={(e) => setEvent({ ...evt, location: e.target.value })} />
                <LocalizationProvider dateAdapter={AdapterLuxon}>
                    <DateTimePicker
                        renderInput={(props) => <TextField {...props} />}
                        label="Date"
                        value={date}
                        disableMaskedInput
                        onChange={(newValue) => setDate(newValue)}
                    />
                </LocalizationProvider>

                <div className="EventEditorForm__ActionButtons">
                    <Button onClick={hide}>Cancel</Button>
                    <Button type="submit">Save</Button>
                </div>
            </form>
        </DialogContent>
    </Dialog>
}