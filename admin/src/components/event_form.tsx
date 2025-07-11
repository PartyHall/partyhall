import { Button, Checkbox, Flex, Form, Input, notification } from 'antd';
import { DateTime } from 'luxon';
import { FormItem } from 'react-hook-form-antd';
import { PhEvent } from '@partyhall/sdk';

import generatePicker from 'antd/es/date-picker/generatePicker';
import luxonGenerateConfig from 'rc-picker/lib/generate/luxon';

import { useAuth } from '../hooks/auth';
import { useForm } from 'react-hook-form';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

const DatePicker = generatePicker<DateTime>(luxonGenerateConfig);

type EventFormData = {
    name: string;
    author: string;
    date: DateTime;
    location: string;
    nexusId: string;
    registrationUrl: string;
    displayText: string | null;
    displayTextAppliance: boolean;
};

type Props = {
    event: PhEvent | null;
    onSaved?: (event: PhEvent) => void;
};

export default function EventForm({ event, onSaved }: Props) {
    const { t } = useTranslation('', { keyPrefix: 'events.editor' });
    const { t: tG } = useTranslation('', { keyPrefix: 'generic' });

    const [notif, ctxHolder] = notification.useNotification();
    const { api, setEvent } = useAuth();

    const [isCreatingNexus, setCreatingNexus] = useState<boolean>(false);

    console.log('EventForm', event);
    const { control, handleSubmit, setValue } = useForm<EventFormData>({
        defaultValues: {
            name: event?.name || '',
            author: event?.author || '',
            date: event?.date || DateTime.now(),
            location: event?.location || '',
            nexusId: event?.nexusId || '',
            registrationUrl: event?.registrationUrl || '',
            displayText: event?.displayText || '',
            displayTextAppliance: event?.displayTextAppliance || false,
        },
    });

    const submit = async (data: EventFormData) => {
        // Yeah i'm not ready to put react-hook-resolvers
        // just because react-hook-form-antd does not
        // support validation rules
        const hasName = data.name.trim().length > 0;
        if (!hasName) {
            notif.error({
                message: 'Invalid event',
                description: 'The event is missing a name',
            });
            return;
        }

        let formEvent = event;
        if (!formEvent) {
            formEvent = new PhEvent({
                id: null,
                name: data.name,
                author: data.author,
                date: data.date,
                location: data.location,
                nexus_id: data.nexusId,
                registration_url: data.registrationUrl,
                display_text: data.displayText,
                display_text_appliance: data.displayTextAppliance,
            });
        } else {
            formEvent.name = data.name;
            formEvent.author = data.author;
            formEvent.date = data.date;
            formEvent.location = data.location;
            formEvent.nexusId = data.nexusId;
            formEvent.registrationUrl = data.registrationUrl;
            formEvent.displayText = data.displayText;
            formEvent.displayTextAppliance = data.displayTextAppliance;
        }

        const isEdit = !!formEvent.id;
        if (isEdit) {
            formEvent = await api.events.update(formEvent);
        } else {
            formEvent = await api.events.create(formEvent);
        }

        if (onSaved && formEvent) {
            onSaved(formEvent);
        }

        notif.success({
            message: (isEdit ? 'Editing ' : 'Creating ') + formEvent?.name,
            description: 'The event was ' + (isEdit ? 'edited' : 'created'),
        });
    };

    const createOnPartyNexus = async () => {
        if (!event || !event.id) {
            return;
        }

        setCreatingNexus(true);

        try {
            const resp = await api.nexus.createEvent(event.id);
            if (resp) {
                setEvent(resp);
                setValue('nexusId', resp.nexusId ?? '');
            }
        } catch (e) {
            // @TODO: Handle properly
            console.error(e);
            notif.error({
                message: 'Failed to create event on PartyNexus',
                description: 'The event was not created due to an issue. See console for more detail',
            });
        }

        setCreatingNexus(false);
        notif.success({
            message: 'PartyNexus event created',
            description: 'The event was created on PartyNexus.',
        });
    };

    return (
        <Form layout="vertical" style={{ width: 500 }} variant="filled" onFinish={handleSubmit(submit)}>
            <FormItem control={control} name="name" label={tG('name')} required>
                <Input />
            </FormItem>
            <FormItem control={control} name="author" label={tG('author')}>
                <Input />
            </FormItem>
            <FormItem control={control} name="date" label={tG('date')} required>
                <DatePicker style={{ width: '100%' }} format="DD/MM/YYYY hh:mm" showTime />
            </FormItem>
            <FormItem control={control} name="location" label={tG('location')}>
                <Input />
            </FormItem>

            <Flex align="center" gap={8}>
                <FormItem control={control} name="nexusId" label="Nexus ID" style={{ flex: '1' }}>
                    <Input />
                </FormItem>

                <Button
                    style={{ marginTop: '6px' }}
                    onClick={createOnPartyNexus}
                    disabled={isCreatingNexus || !event?.id}
                >
                    {t('create_on_pn')}
                </Button>
            </Flex>

            <FormItem control={control} name="registrationUrl" label="Registration URL">
                <Input />
            </FormItem>

            <FormItem control={control} name="displayText" label="Display Text">
                <Input />
            </FormItem>

            <FormItem
                control={control}
                name="displayTextAppliance"
                label="Display Text on the appliance"
                valuePropName="checked"
            >
                <Checkbox />
            </FormItem>

            <Form.Item>
                <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
                    {tG('actions.' + (event?.id ? 'save' : 'create'))}
                </Button>
            </Form.Item>

            {ctxHolder}
        </Form>
    );
}
