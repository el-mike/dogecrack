import React, {
  useEffect,
  useState,
} from 'react';

import styled from 'styled-components';

import { Typography } from '@mui/material';

import { millisecondsInSecond } from 'date-fns';

import { Spacer } from 'common/components';

import { useTimeAgo } from 'core/hooks';

 import { usePitbullJobs } from '../pitbull-jobs.context';

import { PitbullJob } from './PitbullJob';

const HeaderWrapper = styled.div`
  display: flex;
  justify-content: space-between;
`;

export const PitbullJobsList: React.FC = props => {
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

  const refreshedAgo = useTimeAgo(refreshedAt.toISOString());

  return (
    <>
    <HeaderWrapper>
      <Typography variant='h5'>Pitbull Jobs</Typography>
      <Typography variant='h6'>Last refreshed: {refreshedAgo}</Typography>
    </HeaderWrapper>

    <Spacer mb={2} />

    {jobs.map(job => (
      <React.Fragment key={job.id}>
        <PitbullJob job={job} />

        <Spacer mb={3} />
      </React.Fragment>
    ))}
    </>
  );
};
