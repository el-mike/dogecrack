import styled from 'styled-components';

import { Typography } from '@mui/material';

import { Spacer } from 'common/components';

import {
  PitbullJobsProvider,
  PitbullJobsList,
} from '../components';

const DashboardWrapper = styled.div``;

export const DashboardPage: React.FC = () => {

  return (
    <DashboardWrapper>
      <PitbullJobsProvider>
        <Typography variant='h4'>
          Dashboard
        </Typography>

        <Spacer mb={4} />

        <PitbullJobsList />
      </PitbullJobsProvider>
    </DashboardWrapper>
  );
};
