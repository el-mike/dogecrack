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
export type ResetFiltersFn = () => void;
export type CancelFn = (jobId: string) => void;
export type RecreateFn = (jobId: string) => void;

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
  cancel: CancelFn;
  recreate: RecreateFn;
  resetFilters: ResetFiltersFn;
}

export const crackJobsContext = createContext<CrackJobsContext>(null!);

export const useCrackJobsContext = () => useContext(crackJobsContext);
