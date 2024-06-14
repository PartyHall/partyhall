import { Button, Card, CardActions, CardContent, CircularProgress, IconButton, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography } from "@mui/material";
import { DateTime } from "luxon";
import { useState } from "react";
import DownloadIcon from '@mui/icons-material/Download'

import useAsyncEffect from "use-async-effect";
import { useApi } from "../../hooks/useApi";
import { useAdminSocket } from "../../hooks/adminSocket";

import { EventExport } from "../../types/event_export";
import { useSnackbar } from "../../hooks/snackbar";
import { useConfirmDialog } from "../../hooks/dialog";
import { useTranslation } from "react-i18next";

export default function ExportListing() {
    const {t} = useTranslation();
    const { showSnackbar } = useSnackbar();
    const { showDialog } = useConfirmDialog();
    const { api } = useApi();
    const { appState, sendMessage, lastMessage } = useAdminSocket();

    const [downloadInProgress, setDownloadInProgress] = useState<boolean>(false);

    const [lastExports, setLastExports] = useState<EventExport[]>([]);

    const fetchLastExports = async () => {
        if (!appState?.app_state?.current_event) {
            return
        }

        setLastExports(await api.events.getLastExports(appState.app_state.current_event.id));
    };

    useAsyncEffect(async () => {
        if (lastMessage?.type === 'EXPORT_COMPLETED') {
            await fetchLastExports();
        }
    }, [lastMessage]);

    useAsyncEffect(async () => {
        await fetchLastExports();
    }, []);

    const exportAsZip = () => {
        sendMessage('EXPORT_ZIP', appState?.app_state.current_event?.id);
    }

    const download = async (id: number) => {
        setDownloadInProgress(true);
        try {
            const resp = await api.events.downloadExport(id);
            if (resp.status != 200) {
                throw t('exports.failed_to_download');
            }

            /** @TODO: send the Content-Disposition maybe 😂 */
            const filename = resp.headers.get('Content-Disposition')?.split('filename=')[1] ?? 'partyhall.zip';
            const data = await resp.blob();
            const anchor = document.createElement('a');

            anchor.download = filename;
            anchor.href = window.URL.createObjectURL(data);
            anchor.click();
        } catch (e) {
            showSnackbar(t('general.export_occured') + ': ' + e, 'error');
        }
        setDownloadInProgress(false);
    };

    const askExport = () => showDialog(
        t('exports.export_as_zip'),
        <div dangerouslySetInnerHTML={{__html: t('exports.modal_infos', {name: appState.app_state.current_event?.name})}} />,
        t('exports.export'),
        async () => {
            exportAsZip();
        },
    );

    return <>
        <Card>
            <CardContent>
                <Typography variant="h2" fontSize={18}>{t('exports.last_exports')}</Typography>
                <TableContainer component={Paper}>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>{t('exports.file')}</TableCell>
                                <TableCell>{t('exports.date')}</TableCell>
                                <TableCell></TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {
                                lastExports.map(k => <TableRow key={k.id}>
                                    <TableCell>{k.filename}</TableCell>
                                    <TableCell>{DateTime.fromSeconds(k.date).toFormat("yyyy-MM-dd HH:mm:ss")}</TableCell>
                                    <TableCell>
                                        <IconButton onClick={() => download(k.id)} disabled={downloadInProgress}>
                                            {!downloadInProgress && <DownloadIcon />}
                                            {downloadInProgress && <CircularProgress />}
                                        </IconButton>
                                    </TableCell>
                                </TableRow>)
                            }
                        </TableBody>
                    </Table>
                </TableContainer>
            </CardContent>
            <CardActions>
                <Button style={{ width: '100%' }} color="error" onClick={askExport}>{t('exports.export_as_zip')}</Button>
            </CardActions>
        </Card>
    </>
}