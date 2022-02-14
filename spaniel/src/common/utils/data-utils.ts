import { Dictionary } from 'models';

import { isArray } from './base-type-guards';

export interface StripNullishOptions {
  skipNaN: boolean;
  skipUndefined: boolean;
  skipNull: boolean;
  skipFalse: boolean;
  skipEmptyString: boolean;
  skipZero: boolean;
}

const DEFAULT_STRIP_OPTIONS = {
  skipNaN: false,
  skipUndefined: false,
  skipNull: false,
  skipFalse: false,
  skipEmptyString: false,
  skipZero: false,
} as StripNullishOptions;

export const isNullish = (value: any, options: Partial<StripNullishOptions> = {} as StripNullishOptions) =>
  (!options.skipNaN && Number.isNaN(value))
  || (!options.skipUndefined && value === undefined)
  || (!options.skipNull && value === null)
  || (!options.skipFalse && value === false)
  || (!options.skipEmptyString && value === '')
  || (!options.skipZero && value === 0);

export const stripNullish = (
  target: Dictionary,
  options: Partial<StripNullishOptions> = {},
) => {
  const opts = {
    ...DEFAULT_STRIP_OPTIONS,
    ...options,
  } as StripNullishOptions;

  const res = Object.keys(target).reduce((acc, current) =>
    isNullish(target[current], opts)
      ? { ...acc }
      : { ...acc, [current]: target[current] },
    {}
  );

  return res;
};

export const ensureArrayValue = <TValue>(
  value: TValue | TValue[] | undefined
) => isArray(value)
    ? [...value]
    : value === undefined ? [] : [value];

export const toggleValueInArray = <TValue>(
  value: TValue,
  collection: TValue[],
) => collection.includes(value)
    ? collection.filter(v => v !== value)
    : [...collection, value];
