import { useTranslation } from "react-i18next";

export default function LockedModal() {
    const {t} = useTranslation();

    return <div id="modal">
        <h1>{t('disabled.partyhall_disabled')}</h1>
        <span>{t('disabled.msg')}</span>
    </div>;
}