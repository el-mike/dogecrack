import {
  useState,
  useEffect,
} from 'react';

import {
  CrackJobsStatistics,
  PitbullInstancesStatistics,
} from 'models';

import {
  usePitbullInstanceService,
  useCrackJobService,
} from 'core/hooks';

import {
  statisticsContext,
  StatisticsContext,
} from '../statistics.context';

export const StatisticsProvider: React.FC = props => {
  const crackJobService = useCrackJobService();
  const pitbullInstanceService = usePitbullInstanceService();

  const [crackJobs, setCrackJobs] =
    useState<CrackJobsStatistics>({} as CrackJobsStatistics);

  const [pitbullInstances, setPitbullInstances] =
    useState<PitbullInstancesStatistics>({} as PitbullInstancesStatistics);

  const load = () => {
    Promise.all([
      crackJobService.getStatistics(),
      pitbullInstanceService.getStatistics(),
    ]).then(([crackJobsStatistics, pitbullInstancesStatistics]) => {
      setCrackJobs(crackJobsStatistics);
      setPitbullInstances(pitbullInstancesStatistics);
    });
  };

  const value = {
    load,
    crackJobs,
    pitbullInstances,
  } as StatisticsContext;

  useEffect(() => {
    load();
  }, []);

  return (
    <statisticsContext.Provider value={value}>{props.children}</statisticsContext.Provider>
  );
};
