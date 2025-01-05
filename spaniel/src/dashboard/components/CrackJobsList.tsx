import React from 'react';

import styled from 'styled-components';

import {
  Box,
  Typography,
  Pagination,
} from '@mui/material';

import { Spacer } from 'common/components';

import { TimeAgo } from 'core/components';

import { useCrackJobsContext } from '../crack-jobs.context';

import { CrackJob } from './CrackJob';
import { CrackJobsFilters } from './CrackJobsFilters';

const HeaderWrapper = styled.div`
  display: flex;
  justify-content: space-between;
`;

const NoJobsWrapper = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
`;

export const CrackJobsList: React.FC = () => {
  const {
    jobs,
    lastLoaded,
    totalCount,
    pageSize,
    page,
    changePage,
  } = useCrackJobsContext();


  const handlePageChange = (_: React.ChangeEvent<unknown>, value: number) => {
    changePage(value);
  }

  return (
    <Box>
      <HeaderWrapper>
        <Typography variant='h5'>Crack Jobs</Typography>
        <Typography variant='h6'>Last refreshed: <TimeAgo from={lastLoaded.toISOString()} /></Typography>
      </HeaderWrapper>

      <Spacer mb={4} />

      <CrackJobsFilters />

      <Spacer mb={4} />

      {!jobs.length
        ? <NoJobsWrapper>
            <Typography variant='h5'>No jobs found.</Typography>
          </NoJobsWrapper>
        : jobs.map(job => (
          <React.Fragment key={job.id}>
            <CrackJob job={job} />

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
