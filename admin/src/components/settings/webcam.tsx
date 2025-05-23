import { Col, Flex, InputNumber, Row, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useSettings } from '../../hooks/settings';
import { useTranslation } from 'react-i18next';

export type SettingsWebcamValues = {
    width: number;
    height: number;
};

type Props = {
    showTitle?: boolean;
    onSettingsChanged: (values: SettingsWebcamValues) => void;
};

export default function SettingsWebcam({ showTitle, onSettingsChanged }: Props) {
    const { userSettings } = useSettings();

    const [width, setWidth] = useState<number>(userSettings?.photobooth.resolution.width ?? 1920);
    const [height, setHeight] = useState<number>(userSettings?.photobooth.resolution.height ?? 1080);

    const { t } = useTranslation();

    useEffect(() => {
        if (userSettings?.photobooth.resolution) {
            setWidth(userSettings?.photobooth.resolution.width);
            setHeight(userSettings?.photobooth.resolution.height);
        }
    }, [userSettings]);

    return (
        <Flex vertical gap={8}>
            {showTitle && (
                <Typography.Title level={3} style={{ margin: 0 }}>
                    {t('settings.webcam.resolution.subtitle')}
                </Typography.Title>
            )}
            <Typography.Paragraph>{t('settings.webcam.resolution.desc')}</Typography.Paragraph>

            <Row gutter={8} align="middle">
                <Col>
                    <Typography.Text>{t('settings.webcam.resolution.width')}:</Typography.Text>
                </Col>
                <Col flex="auto">
                    <InputNumber
                        style={{ width: '100%' }}
                        value={width}
                        onChange={(x) => {
                            if (!x) {
                                return;
                            }

                            setWidth(x);
                            onSettingsChanged({ width: x, height });
                        }}
                    />
                </Col>
            </Row>

            <Row gutter={8} align="middle">
                <Col>
                    <Typography.Text>{t('settings.webcam.resolution.height')}:</Typography.Text>
                </Col>
                <Col flex="auto">
                    <InputNumber
                        style={{ width: '100%' }}
                        value={height}
                        onChange={(x) => {
                            if (!x) {
                                return;
                            }

                            setHeight(x);
                            onSettingsChanged({ width, height: x });
                        }}
                    />
                </Col>
            </Row>
        </Flex>
    );
}
