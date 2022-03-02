import {
  Box,
  Grid,
} from '@mui/material';

import { CrackJob as CrackJobModel } from 'models';

import {
  Spacer,
  Terminal,
  Accordion,
} from 'common/components';

import { toDateTimeString } from 'core/utils';

import { timeForPipe } from 'core/pipes';

import { LabeledInfo } from 'common/components';

export type CrackJobInfoProps = {
  job: CrackJobModel;
};

export const CrackJobInfo: React.FC<CrackJobInfoProps> = props => {
  const { job } = props;

  const lastFinishedAt = job.acknowledgedAt || job.rejectedAt;

  return (
    <Accordion title='Job info'>
      
      <Grid container spacing={2}>
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
    </Accordion>
  );
};
