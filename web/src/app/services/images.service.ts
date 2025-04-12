import { Injectable } from "@angular/core";
import {
	InfiniteData,
	InfiniteQueryObserverResult,
	injectInfiniteQuery,
	injectQuery,
	QueryKey,
	QueryObserverResult,
} from "@ngneat/query";
import { client } from "./http/kiota";
import { undefToErr } from "./http/undefToErr";

@Injectable({
	providedIn: "root",
})
export class ImagesService {
	#query = injectQuery();

	getImagePage(pageParam: number, pageSize: number, youtube_id: string) {
		return this.#query({
			queryKey: ["images", youtube_id, pageParam, pageSize] as const,
			enabled: false,
			refetchInterval: 5000,
			queryFn: () => this.getImagePageApi(pageParam, pageSize, youtube_id),
		});
	}

	getImagePageApi(pageNumber: number, pageSize: number, youtube_id: string) {
		const offset = pageNumber * pageSize;
		const limit = pageSize;

		return undefToErr(
			client.api
				.getApiImages({
					limit,
					offset,
					youtube_id,
				})
				.then((x) => x.data),
		);
	}
}
