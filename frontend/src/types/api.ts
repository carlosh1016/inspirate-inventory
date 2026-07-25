// Backend response envelopes. Success payloads are wrapped in { data } or
// { data, meta }; errors follow RFC 7807 Problem Details (see lib/errors.ts).

export interface Meta {
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ApiEnvelope<T> {
  data: T;
}

export interface ApiListEnvelope<T> {
  data: T[];
  meta: Meta;
}
