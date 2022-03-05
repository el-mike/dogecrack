import {
  createContext,
  useContext,
} from 'react';

import {
  CrackJob,
  CrackJobsFilters,
  RunCrackJobPayload,
} from 'models';

export type ReloadJobsFn = () => void;
export type FilterFn = (filters: CrackJobsFilters) => void;
export type ChangePageFn = (page: number) => void;
export type RunFn = (payload: RunCrackJobPayload) => void;

export type CrackJobsContext = {
  filters: CrackJobsFilters;
  page: number;
  pageSize: number;
  totalCount: number;
  jobs: CrackJob[];
  loading: boolean;
  lastLoaded: Date;
  reload: ReloadJobsFn;
  filter: FilterFn;
  changePage: ChangePageFn;
  run: RunFn;
}

export const crackJobsContext = createContext<CrackJobsContext>(null!);

export const useCrackJobsContext = () => useContext(crackJobsContext);
