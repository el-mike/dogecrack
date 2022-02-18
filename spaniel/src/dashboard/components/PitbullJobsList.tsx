import React, {
  useEffect,
} from 'react';

import styled from 'styled-components';

import {
  Box,
  Typography,
  Pagination,
} from '@mui/material';

import { millisecondsInSecond } from 'date-fns';

import { Spacer } from 'common/components';

import { TimeAgo } from 'core/components';

 import { usePitbullJobs } from '../pitbull-jobs.context';

import { PitbullJob } from './PitbullJob';
import { PitbullJobsFilters } from './PitbullJobsFilters';

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
  const {
    jobs,
    reload,
    lastLoaded,
    totalCount,
    pageSize,
    page,
    changePage,
  } = usePitbullJobs();

  useEffect(() => {
    const interval = setInterval(
      () => {
        reload();
      },
      /**
       * Every 30 seconds.
       */
      millisecondsInSecond * 30,
    );

    return () => clearInterval(interval);
    /* eslint-disable-next-line */
  }, [reload]);

  const handlePageChange = (_: React.ChangeEvent<unknown>, value: number) => {
    changePage(value);
  }

  return (
    <Box>
      <HeaderWrapper>
        <Typography variant='h5'>Pitbull Jobs</Typography>
        <Typography variant='h6'>Last refreshed: <TimeAgo from={lastLoaded.toISOString()} /></Typography>
      </HeaderWrapper>

      <Spacer mb={4} />

      <PitbullJobsFilters />

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

      <Spacer mb={4} />

      <Box display='flex' justifyContent='flex-end'>
      <Pagination page={page} count={Math.floor((totalCount || 0) / pageSize)} onChange={handlePageChange} />
      </Box>

    </Box>
  );
};
