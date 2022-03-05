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
  CrackJobsStatistics as CrackJobsStatisticsModel,
} from 'models';

import {
  CardHeader,
  LabeledInfo,
  LinearProgressWithLabel,
} from 'common/components';

const SummaryRow = styled(Box)`
  display: flex;
  justify-content: space-between;
  padding: ${props => props.theme.spacing(1)} ${props => props.theme.spacing(2)};
`;

const SummaryRowPart = styled(Box)`
  display: flex;
`;

const CompletionRateContainer = styled(SummaryRowPart)`
  display: flex;
  width: ${props => props.theme.spacing(15)};
`;

export type CrackJobsStatisticsProps = {
  statistics: CrackJobsStatisticsModel;
}

export const CrackJobsStatistics: React.FC<CrackJobsStatisticsProps> = props => {
  const {
    statistics,
  } = props;

  const completed = statistics.acknowledged + statistics.rejected;

  const successRate = completed === 0
    ? 0
    : (statistics.acknowledged / completed) * 100;

  return (
    <Card>
      <CardHeader>
        <Typography variant='subtitle1'>Crack Jobs Statistics</Typography>
      </CardHeader>

      <Divider />

      <CardContent>
        <Grid container spacing={3}>
          <Grid item xs={12} sm={6}>
            <LabeledInfo
              title='Acknowledged'
              value={statistics.acknowledged}
              valueVariant='h5'
              useZero={true}
            />
          </Grid>

          <Grid item xs={12} sm={6}>
            <LabeledInfo
              title='Rejected'
              value={statistics.rejected}
              valueVariant='h5'
              useZero={true}
            />
          </Grid>

          <Grid item xs={12} sm={6}>
            <LabeledInfo
              title='Queued'
              value={statistics.queued}
              valueVariant='h5'
              useZero={true}
              />
          </Grid>

          <Grid item xs={12} sm={6}>
            <LabeledInfo
              title='Processing'
              value={statistics.processing}
              valueVariant='h5'
              useZero={true}
            />
          </Grid>
          
        </Grid>
      </CardContent>

      <Divider />

      <SummaryRow>
        <SummaryRowPart>
          <Typography variant='subtitle1'>All</Typography>
        </SummaryRowPart>

        <SummaryRowPart>
          <Typography variant='subtitle1' fontWeight='bold'>{statistics.all}</Typography>
        </SummaryRowPart>

      </SummaryRow>

      <Divider />

      <SummaryRow>
        <SummaryRowPart>
          <Typography variant='subtitle1'>Acknowledge rate</Typography>
        </SummaryRowPart>

        <CompletionRateContainer>
          <LinearProgressWithLabel value={successRate} />
        </CompletionRateContainer>

      </SummaryRow>
    
    </Card>
  );
}
