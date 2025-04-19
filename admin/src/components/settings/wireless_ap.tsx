import { Checkbox, Col, Flex, Input, Row, Select, Spin, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { InterfacesSettings } from '@partyhall/sdk/dist/models/interfaces';
import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../../hooks/auth';
import { useTranslation } from 'react-i18next';

export type SettingsWirelessApValues = {
    ethIface: string;
    wifiIface: string;
    enabled: boolean;
    name: string;
    password: string;
};

type Props = {
    showTitle?: boolean;
    onSettingsChanged: (values: SettingsWirelessApValues) => void;
    error?: string | null;
};

export default function SettingsWirelessAp({ showTitle, onSettingsChanged, error }: Props) {
    const { t } = useTranslation();
    const { api } = useAuth();

    const [loading, setLoading] = useState<boolean>(true);
    const [apiSettings, setApiSettings] = useState<InterfacesSettings | null>(null);
    const [settings, setSettings] = useState<SettingsWirelessApValues>({
        ethIface: '',
        wifiIface: '',
        enabled: false,
        name: '',
        password: '',
    });

    useAsyncEffect(async () => {
        setLoading(true);

        const data = await api.settings.getWirelessAp();
        setApiSettings(data);
        if (!data) {
            return;
        }

        const apiSettings = data.accessPointSettings;
        setSettings({
            ethIface: apiSettings.wiredInterface,
            wifiIface: apiSettings.wirelessInterface,
            enabled: apiSettings.enabled,
            name: apiSettings.ssid,
            password: apiSettings.password,
        });

        setLoading(false);
    }, []);

    useEffect(() => onSettingsChanged(settings), [settings]);

    return (
        <Flex vertical gap={8} style={{ maxWidth: '50ch' }}>
            {showTitle && (
                <Typography.Title level={3} style={{ margin: 0 }}>
                    {t('settings.wireless_ap.title')}
                </Typography.Title>
            )}
            <Typography.Paragraph>{t('settings.wireless_ap.desc')}</Typography.Paragraph>
            <Typography.Paragraph>{t('settings.wireless_ap.desc_2')}</Typography.Paragraph>
            <Typography.Paragraph>{t('settings.wireless_ap.desc_3')}</Typography.Paragraph>

            <Spin spinning={loading}>
                <Flex vertical gap={8}>
                    {error && (
                        <Typography.Paragraph style={{ textAlign: 'center', color: '#ff7875' }}>
                            {error}
                        </Typography.Paragraph>
                    )}

                    <Row gutter={8} align="middle">
                        <Col>
                            <Typography.Text>{t('settings.wireless_ap.wired_iface')}:</Typography.Text>
                        </Col>
                        <Col flex="auto">
                            <Select
                                style={{ width: '100%' }}
                                options={[
                                    { value: '', label: '-- No interface --' },
                                    ...(apiSettings?.wiredInterfaces.map((x) => ({
                                        value: x.name,
                                        label: x.friendlyName,
                                    })) ?? []),
                                ]}
                                value={settings.ethIface}
                                onChange={(x) => setSettings((oldSettings) => ({ ...oldSettings, ethIface: x }))}
                            />
                        </Col>
                    </Row>

                    <Row gutter={8} align="middle">
                        <Col flex="auto">
                            <Checkbox
                                checked={settings.enabled}
                                onChange={(x) =>
                                    setSettings((oldSettings) => ({ ...oldSettings, enabled: x.target.checked }))
                                }
                            >
                                {t('generic.enabled')} ?
                            </Checkbox>
                        </Col>
                    </Row>

                    <Row gutter={8} align="middle">
                        <Col>
                            <Typography.Text>{t('settings.wireless_ap.wireless_iface')}:</Typography.Text>
                        </Col>
                        <Col flex="auto">
                            <Select
                                disabled={!settings.enabled}
                                style={{ width: '100%' }}
                                value={settings.wifiIface}
                                options={[
                                    { value: '', label: '-- No interface --' },
                                    ...(apiSettings?.wirelessInterfaces.map((x) => ({
                                        value: x.name,
                                        label: x.friendlyName,
                                    })) ?? []),
                                ]}
                                onChange={(x) => setSettings((oldSettings) => ({ ...oldSettings, wifiIface: x }))}
                            />
                        </Col>
                    </Row>

                    <Row gutter={8} align="middle">
                        <Col>
                            <Typography.Text>{t('settings.wireless_ap.ap_name')}:</Typography.Text>
                        </Col>
                        <Col flex="auto">
                            <Input
                                disabled={!settings.enabled}
                                value={settings.name}
                                onChange={(x) =>
                                    setSettings((oldSettings) => ({ ...oldSettings, name: x.target.value }))
                                }
                            />
                        </Col>
                    </Row>

                    <Row gutter={8} align="middle">
                        <Col>
                            <Typography.Text>{t('settings.wireless_ap.ap_password')}:</Typography.Text>
                        </Col>
                        <Col flex="auto">
                            <Input.Password
                                disabled={!settings.enabled}
                                value={settings.password}
                                onChange={(x) =>
                                    setSettings((oldSettings) => ({ ...oldSettings, password: x.target.value }))
                                }
                            />
                        </Col>
                    </Row>

                    <Row gutter={8} align="middle">
                        <Col>
                            <Typography.Text>{t('settings.wireless_ap.ap_password_caution')}</Typography.Text>
                        </Col>
                    </Row>
                </Flex>
            </Spin>
        </Flex>
    );
}
