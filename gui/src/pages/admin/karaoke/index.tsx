import { Box, Stack, Tab, Tabs } from "@mui/material";
import { useState } from "react";
import KaraokeSearch from "./search";
import KaraokeQueue from "./queue";
import KaraokeSettings from "./new_song/index";
import { useApi } from "../../../hooks/useApi";

export default function AdminKaraoke() {
    const {hasRole} = useApi();
    const [currentTab, setCurrentTab] = useState<number>(0);

    return <Stack direction="column" flex={1} style={{ marginTop: 0, height: "100%"}} gap={2}>
        <Box>
            <Tabs value={currentTab} onChange={(_, x) => setCurrentTab(x)}>
                <Tab label="Search"/>
                <Tab label="Queue"/>
                { hasRole('ADMIN_KARAOKE') && <Tab label="Admin"/> }
            </Tabs>
        </Box>

        <Stack flex={1}>
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