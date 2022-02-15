import { useState } from 'react';

import { PitbullJob } from 'models';

import {
  pitbullJobsContext,
  PitbullJobsContext,
} from '../pitbull-jobs.context';

import { usePitbullJobService } from 'core/hooks';

export const PitbullJobsProvider: React.FC = props => {
  const pitbullJobService = usePitbullJobService();

  const [jobs, setJobs] = useState<PitbullJob[]>([]);
  const [loading, setLoading] = useState(false);

  const load = () => {
    setLoading(true);
  
    pitbullJobService.getJobs()
      .then(jobs => setJobs(jobs || []))
      .finally(() => setLoading(false));
  };

  const value = {
    jobs,
    loading,
    load,
  } as PitbullJobsContext;

  return (
    <pitbullJobsContext.Provider value={value}>{props.children}</pitbullJobsContext.Provider>
  );
};
