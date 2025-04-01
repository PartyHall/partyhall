import { Col, Flex, Row, Slider, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { HardwareFlashToggle } from '../hwflash_toggle';
import { useSettings } from '../../hooks/settings';
import { useTranslation } from 'react-i18next';

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
    const { user_settings } = useSettings();

    const [powered, setPowered] = useState<boolean>(false);
    const [brightness, setBrightness] = useState<number>(user_settings?.photobooth.flashBrightness ?? 100);

    const { t } = useTranslation();

    useEffect(() => {
        if (user_settings?.photobooth.resolution) {
            setBrightness(user_settings?.photobooth.flashBrightness ?? 100);
        }
    }, [user_settings]);

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
                            onSettingsChanged({ powered, brightness: x });
                        }}
                    />
                </Col>
                <Col>
                    <Typography.Text>{brightness}%</Typography.Text>
                </Col>
                <Col>
                    <HardwareFlashToggle setPoweredState={(pw) => setPowered(pw)} />
                </Col>
            </Row>
        </Flex>
    );
}
