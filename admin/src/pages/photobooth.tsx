import { Button, Flex, Slider, Tooltip, Typography, notification } from 'antd';
import { IconSun, IconSunOff } from '@tabler/icons-react';
import { useEffect, useState } from 'react';
import KeyVal from '../components/keyval';
import Title from 'antd/es/typography/Title';
import { useAuth } from '../hooks/auth';
import { useSettings } from '../hooks/settings';
import { useTranslation } from 'react-i18next';

export default function Photobooth() {
    const { t } = useTranslation('', { keyPrefix: 'photobooth' });

    const [notif, ctxHolder] = notification.useNotification();
    const { setPageName } = useSettings();
    const { event, api, hardwareFlash } = useAuth();

    const [hwFlashBrightness, setHwFlashBrightness] = useState<number>(hardwareFlash.brightness);

    useEffect(() => setPageName('photobooth'), []);

    const flashChange = async (powered: boolean, brightness: number) => {
        setHwFlashBrightness(brightness);
        await api.photobooth.setFlash(powered, brightness);
    };

    return (
        <Flex vertical gap={2}>
            <Typography>
                <Title>
                    <span className="green">{t('title')}</span>
                </Title>
            </Typography>
            {event && (
                <>
                    <KeyVal label={t('amt_manually_taken')}>{event.amtImagesHandtaken}</KeyVal>
                    <KeyVal label={t('amt_unattended')}>{event.amtImagesUnattended}</KeyVal>

                    <span className="red">{t('flash')}: </span>
                    <Flex gap={8} align="center">
                        <span>{hwFlashBrightness}%</span>

                        <Slider
                            style={{ flex: 1 }}
                            min={0}
                            max={100}
                            value={hwFlashBrightness}
                            onChange={(x) => flashChange(hardwareFlash.powered, x)}
                        />

                        <Tooltip title={t(hardwareFlash.powered ? 'toggle_off' : 'toggle_on')}>
                            <Button
                                icon={hardwareFlash.powered ? <IconSunOff /> : <IconSun />}
                                onClick={() => flashChange(!hardwareFlash.powered, hwFlashBrightness)}
                            />
                        </Tooltip>
                    </Flex>

                    <Flex style={{ marginTop: '2em' }} vertical>
                        <Button
                            onClick={async () => {
                                await api.photobooth.takePicture();
                                notif.success({
                                    message: t('success_notification.title'),
                                    description: t('success_notification.description'),
                                });
                            }}
                        >
                            {t('remote_take_picture')}
                        </Button>
                    </Flex>
                </>
            )}

            {ctxHolder}
        </Flex>
    );
}
