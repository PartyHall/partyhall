import { Button, Card, Flex, Input, Select } from 'antd';

import EventCard from '../components/event_card';
import KeyVal from '../components/keyval';

import { useAuth } from '../hooks/auth';
import { useEffect } from 'react';
import { useMercure } from '../hooks/mercure';
import { useSettings } from '../hooks/settings';
import { useTranslation } from 'react-i18next';

export default function Index() {
    const { t } = useTranslation('', { keyPrefix: 'home' });
    const { t: tG } = useTranslation('', { keyPrefix: 'generic' });

    const { version, commit, setPageName } = useSettings();
    const {
        api,
        displayName,
        setDisplayName,
        mode,
        event,
        logout,
        syncInProgress,
    } = useAuth();
    const { time } = useMercure();

    const changeMode = async (val: string) => {
        await api.settings.setMode(val);
    };

    useEffect(() => {
        setPageName('home', []);
    }, []);

    /** @TODO: Make responsive for < 400 => width=95% or something like that */
    return (
        <Flex vertical style={{ width: '400px' }} gap="2em">
            <Card title={t('about_you.title')}>
                <Flex vertical gap=".5em">
                    <KeyVal label={t('about_you.your_display_name')}>
                        <Input
                            type="text"
                            value={displayName ?? ''}
                            onChange={(x) => setDisplayName(x.target.value)}
                        />
                    </KeyVal>
                </Flex>
            </Card>

            {event && <EventCard event={event} />}

            <Card title={t('actions.title')}>
                <Flex vertical gap={'.25em'}>
                    <KeyVal label={t('actions.ph_version')}>{version}</KeyVal>
                    <KeyVal label={t('actions.ph_commit')}>{commit}</KeyVal>
                    <KeyVal label={t('actions.appliance_time')}>
                        {time && time.toFormat('HH:mm:ss (yyyy-MM-dd)')}
                    </KeyVal>
                    {api.tokenUser?.roles.includes('ADMIN') && (
                        <>
                            <KeyVal label={tG('mode')}>
                                <Select value={mode} onChange={changeMode}>
                                    <Select.Option value="photobooth">
                                        {tG('modes.photobooth')}
                                    </Select.Option>
                                    <Select.Option value="disabled">
                                        {tG('modes.disabled')}
                                    </Select.Option>
                                </Select>
                            </KeyVal>

                            <KeyVal label={t('actions.sync.title')}>
                                {syncInProgress &&
                                    t('actions.sync.in_progress')}
                                {!syncInProgress && t('actions.sync.idle')}
                            </KeyVal>
                        </>
                    )}

                    <Flex
                        gap="1em"
                        style={{ marginTop: '2em' }}
                        wrap="wrap"
                        justify="center"
                    >
                        {api.tokenUser?.roles.includes('ADMIN') && (
                            <>
                                <Button
                                    color="danger"
                                    onClick={() => api.settings.showDebug()}
                                >
                                    {t('actions.show_debug')}
                                </Button>
                                <Button
                                    color="danger"
                                    disabled={syncInProgress}
                                    onClick={() => api.global.forceSync()}
                                >
                                    {t('actions.force_sync')}
                                </Button>
                                <Button color="danger">
                                    {t('actions.shutdown')}
                                </Button>
                            </>
                        )}
                        <Button color="danger" onClick={logout}>
                            {t('actions.logout')}
                        </Button>
                    </Flex>
                </Flex>
            </Card>
        </Flex>
    );
}