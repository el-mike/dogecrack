import React, {
  useEffect,
  useState,
} from 'react';

import styled from 'styled-components';

import { Typography } from '@mui/material';

import { millisecondsInSecond } from 'date-fns';

import { Spacer } from 'common/components';

import { TimeAgo } from 'core/components';

 import { usePitbullJobs } from '../pitbull-jobs.context';

import { PitbullJob } from './PitbullJob';

const HeaderWrapper = styled.div`
  display: flex;
  justify-content: space-between;
`;

const NoJobsWrapper = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
`;

export const PitbullJobsList: React.FC = () => {
  const [refreshedAt, setRefreshedAt] = useState(new Date());
  
  const {
    jobs,
    loading,
    load,
  } = usePitbullJobs();

  useEffect(() => {
    load();

    const interval = setInterval(
      () => {
        load();
        setRefreshedAt(new Date());
      },
      /**
       * Every 30 seconds.
       */
      millisecondsInSecond * 30,
    );

    return () => clearInterval(interval);
  }, []);

  return (
    <>
    <HeaderWrapper>
      <Typography variant='h5'>Pitbull Jobs</Typography>
      <Typography variant='h6'>Last refreshed: <TimeAgo from={refreshedAt.toISOString()} /></Typography>
    </HeaderWrapper>

    <Spacer mb={4} />

    {!jobs.length
      ? <NoJobsWrapper>
          <Typography variant='h5'>No jobs found.</Typography>
        </NoJobsWrapper>
      : jobs.map(job => (
        <React.Fragment key={job.id}>
          <PitbullJob job={job} />
  
          <Spacer mb={3} />
        </React.Fragment>
      ))
    }
    </>
  );
};
