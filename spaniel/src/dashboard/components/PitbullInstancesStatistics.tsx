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

import {
  Spacer,
  CardHeader,
  LabeledInfo,
  CircularStatusIndicator,
  LinearProgressWithLabel,
} from 'common/components';

const StatisticContainer = styled(Box)`
  display: flex;
  align-items: center;
`;


const SummaryRow = styled(Box)`
  display: flex;
  justify-content: space-between;
  padding: ${props => props.theme.spacing(1)} ${props => props.theme.spacing(2)};
`;

const SummaryRowPart = styled(Box)`
  display: flex;
`;

const AcknowledgeRateContainer = styled(SummaryRowPart)`
  display: flex;
  width: ${props => props.theme.spacing(30)};
`;

export type PitbullInstancesStatisticsProps = {
  statistics: PitbullInstancesStatisticsModel;
}

export const PitbullInstancesStatistics: React.FC<PitbullInstancesStatisticsProps> = props => {
  const {
    statistics,
  } = props;

  const allCoompleted = statistics.completed
    + statistics.failed
    + statistics.success;

  const successRate = allCoompleted === 0
    ? 0
    : ((statistics.completed + statistics.success) / allCoompleted) * 100;

  return (
    <Card>
      <CardHeader>
        <Typography variant='subtitle1'>Pitbull Instances Statistics</Typography>
      </CardHeader>

      <Divider />

      <CardContent>
        <Grid container spacing={3}>
          <Grid item xs={12} sm={4}>
            <StatisticContainer>
              <CircularStatusIndicator waiting={true} size='small' />

              <Spacer mr={2} />
            
              <LabeledInfo
                title='Waiting for host'
                value={statistics.waitingForHost}
                valueVariant='h5'
                useZero={true}
              />
            </StatisticContainer>
          </Grid>

          <Grid item xs={12} sm={4}>
            <StatisticContainer>
              <CircularStatusIndicator pending={true} size='small' />

              <Spacer mr={2} />
            
              <LabeledInfo
                title='Running / Host starting'
                value={statistics.running + statistics.hostStarting}
                valueVariant='h5'
                useZero={true}
              />
            </StatisticContainer>
          </Grid>

          {/* Ofsset element - we want to push "Completed" status to the next line. */}
          <Grid item xs={12} sm={4}></Grid>

          <Grid item xs={12} sm={4}>
            <StatisticContainer>
              <CircularStatusIndicator finished={true} size='small' />

              <Spacer mr={2} />
            
              <LabeledInfo
                title='Completed'
                value={statistics.completed}
                valueVariant='h5'
                useZero={true}
              />
            </StatisticContainer>
          </Grid>

          <Grid item xs={12} sm={4}>
            <StatisticContainer>
              <CircularStatusIndicator success={true} size='small' />

              <Spacer mr={2} />
            
              <LabeledInfo
                title='Success'
                value={statistics.success}
                valueVariant='h5'
                useZero={true}
              />
            </StatisticContainer>
          </Grid>

          <Grid item xs={12} sm={4}>
            <StatisticContainer>
              <CircularStatusIndicator error={true} size='small' />

              <Spacer mr={2} />
            
              <LabeledInfo
                title='Failed'
                value={statistics.failed}
                valueVariant='h5'
                useZero={true}
              />
            </StatisticContainer>
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
          <Typography variant='subtitle1'>Completion rate</Typography>
        </SummaryRowPart>

        <AcknowledgeRateContainer>
          <LinearProgressWithLabel value={successRate} />
        </AcknowledgeRateContainer>
      </SummaryRow>
    
    </Card>
  );
}
