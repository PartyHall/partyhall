import { Flex, Typography } from "antd";
import LinkButton from "../../../components/link_button";
import { useEffect } from "react";
import { useSettings } from "../../../hooks/settings";
import { useTranslation } from "react-i18next";

export default function SettingsPage() {
    const { t } = useTranslation();
    const { setPageName } = useSettings();

    useEffect(() => {
        setPageName('settings');
    }, []);

    return <Flex vertical style={{ width: '400px' }} gap="2em">
        <Typography.Title>{t('settings.title')}</Typography.Title>
        <Flex vertical gap={16} align="center">
            <LinkButton to="/settings/photobooth">{t('settings.photobooth.title')}</LinkButton>
            <LinkButton to="/settings/audio">{t('settings.audio.title')}</LinkButton>
            <LinkButton to="/settings/button-mappings">Physical button mappings</LinkButton>
            <LinkButton to="/settings/nexus">PartyNexus</LinkButton>
            <LinkButton to="/settings/users">Users management</LinkButton>
            <LinkButton to="/settings/wifi-ap">Wifi access point</LinkButton>
            <LinkButton to="/settings/3rd-party">Third-party integrations</LinkButton>
        </Flex>
    </Flex>
}