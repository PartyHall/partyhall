import { Box, FormControl, InputLabel, MenuItem, Select, Stack } from "@mui/material";
import { useState } from "react";

import ImportPhk from "./import_phk";
import CreateSong from "./create_song";

type METHOD = 'CREATE'|'IMPORT';

export default function KaraokeAddSong() {
    const [method, setMethod] = useState<METHOD>('CREATE');
 
    return <Stack direction="column" gap={2} flex="1">
        <Box mb={3}>
            <FormControl fullWidth>
                <InputLabel id="label_select_method">What to do?</InputLabel>
                <Select labelId="label_select_method" value={method} onChange={x => setMethod(x.target.value as METHOD)} style={{marginBottom: 3}}>
                    <MenuItem value='CREATE'>Create</MenuItem>
                    <MenuItem value='IMPORT'>Import .phk</MenuItem>
                </Select>
            </FormControl>
        </Box>

        {method === 'IMPORT' && <ImportPhk />}
        {method === 'CREATE' && <CreateSong />}
    </Stack>;
}