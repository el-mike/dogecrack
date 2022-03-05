import { useEffect } from 'react';

import styled from 'styled-components';

import {
  Typography,
  Grid,
} from '@mui/material';

import { millisecondsInSecond } from 'date-fns';

import { Spacer } from 'common/components';

import {
  CrackJobsProvider,
  StatisticsProvider,
  CrackJobsList,
  RunCrackJob,
  Statistics,
} from '../components';

import { useCrackJobsContext } from '../crack-jobs.context';
import { useStatisticsContext } from '../statistics.context';

const DashboardWrapper = styled.div``;

const InnerDashboardPage: React.FC = () => {
  const { reload: reloadJobs } = useCrackJobsContext();
  const { load: loadStatistics } = useStatisticsContext();

  useEffect(() => {
    const interval = setInterval(
      () => {
        reloadJobs();
        loadStatistics();
      },
      /**
       * Every 30 seconds.
       */
      millisecondsInSecond * 30,
    );

    return () => clearInterval(interval);
    /* eslint-disable-next-line */
  }, [reloadJobs]);

  return (
    <DashboardWrapper>
      <Typography variant='h4'>
        Dashboard
      </Typography>

      <Spacer mb={4} />

      <Grid container>
        <Grid item xs={12}>
          <Statistics />

          <Spacer mb={4} />
        
          <RunCrackJob />

          <Spacer mb={4} />
        
          <CrackJobsList />
        </Grid>

      </Grid>
    </DashboardWrapper>
  );
}

export const DashboardPage: React.FC = () => {
  return (
      <CrackJobsProvider>
        <StatisticsProvider>
          <InnerDashboardPage />
        </StatisticsProvider>
      </CrackJobsProvider>
  );
};
