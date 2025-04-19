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
                <LinkButton to="/settings/button-mappings">{t('settings.btn_mappings.title')}</LinkButton>
                <LinkButton to="/settings/nexus">{t('settings.nexus.title')}</LinkButton>
                <span>Users management</span>
                <LinkButton to="/settings/wireless-ap">{t('settings.wireless_ap.title')}</LinkButton>
                <LinkButton to="/settings/third-party">{t('settings.third_party')}</LinkButton>
            </Flex>
        </Flex>
    );
}
