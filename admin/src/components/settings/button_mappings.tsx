import { Col, Flex, Row, Select, Spin, Typography } from 'antd';
import { useEffect, useState } from 'react';
import useAsyncEffect from 'use-async-effect';
import { useAuth } from '../../hooks/auth';
import { useMercureTopic } from '../../hooks/mercure';
import { useTranslation } from 'react-i18next';

export type SettingsButtonMappingsValues = Record<number, string>;

type Props = {
    showTitle?: boolean;
    onSettingsChanged: (values: SettingsButtonMappingsValues) => void;
};

export default function SettingsButtonMappings({ showTitle, onSettingsChanged }: Props) {
    const { t } = useTranslation();
    const { api } = useAuth();

    const [loading, setLoading] = useState<boolean>(true);

    const [availableActions, setAvailableActions] = useState<string[]>([]);

    const [lastButtonPress, setLastButtonPress] = useState<number | null>(null);
    const [mappings, setMappings] = useState<SettingsButtonMappingsValues>({});

    useMercureTopic('/btn-press', (x) => setLastButtonPress(x as number));

    useAsyncEffect(async () => {
        setLoading(true);
        setMappings(await api.settings.getButtonMappings());
        setAvailableActions(await api.settings.getButtonMappingsActions());
        setLoading(false);
    }, []);

    // This pings the appliance to keep the "btn_setup" mode active
    // while we're on this page
    useEffect(() => {
        const itv = setInterval(() => api.settings.getButtonMappings(), 2000);

        return () => {
            clearInterval(itv);
        };
    }, []);

    useEffect(() => {
        if (lastButtonPress === null) {
            return;
        }

        if (Object.keys(mappings).includes(`${lastButtonPress}`)) {
            return;
        }

        setMappings({
            ...mappings,
            [lastButtonPress]: '',
        });
    }, [lastButtonPress]);

    if (loading) {
        return <Spin spinning />;
    }

    return (
        <Flex vertical>
            {showTitle && (
                <Typography.Title level={3} style={{ margin: 0 }}>
                    {t('settings.btn_mappings.title')}
                </Typography.Title>
            )}

            <Typography.Paragraph>{t('settings.btn_mappings.desc')}</Typography.Paragraph>
            <Typography.Paragraph>{t('settings.btn_mappings.desc2')}</Typography.Paragraph>

            <Flex vertical gap={8}>
                {mappings &&
                    Object.keys(mappings).map((x) => (
                        <Row gutter={8} align="middle" key={x}>
                            <Col flex="30%">
                                <Typography.Text className={`${lastButtonPress}` === x ? 'blue' : ''}>
                                    BTN_{x}
                                </Typography.Text>
                            </Col>
                            <Col flex="auto">
                                <Select
                                    style={{ width: '100%' }}
                                    value={mappings[parseInt(x)]}
                                    options={availableActions.map((x) => ({
                                        label: t(`settings.btn_mappings.actions.${x}`),
                                        value: x,
                                    }))}
                                    onChange={(val) => {
                                        const newMappings = {
                                            ...mappings,
                                            [x]: val,
                                        };

                                        setMappings(newMappings);
                                        onSettingsChanged(newMappings);
                                    }}
                                />
                            </Col>
                        </Row>
                    ))}
            </Flex>
        </Flex>
    );
}
