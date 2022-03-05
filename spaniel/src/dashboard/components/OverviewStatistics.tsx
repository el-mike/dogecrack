import styled from 'styled-components';

import {
  Grid,
  Card,
  CardContent,
  Typography,
  Divider,
  Box,
} from '@mui/material';

import {
  PitbullInstancesStatistics as PitbullInstancesStatisticsModel,
} from 'models';

import { roundDecimals } from 'common/utils';

import {
  CardHeader,
  LabeledInfo,
} from 'common/components';

const StatisticContainer = styled(Box)`
  display: flex;
  align-items: center;
`;

export type OverviewStatisticsProps = {
  statistics: PitbullInstancesStatisticsModel;
}

export const OverviewStatistics: React.FC<OverviewStatisticsProps> = props => {
  const {
    statistics,
  } = props;

  return (
    <Card>
      <CardHeader>
        <Typography variant='subtitle1'>General Overview</Typography>
      </CardHeader>

      <Divider />

      <CardContent>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <StatisticContainer>
              <LabeledInfo
                title='Passwords checked'
                value={statistics.passwordsChecked}
                valueVariant='h5'
                useZero={true}
              />
            </StatisticContainer>
          </Grid>

          <Grid item xs={12} sm={6}>
            <StatisticContainer>            
              <LabeledInfo
                title='Total cost'
                value={`${roundDecimals(statistics.totalCost)} $`}
                valueVariant='h5'
                useZero={true}
              />
            </StatisticContainer>
          </Grid>

          <Grid item xs={12} sm={6}>
            <StatisticContainer>
              <LabeledInfo
                title='Avergage cost'
                value={`${roundDecimals(statistics.averageCost)} $`}
                valueVariant='h5'
                useZero={true}
              />
            </StatisticContainer>
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  );
}
