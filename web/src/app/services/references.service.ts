import { HttpClient } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import { injectMutation, injectQuery, injectQueryClient } from "@ngneat/query";
import { DefaultHttpProxyService } from "./http/default-http-proxy.service";
import { Observable } from "rxjs";
import { client } from "./http/kiota";

@Injectable({
	providedIn: "root",
})
export class ReferencesService {
	#http = inject(DefaultHttpProxyService);
	#query = injectQuery();
	#mutate = injectMutation();
	#queryClient = injectQueryClient();

	upload = this.#mutate({
		mutationFn: (params: {
			files: File[];
			minSURFMatches: number;
			minThresholdForSURFMatches: number;
			mseSkip: number;
			ratioTestThreshold: number;
		}) => {
			return client.api.referencesCreate({
				file: params.files[0],
				minSURFMatches: params.minSURFMatches,
				minThresholdForSURFMatches: params.minThresholdForSURFMatches,
				mseSkip: params.mseSkip,
				ratioTestThreshold: params.ratioTestThreshold,
			});
		},
		onSuccess: () => {
			this.#queryClient.invalidateQueries({
				queryKey: ["references"],
			});
			this.#queryClient.invalidateQueries({
				queryKey: ["filters"],
			});
		},
	});

	getReferences() {
		return this.#query({
			queryKey: ["references"] as const,
			queryFn: () => {
				return this.#http.get("/references") as Observable<{ id: string }[]>;
			},
		});
	}

	getReferenceById(id: string) {
		return this.#query({
			queryKey: ["references", id],
			queryFn: () => client.api.getApiReferencesById(id).then((x) => x.data),
		});
	}

	constructor() {}
}
