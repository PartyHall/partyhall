export class PaginatedResponse<T> {
    results: T[];
    meta: {
        last_page: number;
        total: number;
    };

    constructor(resp: any) {
        this.results = resp.results;
        this.meta = resp.meta;
    }
}