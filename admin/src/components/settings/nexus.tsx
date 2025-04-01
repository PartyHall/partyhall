import { Checkbox, Col, Collapse, Flex, Input, Row, Typography } from 'antd';
import NexusSettings from '@partyhall/sdk/dist/models/nexus';
import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../../hooks/auth';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

export type SettingsNexusValues = {
    nexusUrl: string;
    hardwareId: string;
    apiKey: string;
    bypassSsl: boolean;
};

type Props = {
    showTitle?: boolean;
    onSettingsChanged: (values: SettingsNexusValues) => void;
    errorMessage?: string | null;
    initialData?: NexusSettings | null;
};

export default function SettingsNexus({ showTitle, onSettingsChanged, errorMessage, initialData }: Props) {
    const { t } = useTranslation();
    const { api } = useAuth();

    const [nexusUrl, setNexusUrl] = useState<string>('');
    const [hardwareId, setHardwareId] = useState<string>('');
    const [apiKey, setApiKey] = useState<string>('');
    const [bypassSsl, setBypassSsl] = useState<boolean>(false);

    useAsyncEffect(async () => {
        if (initialData) {
            setNexusUrl(initialData.baseUrl);
            setHardwareId(initialData.hardwareId);
            setBypassSsl(initialData.bypassSsl);

            return;
        }

        // Initial fetch data
        const nexusSettings = await api.settings.getNexus();
        if (!nexusSettings) {
            // @TODO: display snackbar
            return;
        }

        setNexusUrl(nexusSettings.baseUrl);
        setHardwareId(nexusSettings.hardwareId);
        setBypassSsl(nexusSettings.bypassSsl);
    }, [initialData]);

    return (
        <Flex vertical gap={8}>
            {showTitle && (
                <Typography.Title level={3} style={{ margin: 0 }}>
                    {t('settings.nexus.title')}
                </Typography.Title>
            )}

            <Typography.Paragraph>{t('settings.nexus.desc')}</Typography.Paragraph>
            <Typography.Paragraph>{t('settings.nexus.desc2')}</Typography.Paragraph>

            <Flex vertical gap={8}>
                <Row gutter={8} align="middle">
                    <Col flex="auto">
                        <Typography.Text>{t('settings.nexus.url')}</Typography.Text>
                    </Col>
                    <Col>
                        <Input
                            value={nexusUrl}
                            onChange={(x) => {
                                setNexusUrl(x.target.value);
                                onSettingsChanged({ nexusUrl: x.target.value, hardwareId, apiKey, bypassSsl });
                            }}
                        />
                    </Col>
                </Row>

                <Row gutter={8} align="middle">
                    <Col flex="auto">
                        <Typography.Text>{t('settings.nexus.hardware_id')}</Typography.Text>
                    </Col>
                    <Col>
                        <Input
                            value={hardwareId}
                            onChange={(x) => {
                                setHardwareId(x.target.value);
                                onSettingsChanged({ hardwareId: x.target.value, nexusUrl, apiKey, bypassSsl });
                            }}
                        />
                    </Col>
                </Row>

                <Row gutter={8} align="middle">
                    <Col flex="auto">
                        <Typography.Text>{t('settings.nexus.api_key')}</Typography.Text>
                    </Col>
                    <Col>
                        <Input
                            value={apiKey}
                            onChange={(x) => {
                                setApiKey(x.target.value);
                                onSettingsChanged({ apiKey: x.target.value, nexusUrl, hardwareId, bypassSsl });
                            }}
                        />
                    </Col>
                </Row>
            </Flex>

            {errorMessage && (
                <Typography.Paragraph style={{ textAlign: 'center', color: '#ff7875' }}>
                    {errorMessage}
                </Typography.Paragraph>
            )}

            <Collapse
                items={[
                    {
                        key: 'advanced_settings',
                        label: t('settings.advanced_settings'),
                        children: (
                            <Flex vertical gap={8}>
                                <Checkbox
                                    checked={bypassSsl}
                                    onChange={(x) => {
                                        setBypassSsl(x.target.checked);
                                        onSettingsChanged({
                                            bypassSsl: x.target.checked,
                                            nexusUrl,
                                            hardwareId,
                                            apiKey,
                                        });
                                    }}
                                >
                                    {t('settings.nexus.bypass_ssl')}
                                </Checkbox>
                            </Flex>
                        ),
                    },
                ]}
            />
        </Flex>
    );
}
