import { Button, Flex, Tooltip, Typography, notification } from 'antd';
import { IconSun, IconSunOff } from '@tabler/icons-react';
import BackdropSelector from '../components/backdrop_selector';
import KeyVal from '../components/keyval';
import Title from 'antd/es/typography/Title';
import { useAuth } from '../hooks/auth';
import { useEffect } from 'react';
import { useSettings } from '../hooks/settings';
import { useTranslation } from 'react-i18next';

export default function Photobooth() {
    const { t: tGen } = useTranslation('', { keyPrefix: 'generic' });
    const { t } = useTranslation('', { keyPrefix: 'photobooth' });

    const [notif, ctxHolder] = notification.useNotification();
    const { setPageName } = useSettings();
    const { event, api, hardwareFlashPowered, backdropAlbum, selectedBackdrop } = useAuth();

    useEffect(() => setPageName('photobooth'), []);

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

                    <Flex gap={8} align="center" justify='space-between'>
                        <span className="red">{t('flash')}: </span>
                        <Tooltip title={t(hardwareFlashPowered ? 'toggle_off' : 'toggle_on')}>
                            <Button
                                icon={hardwareFlashPowered ? <IconSunOff /> : <IconSun />}
                                onClick={async () => await api.state.setFlash(!hardwareFlashPowered)}
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
