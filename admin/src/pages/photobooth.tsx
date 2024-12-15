import { Button, Flex, Typography, notification } from 'antd';
import KeyVal from '../components/keyval';
import Title from 'antd/es/typography/Title';
import { useAuth } from '../hooks/auth';
import { useEffect } from 'react';
import { useSettings } from '../hooks/settings';
import { useTranslation } from 'react-i18next';

export default function Photobooth() {
    const { t } = useTranslation('', { keyPrefix: 'photobooth' });

    const [notif, ctxHolder] = notification.useNotification();
    const { setPageName } = useSettings();
    const { event, api } = useAuth();

    useEffect(() => setPageName('photobooth'), []);

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
