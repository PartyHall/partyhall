import { Box, FormControl, InputLabel, MenuItem, Select, Stack } from "@mui/material";
import { useState } from "react";

import ImportPhk from "./import_phk";
import CreateSong from "./create_song";
import SettingsKaraoke from "./settings";

type METHOD = 'CREATE'|'IMPORT'|'SETTINGS';

export default function KaraokeSettings() {
    const [method, setMethod] = useState<METHOD>('CREATE');
 
    return <Stack direction="column" gap={2} flex="1">
        <Box mb={3}>
            <FormControl fullWidth>
                <InputLabel id="label_select_method">What to do?</InputLabel>
                <Select labelId="label_select_method" value={method} onChange={x => setMethod(x.target.value as METHOD)} style={{marginBottom: 3}}>
                    <MenuItem value='CREATE'>Add a song</MenuItem>
                    <MenuItem value='IMPORT'>Import .phk</MenuItem>
                    <MenuItem value='SETTINGS'>Settings</MenuItem>
                </Select>
            </FormControl>
        </Box>

        {method === 'IMPORT' && <ImportPhk />}
        {method === 'CREATE' && <CreateSong />}
        {method === 'SETTINGS' && <SettingsKaraoke />}
    </Stack>;
}