import {
    Button,
    Card,
    Flex,
    Form,
    FormProps,
    Input,
    Switch,
    Typography,
    notification,
} from 'antd';

import { useEffect, useState } from 'react';

import PhLogo from '../assets/ph_logo_sd.webp';

import { useAuth } from '../hooks/auth';
import { useNavigate } from 'react-router-dom';
import { useSettings } from '../hooks/settings';
import { useTranslation } from 'react-i18next';

type FormType = {
    username?: string;
    password?: string;
};

export default function LoginPage() {
    const { t } = useTranslation('', { keyPrefix: 'login' });
    const { guests_allowed } = useSettings();
    const [admin, setAdmin] = useState<boolean>(!guests_allowed);
    const [notifApi, notifCtx] = notification.useNotification();
    const { api, isLoggedIn, login, loginGuest } = useAuth();
    const navigate = useNavigate();

    const onFinish: FormProps<FormType>['onFinish'] = async (values) => {
        if (!values.username || (admin && !values.password)) {
            return;
        }

        try {
            if (admin) {
                await login(values.username, values.password!);
            } else {
                await loginGuest(values.username);
            }
        } catch (e: any) {
            console.error(e);

            const isBadLogin =
                e.message?.type ===
                'https://github.com/partyhall/partyhall/bad-login';

            notifApi.error({
                message: t('notification_failed_login.title'),
                description: t(
                    'notification_failed_login.' +
                        (e.status === 400 && isBadLogin
                            ? 'bad_login'
                            : 'description')
                ),
            });
        }
    };

    const onFinishFailed: FormProps<FormType>['onFinishFailed'] = (
        errorInfo
    ) => {
        notifApi.error({
            message: 'Failed to login',
            description: JSON.stringify(errorInfo), // @TODO make it clean
        });
    };

    useEffect(() => {
        if (isLoggedIn()) {
            navigate('/');
        }
    }, [api]);

    return (
        <Card>
            <Flex gap="middle" align="center" justify="stretch" vertical>
                <Typography>
                    <img
                        src={PhLogo}
                        alt="PartyHall logo"
                        style={{ display: 'block', maxHeight: '4em' }}
                    />
                </Typography>
                <Form
                    name="login"
                    labelCol={{ span: 8 }}
                    wrapperCol={{ span: 16 }}
                    style={{ maxWidth: 600 }}
                    onFinish={onFinish}
                    onFinishFailed={onFinishFailed}
                    autoComplete="off"
                    className="hide-asterisk"
                >
                    <Flex
                        vertical
                        style={{ marginTop: '2em', marginBottom: '2em' }}
                    >
                        <Form.Item<FormType>
                            label={t('username')}
                            name="username"
                            rules={[
                                {
                                    required: true,
                                    message: t('username_required'),
                                },
                            ]}
                        >
                            <Input />
                        </Form.Item>

                        {(!guests_allowed || admin) && (
                            <Form.Item<FormType>
                                label={t('password')}
                                name="password"
                                rules={[
                                    {
                                        required: true,
                                        message: t('password_required'),
                                    },
                                ]}
                            >
                                <Input.Password />
                            </Form.Item>
                        )}
                    </Flex>

                    <Flex
                        gap="middle"
                        style={{
                            alignItems: 'center',
                            justifyContent: 'space-around',
                        }}
                    >
                        {guests_allowed && (
                            <Switch
                                checkedChildren={t('admin')}
                                unCheckedChildren={t('anonymous')}
                                onChange={(x) => setAdmin(x.valueOf())}
                                checked={admin}
                            />
                        )}
                        <Button type="primary" htmlType="submit">
                            {t('connect')}
                        </Button>
                    </Flex>
                </Form>
            </Flex>
            {notifCtx}
        </Card>
    );
}
