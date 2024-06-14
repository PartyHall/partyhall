import useAsyncEffect from "use-async-effect";

import { Pagination, Stack, TextField, Typography } from "@mui/material";
import { useEffect, useRef, useState } from "react";
import { KaraokeSong } from "../../../types/appstate";
import { useApi } from "../../../hooks/useApi";
import { PaginatedResponse } from "../../../sdk/responses/paginated_responses";
import Loader from "../../../components/loader";
import Song from "./song";
import { useTranslation } from "react-i18next";
import { debounce } from "lodash";

export default function KaraokeSearch() {
    const {t} = useTranslation();
    const {api} = useApi();

    const [loading, setLoading] = useState<boolean>(true);
    const [page, setPage] = useState<number>(1);
    const [search, setSearch] = useState<string>("");

    const [apiResponse, setApiResponse] = useState<PaginatedResponse<KaraokeSong>|null>(null);

    const doSearch = async (query: string, requestedPage: number) => {
        setLoading(true);
        setApiResponse(await api.karaoke.search(query, requestedPage));
        setLoading(false);
    };

    const debouncedSearch = useRef(
        debounce(async (query: string, requestedPage: number) => await doSearch(query, requestedPage), 250)
    ).current;

    useAsyncEffect(async () => {
        await debouncedSearch('', 1);
    }, []);

    useEffect(() => {
        return () => {
            debouncedSearch.cancel();
        }
    }, [debouncedSearch]);

    return <Stack direction="column" alignItems="stretch" flex={1} gap={1}>
        <TextField placeholder={t('karaoke.search') + '...'} value={search} onChange={x => {
            setSearch(x.target.value);
            setPage(1);
            debouncedSearch(x.target.value, 1);
        }} />
        <Loader loading={loading}>
            <Stack direction="column" gap={1} pt={1} flex="1 1 0" style={{overflowY: 'scroll'}}>
                { 
                    apiResponse &&
                    apiResponse.results.map(x => <Song key={x.id} song={x} type="SEARCH" mb={0} />) 
                }
            </Stack>
        </Loader>
        <Stack direction="column" alignItems="center" justifyContent="center" gap={2}>
            <Typography variant="body1">{t('karaoke.amt_songs')}: {apiResponse?.meta.total}</Typography>
            <Pagination 
                count={apiResponse?.meta.last_page}
                shape="rounded"
                variant="outlined"
                page={page}
                onChange={async (event: React.ChangeEvent<unknown>, value: number) => {
                    setPage(value);
                    await doSearch(search, value);
                }}
            />
        </Stack>
    </Stack>
}