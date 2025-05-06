import { Button, Card, Flex, Form, Input, Typography, notification } from 'antd';
import { Controller, useForm, useWatch } from 'react-hook-form';
import PhLogo from '../assets/ph_logo_sd.webp';

import { useAuth } from '../hooks/auth';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSettings } from '../hooks/settings';
import { useTranslation } from 'react-i18next';

type FormType = {
    username?: string;
    display_name?: string;
    password?: string;
    password2?: string;
};

export default function CreateAdminPage() {
    const { adminCreated, setAdminCreated } = useSettings();
    const { t } = useTranslation();
    const [notifApi, notifCtx] = notification.useNotification();
    const { api, isLoggedIn } = useAuth();
    const navigate = useNavigate();

    const {
        control,
        handleSubmit,
        formState: { errors, isSubmitting, isValid },
        trigger,
        clearErrors,
    } = useForm<FormType>({ mode: 'onChange', reValidateMode: 'onChange' });

    const password = useWatch({ control, name: 'password' });
    const password2 = useWatch({ control, name: 'password2' });

    useEffect(() => {
        if (!password || !password2) {
            clearErrors('password2');
            return;
        }

        trigger('password2');
    }, [password, password2, trigger, clearErrors]);

    const onSubmit = async (values: FormType) => {
        if (!values.username || !values.password || !values.password2) {
            return;
        }

        try {
            const resp = await api.admin.createAdmin(values.display_name || null, values.username, values.password);

            if (!resp) {
                throw new Error('generic.error.unexpected');
            }

            notifApi.success({
                message: t('create_admin.success.title', { username: resp.username }),
                description: t('create_admin.success.description'),
            });

            setAdminCreated(true);
        } catch (e: any) {
            console.error(e);

            let description = 'generic.error.unexpected';
            if (e.message?.type === 'bad-login') {
                description = 'login.bad_login';
            }

            notifApi.error({
                message: t('generic.error.title'),
                description: t(description),
            });
        }
    };

    const onError = (errorInfo: any) => {
        notifApi.error({
            message: t('generic.error.title'),
            description: JSON.stringify(errorInfo),
        });
    };

    useEffect(() => {
        // If the user is already logged in
        // no need to check for admin creation
        if (isLoggedIn()) {
            navigate('/');

            return;
        }

        if (adminCreated) {
            navigate('/login');

            return;
        }
    }, [api, adminCreated]);

    return (
        <Card>
            <Flex gap={4} align="center" justify="stretch" vertical>
                <img src={PhLogo} alt="PartyHall logo" style={{ display: 'block', maxHeight: '4em' }} />

                <Typography.Title level={2} style={{ margin: 0 }}>
                    {t('create_admin.title')}
                </Typography.Title>

                <Typography.Paragraph>{t('create_admin.p1')}</Typography.Paragraph>

                <Form
                    onFinish={handleSubmit(onSubmit)}
                    onFinishFailed={onError}
                    layout="horizontal"
                    labelCol={{ span: 8 }}
                    wrapperCol={{ span: 16 }}
                    className="hide-asterisk"
                    style={{ maxWidth: 800 }}
                    autoComplete="off"
                >
                    <Flex vertical>
                        <Form.Item
                            label={t('login.username')}
                            validateStatus={errors.username ? 'error' : ''}
                            help={errors.username?.message}
                        >
                            <Controller
                                name="username"
                                control={control}
                                rules={{ required: t('login.username_required') }}
                                render={({ field }) => <Input {...field} disabled={isSubmitting} />}
                            />
                        </Form.Item>

                        <Form.Item label={t('create_admin.display_name')}>
                            <Controller
                                name="display_name"
                                control={control}
                                render={({ field }) => <Input {...field} disabled={isSubmitting} />}
                            />
                        </Form.Item>

                        <Form.Item
                            label={t('login.password')}
                            validateStatus={errors.password ? 'error' : ''}
                            help={errors.password?.message}
                        >
                            <Controller
                                name="password"
                                control={control}
                                rules={{ required: t('login.password_required') }}
                                render={({ field }) => <Input.Password {...field} disabled={isSubmitting} />}
                            />
                        </Form.Item>

                        <Form.Item
                            label={t('create_admin.password_repeat')}
                            validateStatus={errors.password2 ? 'error' : ''}
                            help={errors.password2?.message}
                        >
                            <Controller
                                name="password2"
                                control={control}
                                rules={{
                                    required: t('login.password_required'),
                                    validate: (value) => value === password || t('create_admin.passwords_mismatch'),
                                }}
                                render={({ field }) => <Input.Password {...field} disabled={isSubmitting} />}
                            />
                        </Form.Item>
                    </Flex>

                    <Flex
                        gap="middle"
                        style={{
                            alignItems: 'center',
                            justifyContent: 'space-around',
                        }}
                    >
                        <Button
                            type="primary"
                            htmlType="submit"
                            loading={isSubmitting}
                            disabled={!isValid || isSubmitting}
                        >
                            {t('create_admin.register')}
                        </Button>
                    </Flex>
                </Form>
            </Flex>
            {notifCtx}
        </Card>
    );
}
