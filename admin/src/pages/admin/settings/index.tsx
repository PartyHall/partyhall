import { Flex, Typography } from 'antd';
import LinkButton from '../../../components/link_button';
import { useEffect } from 'react';
import { useSettings } from '../../../hooks/settings';
import { useTranslation } from 'react-i18next';

export default function SettingsPage() {
    const { t } = useTranslation();
    const { setPageName } = useSettings();

    useEffect(() => {
        setPageName('settings');
    }, []);

    return (
        <Flex vertical style={{ width: '400px' }} gap="2em">
            <Typography.Title>{t('settings.title')}</Typography.Title>
            <Flex vertical gap={16} align="center">
                <LinkButton to="/settings/photobooth">{t('settings.photobooth.title')}</LinkButton>
                <LinkButton to="/settings/audio">{t('settings.audio.title')}</LinkButton>
                Physical button mappings
                <LinkButton to="/settings/nexus">PartyNexus</LinkButton>
                Users management Wifi access point
                <LinkButton to="/settings/third-party">Third-party integrations</LinkButton>
            </Flex>
        </Flex>
    );
}
