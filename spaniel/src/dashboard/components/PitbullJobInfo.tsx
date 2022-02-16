import {
  Box,
  Grid,
} from '@mui/material';

import { PitbullJob as PitbullJobModel } from 'models';

import { Typography } from '@mui/material';

import {
  Spacer,
  CollapsibleTerminal,
} from 'common/components';

import { toDateTimeString } from 'core/utils';

import { timeForPipe } from 'core/pipes';

import { LabeledInfo } from 'common/components';

export type PitbullJobInfoProps = {
  job: PitbullJobModel;
};

export const PitbullJobInfo: React.FC<PitbullJobInfoProps> = props => {
  const { job } = props;

  const lastFinishedAt = job.acknowledgedAt || job.rejectedAt;

  return (
    <Box>      
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <Typography variant='overline'>Job info</Typography>
        </Grid>

        <Grid container item xs={12}>
          <Grid item xs={6} md={3}>
            <LabeledInfo
              title='Keyword:'
              value={job.keyword}
            />
          </Grid>

          <Grid item xs={6} md={3}>
            <LabeledInfo
              title='First scheduled at:'
              value={toDateTimeString(new Date(job.firstScheduledAt))}
            />
          </Grid>

          <Grid item xs={6} md={3}>
            <LabeledInfo
              title='updated at:'
              value={toDateTimeString(new Date(job.updatedAt))}
            />
          </Grid>

          <Grid item xs={6} md={3}>
            <LabeledInfo
              title='Rescheduled count:'
              value={job.rescheduleCount}
            />
          </Grid>
        </Grid>

        <Spacer mb={1} />

        <Grid container item xs={12}>
          <Grid item xs={6} md={3}>
            <LabeledInfo
              title='Run for:'
              value={lastFinishedAt && timeForPipe(job.lastScheduledAt, lastFinishedAt)}
            />
          </Grid>

          <Grid item xs={6} md={3}>
            {!!job.rejectedAt && (
              <LabeledInfo
                title='Rejected at:'
                value={toDateTimeString(new Date(job.rejectedAt))}
              />
            )}

            {!!job.acknowledgedAt && (
              <LabeledInfo
                title='Acknowledged at:'
                value={toDateTimeString(new Date(job.acknowledgedAt))}
              />
            )}
          </Grid>
        </Grid>
      </Grid>

      <Grid container>
        <Grid item xs={12}>
          <CollapsibleTerminal
            title='Job errors'
            content={job.errorLog || ''}
          />
        </Grid>
      </Grid>
    </Box>
  );
};
