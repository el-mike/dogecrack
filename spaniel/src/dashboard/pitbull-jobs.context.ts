import {
  createContext,
  useContext,
} from 'react';

import {
  PitbullJob,
} from 'models';

export type LoadJobsFn = () => void;

export type PitbullJobsContext = {
  jobs: PitbullJob[];
  loading: boolean;
  load: LoadJobsFn;
}

export const pitbullJobsContext = createContext<PitbullJobsContext>(null!);

export const usePitbullJobs = () => useContext(pitbullJobsContext);
