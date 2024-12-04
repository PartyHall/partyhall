export class Collection<T> {
    results: T[];
    currentPage: number;
    perPageCount: number;
    pageCount: number;
    totalCount: number;

    constructor(
        results: T[],
        currentPage: number,
        perPageCount: number,
        pageCount: number,
        totalCount: number
    ) {
        this.results = results;
        this.currentPage = currentPage;
        this.perPageCount = perPageCount;
        this.pageCount = pageCount;
        this.totalCount = totalCount;
    }

    static fromJson<U>(
        json: Record<string, any> | null,
        objCreator: (data: Record<string, any> | null) => U
    ) {
        if (!json) {
            return null;
        }

        return new Collection<U>(
            (json['results'] || []).map((x: Record<string, any>) =>
                objCreator(x)
            ),
            json['current_page'],
            json['per_page_count'],
            json['page_count'],
            json['total_count']
        );
    }
}
