import styled from 'styled-components';

import { Typography } from '@mui/material';

import { Spacer } from 'common/components';

import {
  CrackJobsProvider,
  CrackJobsList,
  RunCrackJob
} from '../components';

const DashboardWrapper = styled.div``;

export const DashboardPage: React.FC = () => {

  return (
    <DashboardWrapper>
      <CrackJobsProvider>
        <Typography variant='h4'>
          Dashboard
        </Typography>

        <Spacer mb={4} />

        <RunCrackJob />

        <Spacer mb={4} />

        <CrackJobsList />
      </CrackJobsProvider>
    </DashboardWrapper>
  );
};
