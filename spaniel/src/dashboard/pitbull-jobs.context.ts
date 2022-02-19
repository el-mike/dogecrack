import {
  createContext,
  useContext,
} from 'react';

import {
  PitbullJob,
  PitbullJobsFilters,
  RunPitbullJobPayload,
} from 'models';

export type ReloadJobsFn = () => void;
export type FilterFn = (filters: PitbullJobsFilters) => void;
export type ChangePageFn = (page: number) => void;
export type RunFn = (payload: RunPitbullJobPayload) => void;

export type PitbullJobsContext = {
  filters: PitbullJobsFilters;
  page: number;
  pageSize: number;
  totalCount: number;
  jobs: PitbullJob[];
  loading: boolean;
  lastLoaded: Date;
  reload: ReloadJobsFn;
  filter: FilterFn;
  changePage: ChangePageFn;
  run: RunFn;
}

export const pitbullJobsContext = createContext<PitbullJobsContext>(null!);

export const usePitbullJobs = () => useContext(pitbullJobsContext);
