import { Button, Flex, Slider, Tooltip, Typography, notification } from 'antd';
import { IconSun, IconSunOff } from '@tabler/icons-react';
import { useEffect, useState } from 'react';
import BackdropSelector from '../components/backdrop_selector';
import KeyVal from '../components/keyval';
import Title from 'antd/es/typography/Title';
import { useAuth } from '../hooks/auth';
import { useSettings } from '../hooks/settings';
import { useTranslation } from 'react-i18next';

export default function Photobooth() {
    const { t: tGen } = useTranslation('', { keyPrefix: 'generic' });
    const { t } = useTranslation('', { keyPrefix: 'photobooth' });

    const [notif, ctxHolder] = notification.useNotification();
    const { setPageName } = useSettings();
    const { event, api, hardwareFlash, backdropAlbum, selectedBackdrop } = useAuth();

    const [hwFlashBrightness, setHwFlashBrightness] = useState<number>(hardwareFlash.brightness);

    useEffect(() => setPageName('photobooth'), []);

    const flashChange = async (powered: boolean, brightness: number) => {
        setHwFlashBrightness(brightness);
        await api.photobooth.setFlash(powered, brightness);
    };

    const changeSelectedBackdrop = async (newIdx: number) => {
        if (!backdropAlbum) {
            return;
        }

        if (newIdx < 0) {
            newIdx = backdropAlbum.backdrops.length;
        } else if (newIdx > backdropAlbum.backdrops.length) {
            newIdx = 0;
        }

        await api.state.setBackdrops(backdropAlbum.id, newIdx);
    };

    return (
        <Flex vertical gap={8}>
            <Typography>
                <Title>
                    <span className="green">{t('title')}</span>
                </Title>
            </Typography>
            {event && (
                <>
                    <KeyVal label={t('amt_manually_taken')}>{event.amtImagesHandtaken}</KeyVal>
                    <KeyVal label={t('amt_unattended')}>{event.amtImagesUnattended}</KeyVal>

                    <span className="red">{t('backdrops')}: </span>
                    <BackdropSelector />

                    {backdropAlbum && (
                        <Flex vertical gap={4} align="center">
                            <Typography.Text>
                                {t('selected_backdrop')}:{' '}
                                {selectedBackdrop === 0
                                    ? tGen('none')
                                    : backdropAlbum.backdrops[selectedBackdrop - 1].title}
                            </Typography.Text>
                            {selectedBackdrop > 0 && (
                                <img
                                    src={api.backdrops.getImageLink(backdropAlbum.backdrops[selectedBackdrop - 1].id)}
                                    alt={`Current backdrop: ${backdropAlbum.backdrops[selectedBackdrop - 1].title}`}
                                    style={{
                                        maxWidth: '150px',
                                        maxHeight: '150px',
                                    }}
                                />
                            )}
                            <Flex gap={4} style={{ width: '100%' }}>
                                <Button
                                    style={{ width: '100%' }}
                                    onClick={() => changeSelectedBackdrop(selectedBackdrop - 1)}
                                >
                                    {tGen('actions.previous')}
                                </Button>
                                <Button
                                    style={{ width: '100%' }}
                                    onClick={() => changeSelectedBackdrop(selectedBackdrop + 1)}
                                >
                                    {tGen('actions.next')}
                                </Button>
                            </Flex>
                        </Flex>
                    )}

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
