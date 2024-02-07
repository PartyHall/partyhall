import { ReactNode } from "react";
import { Navigate, useOutlet } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import useKeyPress from "../hooks/useKeyPress";
import { useBoothSocket } from "../hooks/boothSocket";
import { useTranslation } from "react-i18next";

export default function BoothLayout() {
    const {t} = useTranslation();
    const outlet = useOutlet();
    const {socketMode} = useApi();
    const { appState, currentTime, showDebug } = useBoothSocket();

    useKeyPress(['d'], (event: any) => {
        if (event.key === 'd') {
            showDebug();
        }
    });

    if (socketMode != 'booth') {
        return <Navigate to={"/admin/login"} />
    }

    if (!appState) {
        return <div className="debug abstl">{t('general.something_went_wrong')}</div>
    }

    const datetime = currentTime ?? t('general.dt_not_available');
    const eventName = !!appState.app_state?.current_event ? appState.app_state.current_event.name : t('osd.no_event');

    const D = (title: string, child: ReactNode) => <div><span style={{ fontWeight: 'bold' }}>{title}</span>: {child}</div>

    return <>
        <div className="debug abstl">
            {<span>{eventName}</span>}
            {
                appState.debug && <>
                    {D(t('admin_main.mode'), <span>{appState.current_mode}</span>)}
                    {D(t('admin_main.hw_flash'), <span>{appState.modules.photobooth.hardware_flash ? t('general.yes') : t('general.no')}</span>)}
                    {D('IPs', <ul>
                        {
                            appState.ip_addresses && Object.entries(appState.ip_addresses).filter(([_, x]) => x.length > 0).map(([key, inter]) => <li key={key}>
                                {key}: {inter.join(', ')}
                            </li>)
                        }
                    </ul>)}
                </>
            }
        </div>
        <div className="debug abstr">
            <span>{datetime}</span>
            {
                appState.debug && <>
                    {D('HWID', <span>{appState.app_state.hwid}</span>)}
                    {D('Token', <span>{appState.app_state.token}</span>)}
                    {D(t('general.version'), <span>{appState.partyhall_version}</span>)}
                    {D('Commit', <span>{appState.partyhall_commit}</span>)}
                </>
            }
        </div>

        {outlet}
    </>
}