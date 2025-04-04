import { Button, Card, Flex } from 'antd';
import SettingsSpotify, { SettingsSpotifyValues } from '../../../components/settings/spotify';
import { IconDeviceFloppy } from '@tabler/icons-react';
import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../../../hooks/auth';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

// Quel language de merde
const isObj = (x: any) => Object.prototype.toString.call(x) === '[object Object]';

export default function SettingsThirdPartyPage() {
    const { t } = useTranslation();
    const { api } = useAuth();

    const [spotifySettings, setSpotifySettings] = useState<SettingsSpotifyValues | null>(null);
    const [spotifyError, setSpotifyError] = useState<string | null>(null);

    const onSpotifySave = async () => {
        if (!spotifySettings) {
            return;
        }

        try {
            await api.settings.setSpotify(spotifySettings.enabled, spotifySettings.name);
            setSpotifySettings(null);
        } catch (e: any) {
            if (isObj(e.message) && e.message.extra_data?.err) {
                setSpotifyError(e.message.extra_data?.err);
            } else {
                setSpotifyError('Unknown error occured, check the console for more details');
            }
        }
    };

    useAsyncEffect(async () => setSpotifySettings(null), []);

    return (
        <Flex vertical gap={16} style={{ maxWidth: '60ch' }}>
            <Card
                title={t('settings.spotify.title')}
                extra={
                    <Button
                        type="primary"
                        icon={<IconDeviceFloppy size={18} />}
                        disabled={!spotifySettings}
                        onClick={onSpotifySave}
                    />
                }
            >
                <SettingsSpotify onSettingsChanged={setSpotifySettings} error={spotifyError} />
            </Card>
        </Flex>
    );
}
