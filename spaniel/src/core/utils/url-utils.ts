import queryString from 'query-string';

import { Dictionary } from 'models';

import {
  stripNullish,
  StripNullishOptions,
} from 'common/utils';

/**
 * There are some filters that take null, true/false, zero and empty strings as
 * correct values, therefore we need to skip stripping them for query params.
 */
const STRIP_NULLISH_OPTIONS = {
  skipNull: true,
  skipFalse: true,
  skipEmptyString: true,
  skipZero: true,
} as StripNullishOptions;

const DEFAULT_STRINGIFY_OPTIONS = {
  skipEmptyString: true,
} as queryString.StringifyOptions;

const DEFAULT_PARSE_OPTIONS = {
  parseBooleans: true,
  parseNumbers: true,
} as queryString.ParseOptions;

export const toQueryParams = (
  target: Dictionary,
  options: queryString.StringifyOptions = {}
) => `?${queryString.stringify(stripNullish(target, STRIP_NULLISH_OPTIONS), {
  ...DEFAULT_STRINGIFY_OPTIONS,
  ...options,
})}`;

export const fromQueryParams = (
  query: string,
  options: queryString.StringifyOptions = {},
) => queryString.parse(query, {
  ...DEFAULT_PARSE_OPTIONS,
  ...options,
});
