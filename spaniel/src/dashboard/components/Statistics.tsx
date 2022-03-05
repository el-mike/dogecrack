import styled from 'styled-components';

import {
  Grid,
} from '@mui/material';

import { useStatisticsContext } from '../statistics.context';

import { CrackJobsStatistics } from './CrackJobsStatistics';
import { PitbullInstancesStatistics } from './PitbullInstancesStatistics';
import { OverviewStatistics } from './OverviewStatistics';

export const Statistics: React.FC = () => {
  const { statistics } = useStatisticsContext();

  return (
    !!statistics && (
      <Grid container spacing={3}>
        <Grid item xs={12} lg={3}>
          <CrackJobsStatistics statistics={statistics.crackJobs} />
        </Grid>

        <Grid item xs={12} lg={6}>
          <PitbullInstancesStatistics statistics={statistics.pitbullInstances} />
        </Grid>

        <Grid item xs={12} lg={3}>
          <OverviewStatistics statistics={statistics.pitbullInstances} />        
        </Grid>  
      </Grid>
    )
  );
}
