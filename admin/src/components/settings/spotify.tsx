import { Checkbox, Col, Flex, Input, Row, Typography } from 'antd';
import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../../hooks/auth';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

export type SettingsSpotifyValues = {
    enabled: boolean;
    name: string;
};

type Props = {
    showTitle?: boolean;
    onSettingsChanged: (values: SettingsSpotifyValues) => void;
    error?: string | null;
};

export default function SettingsSpotify({ showTitle, onSettingsChanged, error }: Props) {
    const { api } = useAuth();
    const [enabled, setEnabled] = useState<boolean>(false);
    const [name, setName] = useState<string>('');

    const { t } = useTranslation();

    useAsyncEffect(async () => {
        const data = await api.settings.getSpotify();
        if (data) {
            setEnabled(data.enabled);
            setName(data.name);
        }
    }, []);

    return (
        <Flex vertical gap={8}>
            {showTitle && (
                <Typography.Title level={3} style={{ margin: 0 }}>
                    {t('settings.spotify.title')}
                </Typography.Title>
            )}
            <Typography.Paragraph>{t('settings.spotify.desc')}</Typography.Paragraph>

            <Row gutter={8} align="middle">
                <Col flex="auto">
                    <Checkbox
                        checked={enabled}
                        onChange={(x) => {
                            setEnabled(x.target.checked);
                            onSettingsChanged({ enabled: x.target.checked, name });
                        }}
                    >
                        {t('generic.enabled')} ?
                    </Checkbox>
                </Col>
            </Row>

            <Row gutter={8} align="middle">
                <Col>
                    <Typography.Text>{t('settings.spotify.name')}:</Typography.Text>
                </Col>
                <Col flex="auto">
                    <Input
                        value={name}
                        onChange={(x) => {
                            setName(x.target.value);
                            onSettingsChanged({ enabled, name: x.target.value });
                        }}
                    />
                </Col>
            </Row>

            <Row gutter={8} align="middle">
                <Col>
                    <Typography.Text>{t('settings.spotify.name_desc')}</Typography.Text>
                </Col>
            </Row>

            {error && (
                <Row gutter={8} align="middle">
                    <Col>
                        <Typography.Text style={{ color: '#ff7875' }}>{error}</Typography.Text>
                    </Col>
                </Row>
            )}
        </Flex>
    );
}
