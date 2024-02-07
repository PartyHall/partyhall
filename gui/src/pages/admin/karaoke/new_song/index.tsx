import { Box, FormControl, InputLabel, MenuItem, Select, Stack } from "@mui/material";
import { useState } from "react";

import ImportPhk from "./import_phk";
import CreateSong from "./create_song";
import SettingsKaraoke from "./settings";
import { useTranslation } from "react-i18next";

type METHOD = 'CREATE'|'IMPORT'|'SETTINGS';

export default function KaraokeSettings() {
    const [method, setMethod] = useState<METHOD>('CREATE');
    const {t} = useTranslation();
 
    return <Stack direction="column" gap={2} flex="1" pt={2} pb={2}>
        <Box mb={3}>
            <FormControl fullWidth>
                <InputLabel id="label_select_method">{t('karaoke.what_to_do')}</InputLabel>
                <Select labelId="label_select_method" value={method} onChange={x => setMethod(x.target.value as METHOD)} style={{marginBottom: 3}}>
                    <MenuItem value='CREATE'>{t('karaoke.wtd_add_song')}</MenuItem>
                    <MenuItem value='IMPORT'>{t('karaoke.wtd_import')}</MenuItem>
                    <MenuItem value='SETTINGS'>{t('karaoke.wtd_settings')}</MenuItem>
                </Select>
            </FormControl>
        </Box>

        {method === 'IMPORT' && <ImportPhk />}
        {method === 'CREATE' && <CreateSong />}
        {method === 'SETTINGS' && <SettingsKaraoke />}
    </Stack>;
}