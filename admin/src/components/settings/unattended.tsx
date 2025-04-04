import { Checkbox, Col, Flex, InputNumber, Row, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useSettings } from '../../hooks/settings';
import { useTranslation } from 'react-i18next';

export type SettingsUnattendedValues = {
    enabled: boolean;
    interval: number;
};

type Props = {
    showTitle?: boolean;
    onSettingsChanged: (values: SettingsUnattendedValues) => void;
};

export default function SettingsUnattended({ showTitle, onSettingsChanged }: Props) {
    const { userSettings } = useSettings();

    const [enabled, setEnabled] = useState<boolean>(userSettings?.photobooth.unattended.enabled ?? false);
    const [interval, setInterval] = useState<number>(userSettings?.photobooth.unattended.interval ?? 300);

    const { t } = useTranslation();

    useEffect(() => {
        if (userSettings?.photobooth.resolution) {
            setEnabled(userSettings?.photobooth.unattended.enabled ?? false);
            setInterval(userSettings?.photobooth.unattended.interval ?? 300);
        }
    }, [userSettings]);

    return (
        <Flex vertical gap={8}>
            {showTitle && (
                <Typography.Title level={3} style={{ margin: 0 }}>
                    {t('settings.unattended.title')}
                </Typography.Title>
            )}
            <Typography.Paragraph>{t('settings.unattended.desc')}</Typography.Paragraph>

            <Flex>
                <Checkbox
                    checked={enabled}
                    onChange={(x) => {
                        setEnabled(x.target.checked);
                        onSettingsChanged({ enabled: x.target.checked, interval });
                    }}
                >
                    {t('generic.enabled')} ?
                </Checkbox>
            </Flex>

            <Row gutter={8} align="middle">
                <Col>
                    <Typography.Text>{t('settings.unattended.interval')}:</Typography.Text>
                </Col>
                <Col flex="auto">
                    <InputNumber
                        style={{ width: '100%' }}
                        value={interval}
                        onChange={(x) => {
                            if (!x) {
                                return;
                            }

                            setInterval(x);
                            onSettingsChanged({ enabled, interval: x });
                        }}
                    />
                </Col>
            </Row>
        </Flex>
    );
}
