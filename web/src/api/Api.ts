/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface CreateJob {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/CreateJob.json"
   */
  $schema?: string;
  filter_id?: string;
  is_query_based: boolean;
  is_single_video: boolean;
  /** @format int64 */
  limit?: number;
  search_query?: string;
  youtube_id?: string;
}

export interface CreatedImportJob {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/CreatedImportJob.json"
   */
  $schema?: string;
  id: string;
}

export interface CreatedJob {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/CreatedJob.json"
   */
  $schema?: string;
  id: string;
}

export interface ErrorDetail {
  /** Where the error occurred, e.g. 'body.items[3].tags' or 'path.thing-id' */
  location?: string;
  /** Error message text */
  message?: string;
  /** The value at the given location */
  value?: any;
}

export interface ErrorModel {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/ErrorModel.json"
   */
  $schema?: string;
  /**
   * A human-readable explanation specific to this occurrence of the problem.
   * @example "Property foo is required but is missing."
   */
  detail?: string;
  /** Optional list of individual error details */
  errors?: ErrorDetail[] | null;
  /**
   * A URI reference that identifies the specific occurrence of the problem.
   * @format uri
   * @example "https://example.com/error-log/abc123"
   */
  instance?: string;
  /**
   * HTTP status code
   * @format int64
   * @example 400
   */
  status?: number;
  /**
   * A short, human-readable summary of the problem type. This value should not change between occurrences of the error.
   * @example "Bad Request"
   */
  title?: string;
  /**
   * A URI reference to human-readable documentation for the error.
   * @format uri
   * @default "about:blank"
   * @example "https://example.com/errors/example"
   */
  type?: string;
}

export interface FrameMatchingTestBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/FrameMatchingTestBody.json"
   */
  $schema?: string;
  matched: boolean;
}

export interface JobAndVideos {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/JobAndVideos.json"
   */
  $schema?: string;
  id: string;
  videos: JobVideo[] | null;
}

export interface JobDetails {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/JobDetails.json"
   */
  $schema?: string;
  id: string;
  search_query: string;
  /** @format int64 */
  video_target: number;
  /** @format int64 */
  videos_found: number;
}

export interface JobVideo {
  download_status: string;
  import_status: string;
  youtube_id: string;
}

export interface ListJobBody {
  id: string;
  /** @format int64 */
  limit: number;
  search_query: string;
}

export interface ListVideoBody {
  name: string;
  /** @format int64 */
  progress: number;
  video_id: string;
}

export interface Picture {
  blob_id: string;
  id: string;
  youtube_id: string;
}

export interface ReferenceGetFeatureResponse {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/ReferenceGetFeatureResponse.json"
   */
  $schema?: string;
  BlobIds: string[] | null;
  Discriminator: string;
  ID: string;
  /** @format int64 */
  Minsurfmatches: number;
  /** @format double */
  Minthresholdforsurfmatches: number;
  /** @format double */
  Mseskip: number;
  Name: string;
  /** @format double */
  Ratiotestthreshold: number;
}

export interface ReferenceUploadResponseBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/ReferenceUploadResponseBody.json"
   */
  $schema?: string;
  status: string;
}

export interface Response {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/Response.json"
   */
  $schema?: string;
  pictures: Picture[] | null;
  /** @format int64 */
  total: number;
}

export interface UpdateJobLimitRequestBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/UpdateJobLimitRequestBody.json"
   */
  $schema?: string;
  /** @format int64 */
  limit: number;
}

