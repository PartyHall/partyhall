import { Button, Card, Flex, notification } from 'antd';
import SettingsWirelessAp, { SettingsWirelessApValues } from '../../../components/settings/wireless_ap';
import { IconDeviceFloppy } from '@tabler/icons-react';
import { useAuth } from '../../../hooks/auth';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

export default function WirelessApPage() {
    const { t } = useTranslation();
    const [notif, ctxHolder] = notification.useNotification();

    const { api } = useAuth();

    const [isSubmitting, setSubmitting] = useState(false);
    const [wifiApSettings, setWifiApSettings] = useState<SettingsWirelessApValues | null>(null);
    const [error, setError] = useState<string | null>(null);

    const onSave = async () => {
        if (!wifiApSettings) {
            setError(t('settings.wireless_ap.error_missing_settings'));
            return;
        }

        setSubmitting(true);
        setError(null);

        try {
            await api.settings.setWirelessAp(
                wifiApSettings.ethIface,
                wifiApSettings.wifiIface,
                wifiApSettings.enabled,
                wifiApSettings.name,
                wifiApSettings.password
            );

            notif.success({
                message: t('settings.wireless_ap.settings_saved'),
                description: t('settings.wireless_ap.settings_saved_desc'),
            });
        } catch (e: any) {
            console.error(e);
            setError(t('settings.wireless_ap.settings_save_error'));
        }

        setSubmitting(false);
    };

    return (
        <Flex vertical gap={16}>
            <Card
                title={t('settings.wireless_ap.title')}
                extra={
                    <Button
                        type="primary"
                        icon={<IconDeviceFloppy size={18} />}
                        disabled={!wifiApSettings || isSubmitting}
                        onClick={onSave}
                    />
                }
            >
                <SettingsWirelessAp onSettingsChanged={setWifiApSettings} error={error} />
            </Card>

            {ctxHolder}
        </Flex>
    );
}
