import React, { useMemo } from 'react';

import debounce from 'lodash.debounce';

export type HandlerFn = (event: React.ChangeEvent<HTMLInputElement>) => void;

export const useDebouncedInput = (handler: HandlerFn, wait: number, ...deps: any[]) =>
  useMemo(
    () => debounce(handler, wait),
    /* eslint-disable-next-line */
    [...deps],
  );
