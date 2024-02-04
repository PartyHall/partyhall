import useAsyncEffect from "use-async-effect";

import { Pagination, Stack, TextField, Typography } from "@mui/material";
import { useState } from "react";
import { KaraokeSong } from "../../../types/appstate";
import { useApi } from "../../../hooks/useApi";
import { PaginatedResponse } from "../../../sdk/responses/paginated_responses";
import Loader from "../../../components/loader";
import Song from "./song";

// @TODO: debounce
export default function KaraokeSearch() {
    const [wasSearch, setWasSearch] = useState<boolean>(false);
    const {api} = useApi();

    const [loading, setLoading] = useState<boolean>(true);
    const [page, setPage] = useState<number>(1);
    const [search, setSearch] = useState<string>("");

    const [apiResponse, setApiResponse] = useState<PaginatedResponse<KaraokeSong>|null>(null);

    const loadSongs = async () => {
        let currPage = page;
        if ((wasSearch && search.length == 0) || (!wasSearch && search.length > 0)) {
            currPage = 1;
            setWasSearch(!wasSearch);
            setPage(currPage);
        }

        setLoading(true);
        setApiResponse(await api.karaoke.search(search, currPage));
        setLoading(false);
    };

    useAsyncEffect(async () => {
        await loadSongs();
    }, [search, page]);

    return <Stack direction="column" alignItems="stretch" flex={1} gap={1}>
        <TextField placeholder="Search..." value={search} onChange={x => setSearch(x.target.value)}/>
        <Loader loading={loading}>
            <Stack direction="column" gap={1} pt={1} flex="1 1 0" style={{overflowY: 'scroll'}}>
                { 
                    apiResponse &&
                    apiResponse.results.map(x => <Song key={x.id} song={x} type="SEARCH" mb={0} />) 
                }
            </Stack>
        </Loader>
        <Stack direction="column" alignItems="center" justifyContent="center" gap={2}>
            <Typography variant="body1">Amt of songs: {apiResponse?.meta.total}</Typography>
            <Pagination 
                count={apiResponse?.meta.last_page}
                shape="rounded"
                variant="outlined"
                page={page}
                onChange={(event: React.ChangeEvent<unknown>, value: number) => setPage(value)}
            />
        </Stack>
    </Stack>
}