import { Button, Card, Flex, Form, FormProps, Input, Switch, Typography, notification } from 'antd';

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
    const { t } = useTranslation();
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

    /**
     * Not sure it even happens
     */
    const onFinishFailed: FormProps<FormType>['onFinishFailed'] = (errorInfo) => {
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
                    <img src={PhLogo} alt="PartyHall logo" style={{ display: 'block', maxHeight: '4em' }} />
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
                    <Flex vertical style={{ marginTop: '2em', marginBottom: '2em' }}>
                        <Form.Item<FormType>
                            label={t('login.username')}
                            name="username"
                            rules={[
                                {
                                    required: true,
                                    message: t('login.username_required'),
                                },
                            ]}
                        >
                            <Input />
                        </Form.Item>

                        {(!guests_allowed || admin) && (
                            <Form.Item<FormType>
                                label={t('login.password')}
                                name="password"
                                rules={[
                                    {
                                        required: true,
                                        message: t('login.password_required'),
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
                                checkedChildren={t('login.admin')}
                                unCheckedChildren={t('login.anonymous')}
                                onChange={(x) => setAdmin(x.valueOf())}
                                checked={admin}
                            />
                        )}
                        <Button type="primary" htmlType="submit">
                            {t('login.connect')}
                        </Button>
                    </Flex>
                </Form>
            </Flex>
            {notifCtx}
        </Card>
    );
}
