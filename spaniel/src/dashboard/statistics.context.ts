import {
  createContext,
  useContext,
} from 'react';

import {
  CrackJobsStatistics,
  PitbullInstancesStatistics,
} from 'models';

export type LoadFn = () => void;

export type StatisticsContext = {
  load: LoadFn;
  crackJobs: CrackJobsStatistics;
  pitbullInstances: PitbullInstancesStatistics;
};

export const statisticsContext = createContext<StatisticsContext>(null!);

export const useStatisticsContext = () => useContext(statisticsContext);
