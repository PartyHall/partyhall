import { Box, Stack, Tab, Tabs } from "@mui/material";
import { useState } from "react";
import KaraokeSearch from "./search";
import KaraokeQueue from "./queue";

export default function AdminKaraoke() {
    const [currentTab, setCurrentTab] = useState<number>(0);

    return <Stack direction="column" flex={1} style={{ marginTop: 0, height: "100%"}} gap={2}>
        <Box>
            <Tabs value={currentTab} onChange={(_, x) => setCurrentTab(x)}>
                <Tab label="Search"/>
                <Tab label="Queue"/>
            </Tabs>
        </Box>

        <Stack flex={1} overflow={{y: "scroll"}}>
            {
                currentTab == 0 &&
                <KaraokeSearch />
            }
            {
                currentTab == 1 &&
                <KaraokeQueue />
            }
        </Stack>
    </Stack>
}