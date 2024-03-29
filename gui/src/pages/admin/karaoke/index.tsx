import { Box, Stack, Tab, Tabs } from "@mui/material";
import { useState } from "react";
import KaraokeSearch from "./search";
import KaraokeQueue from "./queue";
import KaraokeSettings from "./new_song/index";
import { useApi } from "../../../hooks/useApi";
import { useTranslation } from "react-i18next";

export default function AdminKaraoke() {
    const {hasRole} = useApi();
    const {t} = useTranslation();
    const [currentTab, setCurrentTab] = useState<number>(0);

    return <Stack direction="column" flex={1} style={{ marginTop: 0, height: "100%"}} gap={3}>
        <Box>
            <Tabs value={currentTab} onChange={(_, x) => setCurrentTab(x)}>
                <Tab label={t('karaoke.search')}/>
                <Tab label={t('karaoke.queue')}/>
                { hasRole('ADMIN_KARAOKE') && <Tab label={t('karaoke.admin')}/> }
            </Tabs>
        </Box>

        <Stack flex="1 0 0" overflow="scroll">
            {
                currentTab == 0 &&
                <KaraokeSearch />
            }
            {
                currentTab == 1 &&
                <KaraokeQueue />
            }
            {
                currentTab == 2 &&
                <KaraokeSettings />
            }
        </Stack>
    </Stack>
}