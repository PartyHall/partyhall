import { Content, Footer, Header } from 'antd/es/layout/layout';
import { Flex, Layout, Menu, Typography } from 'antd';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useEffect, useState } from 'react';

import PhLogo from '../../assets/ph_logo_sd.webp';

import Sider from 'antd/es/layout/Sider';

import { useAuth } from '../../hooks/auth';
import { useSettings } from '../../hooks/settings';
import { useTranslation } from 'react-i18next';

export default function AuthedLayout() {
    const { t } = useTranslation('', { keyPrefix: 'generic.menu' });
    const { pageName } = useSettings();
    const { kioskMode, isLoggedIn, api } = useAuth();
    const [collapsed, setCollapsed] = useState(true);
    const navigate = useNavigate();

    const location = useLocation();

    let pages = [
        { key: 'home', label: t('home'), target: '/' },
        { key: 'photobooth', label: t('photobooth'), target: '/photobooth' },
        { key: 'karaoke', label: t('karaoke'), target: '/karaoke' },
    ];

    if (api.tokenUser?.roles.includes('ADMIN')) {
        pages = [...pages.slice(0, 1), { key: 'events', label: t('events'), target: '/events' }, ...pages.slice(1)];

        pages.push({ key: 'settings', label: t('settings'), target: '/settings' });
        pages.push({ key: 'logs', label: t('logs'), target: '/logs' });
    }

    const handleMenuClick = ({ key }: { key: string }) => {
        const { target } = pages.find((item) => item.key === key) || {};
        if (target) {
            // Crappy hack but it works
            const width = window.visualViewport?.width;
            if (width && width < 991) {
                setCollapsed(true);
            }

            navigate(target);
        }
    };

    useEffect(() => {
        if (!isLoggedIn()) {
            navigate('/login');
            return;
        }
    }, [api.token]);

    useEffect(() => {
        if (kioskMode && location.pathname !== '/kiosk') {
            navigate('/kiosk');
        }
    }, [location.pathname, kioskMode, navigate]);

    return (
        <Layout>
            {!kioskMode && (
                <>
                    <Sider
                        style={{
                            display: 'flex',
                            alignItems: 'center',
                            paddingTop: '2em',
                        }}
                        className="custom-sider"
                        breakpoint="lg"
                        collapsedWidth="0"
                        collapsed={collapsed}
                        onCollapse={(val) => setCollapsed(val)}
                    >
                        <Menu
                            theme="dark"
                            mode="inline"
                            defaultSelectedKeys={['home']}
                            selectedKeys={['' + pageName]}
                            items={pages}
                            style={{ flex: 1, minWidth: 0 }}
                            onClick={handleMenuClick}
                        />
                    </Sider>

                    {!collapsed && <div className="menu-backdrop" onClick={() => setCollapsed(true)} />}
                </>
            )}

            <Layout>
                {!kioskMode && (
                    <Header
                        style={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'space-between',
                        }}
                    >
                        {/* @TODO: Make it a real link for accessibility */}
                        <img
                            src={PhLogo}
                            style={{
                                display: 'block',
                                height: '80%',
                                cursor: 'pointer',
                            }}
                            onClick={() => navigate('/')}
                        />
                    </Header>
                )}
                <Content>
                    <Flex
                        vertical
                        style={{
                            width: '100%',
                            height: '100%',
                            overflow: 'auto',
                        }}
                        align="center"
                    >
                        <Outlet />
                    </Flex>
                </Content>
                {!kioskMode && (
                    <Footer style={{ textAlign: 'center' }}>
                        <Typography>
                            PartyHall -{' '}
                            <a href="https://github.com/partyhall/PartyHall" target="_blank" rel="noopener noreferrer">
                                Github
                            </a>
                        </Typography>
                    </Footer>
                )}
            </Layout>
        </Layout>
    );
}
