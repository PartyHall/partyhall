import { Button, Card, CardActions, CardContent, CircularProgress, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, IconButton, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography } from "@mui/material";
import { DateTime } from "luxon";
import { useState } from "react";
import DownloadIcon from '@mui/icons-material/Download'
import useAsyncEffect from "use-async-effect";
import { useApi } from "../hooks/auth";
import { useWebsocket } from "../hooks/ws";
import { EventExport } from "../types/event_export";

export default function ExportListing() {
    const { password, showError, getLastExports } = useApi();
    const { appState, sendMessage, lastMessage } = useWebsocket();
    const [exportZipShown, setExportZipShown] = useState<boolean>(false);

    const [downloadInProgress, setDownloadInProgress] = useState<boolean>(false);

    const [lastExports, setLastExports] = useState<EventExport[]>([]);

    const fetchLastExports = async () => {
        if (!appState?.app_state?.current_event) {
            return
        }

        const exports = await getLastExports(appState.app_state.current_event.id);
        setLastExports(exports);
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
        /**
         * Something goes terribly wrong with this
         */
        try {
            const resp = await fetch(
                `/api/admin/exports/${id}/download`,
                { 'headers': { 'Authorization': password ?? '' } }
            );
            if (resp.status != 200) {
                throw 'Failed to download file';
            }

            const filename = resp.headers.get('Content-Disposition')?.split('filename=')[1] ?? 'partyhall.zip';
            const data = await resp.blob();
            const anchor = document.createElement('a');

            anchor.download = filename;
            anchor.href = window.URL.createObjectURL(data);
            anchor.click();
        } catch (e) {
            showError('An error has occured: ' + e, 'error');
        }
        setDownloadInProgress(false);
    };

    return <>
        <Card>
            <CardContent>
                <Typography variant="h2" fontSize={18}>Last exports</Typography>
                <TableContainer component={Paper}>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>File</TableCell>
                                <TableCell>Date</TableCell>
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
                <Button style={{ width: '100%' }} color="error" onClick={() => setExportZipShown(true)}>Export as zip</Button>
            </CardActions>
        </Card>

        <Dialog open={exportZipShown} onClose={() => setExportZipShown(false)}>
            <DialogTitle>Export as zip</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    You are trying to export the event {appState?.app_state?.current_event?.name}. <br />
                    This will create a zip with all the pictures and let you download, thus it could take a long time. <br />
                    Are you sure you want to continue ?
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => setExportZipShown(false)}>Cancel</Button>
                <Button onClick={() => {
                    exportAsZip();
                    setExportZipShown(false);
                }} color="warning" autoFocus>Export zip</Button>
            </DialogActions>
        </Dialog>

    </>
}