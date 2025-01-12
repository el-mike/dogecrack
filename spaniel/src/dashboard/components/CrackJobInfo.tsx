import { Grid } from '@mui/material';

import {
  CrackJob as CrackJobModel,
  TokenGeneratorVersionKey
} from 'models';

import { Accordion, LabeledInfo } from 'common/components';

import {
  getLabelForEnum,
  toDateTimeString
} from 'core/utils';
import { timeForPipe } from 'core/pipes';
import { useGeneralContext } from 'core/contexts';

export type CrackJobInfoProps = {
  job: CrackJobModel;
};

export const CrackJobInfo: React.FC<CrackJobInfoProps> = props => {
  const { enums } = useGeneralContext();

  const { job } = props;

  const lastFinishedAt = job.acknowledgedAt || job.rejectedAt;

  return (
    <Accordion title='Job info'>
      <Grid container spacing={2}>
        <Grid container item xs={12}>
          <Grid item xs={6} md={2}>
            <LabeledInfo
              title='Token generator version:'
              value={getLabelForEnum(enums.tokenGeneratorVersion, job.tokenGeneratorVersion || 0)}
            />
          </Grid>

          <Grid item xs={6} md={2}>
            <LabeledInfo
              title='First scheduled at:'
              value={toDateTimeString(new Date(job.firstScheduledAt))}
            />
          </Grid>

          <Grid item xs={6} md={2}>
            <LabeledInfo
              title='updated at:'
              value={toDateTimeString(new Date(job.updatedAt))}
            />
          </Grid>

          <Grid item xs={6} md={2}>
            <LabeledInfo
              title='Run for:'
              value={lastFinishedAt && timeForPipe(job.lastScheduledAt, lastFinishedAt)}
            />
          </Grid>

          <Grid item xs={6} md={2}>
            <LabeledInfo
              title='Rescheduled count:'
              value={job.rescheduleCount}
            />
          </Grid>

          <Grid item xs={6} md={2}>
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
