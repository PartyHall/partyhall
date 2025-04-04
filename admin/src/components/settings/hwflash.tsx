import { Col, Flex, Row, Slider, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { HardwareFlashToggle } from '../hwflash_toggle';
import { useSettings } from '../../hooks/settings';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../../hooks/auth';

export type SettingsHardwareFlashValues = {
    powered: boolean;
    brightness: number;
};

type Props = {
    showTitle?: boolean;
    onSettingsChanged: (values: SettingsHardwareFlashValues) => void;
    setPowered?: (powered: boolean) => void;
};

export default function SettingsHardwareFlash({ showTitle, onSettingsChanged }: Props) {
    const { t } = useTranslation();
    const { userSettings } = useSettings();
    const { hardwareFlashPowered } = useAuth();

    const [brightness, setBrightness] = useState<number>(userSettings?.photobooth.flashBrightness ?? 100);

    useEffect(() => {
        setBrightness(userSettings?.photobooth.flashBrightness ?? 100)
    }, [userSettings]);

    return (
        <Flex vertical gap={8}>
            {showTitle && (
                <Typography.Title level={3} style={{ margin: 0 }}>
                    {t('settings.hwflash.title')}
                </Typography.Title>
            )}
            <Typography.Paragraph>{t('settings.hwflash.desc')}</Typography.Paragraph>
            <Typography.Paragraph>{t('settings.hwflash.desc2')}</Typography.Paragraph>

            <Row gutter={16} align="middle">
                <Col>
                    <Typography.Text>{t('settings.hwflash.brightness')}:</Typography.Text>
                </Col>
                <Col flex="auto">
                    <Slider
                        min={0}
                        max={100}
                        value={brightness}
                        tooltip={{ open: false }}
                        onChange={(x) => {
                            setBrightness(x);
                            onSettingsChanged({ powered: hardwareFlashPowered, brightness: x });
                        }}
                    />
                </Col>
                <Col>
                    <Typography.Text>{brightness}%</Typography.Text>
                </Col>
                <Col>
                    <HardwareFlashToggle />
                </Col>
            </Row>
        </Flex>
    );
}
