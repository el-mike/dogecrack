import styled from 'styled-components';

import {
  Grid,
  Box,
  Typography,
} from '@mui/material';

import { PitbullInstance } from 'models';

import {
  Spacer,
  LabeledInfo,
} from 'common/components';

import {
  getInstanceEstimatedCost,
  toDateTimeString,
} from 'core/utils';

import {
  TimeFor,
  TimeAgo,
} from 'core/components';

import { PitbullProgress } from './PitbullProgress';
import { PitbullInstanceStatus } from './PitbullInstanceStatus';

export type PitbullInfoProps = {
  instance: PitbullInstance;
};

const InfoWrapper = styled(Box)`
  padding: ${props => props.theme.spacing(2)};
`;

export const PitbullInstanceInfo: React.FC<PitbullInfoProps> = props => {
  const { instance } = props;

  const { pitbull } = instance;

  return (
    <InfoWrapper>
      <Typography variant='overline'>Pitbull info</Typography>

      <Spacer mb={2} />

      <Grid container spacing={2}>
        <Grid container spacing={2} item xs={12} md={6} lg={4}>
          <Grid item xs={12}>
            <PitbullInstanceStatus status={instance?.status || -1} />
          </Grid>

          <Spacer mb={2} />

          <Grid item xs={12}>
            <PitbullProgress pitbull={pitbull || {}} />
          </Grid>
        </Grid>

        <Grid container spacing={2} item xs={12} md={6} lg={8}>
          <Grid item xs={6} md={4}>
            <LabeledInfo title='Last updated:'>
              <TimeAgo from={instance?.updatedAt} />
            </LabeledInfo>
          </Grid>

          <Grid item xs={6} md={4}>
            <LabeledInfo
              title='Started at:'
              value={instance.startedAt && toDateTimeString(new Date(instance.startedAt))}
            />
          </Grid>

          <Grid item xs={6} md={4}>
            <LabeledInfo
              title='Completed at:'
              value={instance.completedAt && toDateTimeString(new Date(instance.completedAt))}
            />
          </Grid>

          <Grid item xs={6} md={4}>
            <LabeledInfo title='Run for:'>
              {
                !!instance.startedAt
                  ? (<TimeFor from={instance.startedAt} to={instance.completedAt} />)
                  : '-'
              }
            </LabeledInfo>
          </Grid>

          <Grid item xs={6} md={4}>
            <LabeledInfo
              title='Estimated cost:'
              value={`${getInstanceEstimatedCost(instance)} $`}
            />
          </Grid>

          <Grid item xs={6} md={4}>
            <LabeledInfo
              title='Instance ID:'
              value={instance.id}
              toCopy={instance.id}
            />
          </Grid>
        </Grid>
      </Grid>
      </InfoWrapper>
  );
};
