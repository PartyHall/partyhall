import { IconButton, Stack } from "@mui/material";
import { KaraokeSongSession } from "../../../../types/appstate";
import { useAdminSocket } from "../../../../hooks/adminSocket";

import { ArrowUpward as UpIcon, ArrowDownward as DownIcon } from '@mui/icons-material';

export default function QueueMoveButtons({isFirst, isLast, session}: {isFirst: boolean, isLast: boolean, session: KaraokeSongSession}) {
    const { sendMessage } = useAdminSocket();

    return <Stack direction="column" alignItems="center" justifyContent="center" gap={1}>
        <IconButton style={{ padding: ".25em" }} disabled={isFirst} onClick={() => {
            sendMessage('karaoke/QUEUE_MOVE_UP', session.id)
        }}>
            <UpIcon />
        </IconButton>
        <IconButton style={{ padding: ".25em" }} disabled={isLast} onClick={() => {
            sendMessage('karaoke/QUEUE_MOVE_DOWN', session.id)
        }}>
            <DownIcon />
        </IconButton>
    </Stack>;
}