export interface VideoMetadataBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   * @example "https://example.com/schemas/VideoMetadataBody.json"
   */
  $schema?: string;
  /** @format int64 */
  maxframe: number;
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown>
  extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) =>
    fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter(
      (key) => "undefined" !== typeof query[key],
    );
    return keys
      .map((key) =>
        Array.isArray(query[key])
          ? this.addArrayQueryParam(query, key)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.Text]: (input: any) =>
      input !== null && typeof input !== "string"
        ? JSON.stringify(input)
        : input,
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
              ? JSON.stringify(property)
              : `${property}`,
        );
        return formData;
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(
    params1: RequestParams,
    params2?: RequestParams,
  ): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (
    cancelToken: CancelToken,
  ): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(
      `${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`,
      {
        ...requestParams,
        headers: {
          ...(requestParams.headers || {}),
          ...(type && type !== ContentType.FormData
            ? { "Content-Type": type }
            : {}),
        },
        signal:
          (cancelToken
            ? this.createAbortSignal(cancelToken)
            : requestParams.signal) || null,
        body:
          typeof body === "undefined" || body === null
            ? null
            : payloadFormatter(body),
      },
    ).then(async (response) => {
      const r = response.clone() as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title My API
 * @version 1.0.0
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  api = {
    /**
     * No description
     *
     * @tags Images
     * @name ImagesList
     * @summary List images
     * @request GET:/api/images
     */
    imagesList: (
      query?: {
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
        youtube_id?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<Response, ErrorModel>({
        path: `/api/images`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Jobs
     * @name JobsList
     * @summary List all jobs
     * @request GET:/api/jobs
     */
    jobsList: (params: RequestParams = {}) =>
      this.request<ListJobBody[] | null, ErrorModel>({
        path: `/api/jobs`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Jobs
     * @name JobsCreate
     * @summary Create a new job
     * @request POST:/api/jobs
     */
    jobsCreate: (data: CreateJob, params: RequestParams = {}) =>
      this.request<CreatedJob, ErrorModel>({
        path: `/api/jobs`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Jobs
     * @name JobsVideoCreate
     * @summary Create a direct video job
     * @request POST:/api/jobs/video
     */
    jobsVideoCreate: (
      data: {
        /** @format binary */
        file: File;
        filter_id: string;
        name?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<CreatedImportJob, ErrorModel>({
        path: `/api/jobs/video`,
        method: "POST",
        body: data,
        type: ContentType.FormData,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Jobs
     * @name JobsDetail
     * @summary Get job details
     * @request GET:/api/jobs/{id}
     */
    jobsDetail: (id: string, params: RequestParams = {}) =>
      this.request<JobDetails, ErrorModel>({
        path: `/api/jobs/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Jobs
     * @name JobsActionsRestartCreate
     * @summary Restart the job pipeline
     * @request POST:/api/jobs/{id}/actions/restart
     */
    jobsActionsRestartCreate: (id: string, params: RequestParams = {}) =>
      this.request<void, ErrorModel>({
        path: `/api/jobs/${id}/actions/restart`,
        method: "POST",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Jobs
     * @name JobsActionsUpdateLimitCreate
     * @summary Update job limit
     * @request POST:/api/jobs/{id}/actions/update-limit
     */
    jobsActionsUpdateLimitCreate: (
      id: string,
      data: UpdateJobLimitRequestBody,
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorModel>({
        path: `/api/jobs/${id}/actions/update-limit`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * No description
     *
     * @tags Jobs
     * @name JobsVideosList
     * @summary List videos found by job
     * @request GET:/api/jobs/{id}/videos
     */
    jobsVideosList: (id: string, params: RequestParams = {}) =>
      this.request<JobAndVideos, ErrorModel>({
        path: `/api/jobs/${id}/videos`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags References
     * @name ReferencesCreate
     * @summary Upload reference images
     * @request POST:/api/references
     */
    referencesCreate: (
      data: {
        /** @format binary */
        file: File;
        /** @format int64 */
        minSURFMatches: number;
        /** @format double */
        minThresholdForSURFMatches: number;
        /** @format double */
        mseSkip: number;
        /** @format double */
        ratioTestThreshold: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<ReferenceUploadResponseBody, ErrorModel>({
        path: `/api/references`,
        method: "POST",
        body: data,
        type: ContentType.FormData,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags References
     * @name ReferencesDetail
     * @summary Get reference by ID
     * @request GET:/api/references/{id}
     */
    referencesDetail: (id: string, params: RequestParams = {}) =>
      this.request<ReferenceGetFeatureResponse, ErrorModel>({
        path: `/api/references/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Videos
     * @name VideosList
     * @summary List downloaded videos
     * @request GET:/api/videos
     */
    videosList: (params: RequestParams = {}) =>
      this.request<ListVideoBody[] | null, ErrorModel>({
        path: `/api/videos`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
  testsurf = {
    /**
     * No description
     *
     * @tags TestSurf
     * @name TestsurfList
     * @summary Get video metadata
     * @request GET:/testsurf
     */
    testsurfList: (params: RequestParams = {}) =>
      this.request<VideoMetadataBody, ErrorModel>({
        path: `/testsurf`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TestSurf
     * @name FrameList
     * @summary Retrieve a specific frame image
     * @request GET:/testsurf/frame
     */
    frameList: (
      query: {
        /** @format int64 */
        framenum: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorModel>({
        path: `/testsurf/frame`,
        method: "GET",
        query: query,
        ...params,
      }),

    /**
     * No description
     *
     * @tags TestSurf
     * @name ReferenceCreate
     * @summary Upload a reference file
     * @request POST:/testsurf/reference
     */
    referenceCreate: (
      data: {
        /** @format binary */
        video: File;
      },
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorModel>({
        path: `/testsurf/reference`,
        method: "POST",
        body: data,
        type: ContentType.FormData,
        ...params,
      }),

    /**
     * No description
     *
     * @tags TestSurf
     * @name TestList
     * @summary Test frame matching
     * @request GET:/testsurf/test
     */
    testList: (
      query: {
        /** @format int64 */
        framenum: number;
        /** @format double */
        ratiocheck: number;
        /** @format int64 */
        minmatches: number;
        /** @format double */
        goodmatchthreshold: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<FrameMatchingTestBody, ErrorModel>({
        path: `/testsurf/test`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TestSurf
     * @name VideoCreate
     * @summary Upload a video file
     * @request POST:/testsurf/video
     */
    videoCreate: (
      data: {
        /** @format binary */
        video: File;
      },
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorModel>({
        path: `/testsurf/video`,
        method: "POST",
        body: data,
        type: ContentType.FormData,
        ...params,
      }),
  };
}
