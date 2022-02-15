import { useEffect } from 'react';
import { usePitbullJobs } from '../pitbull-jobs.context';

import { PitbullJob } from './PitbullJob';

export const PitbullJobsList: React.FC = props => {
  const {
    jobs,
    loading,
    load,
  } = usePitbullJobs();

  useEffect(() => {
    load();
  }, []);

  return (
    <>
      {jobs.map(job => (
        <PitbullJob key={job.id} job={job} />
      ))}
    </>
  );
};
