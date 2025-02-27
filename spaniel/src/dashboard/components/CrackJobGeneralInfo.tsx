import styled from 'styled-components';

import {
  Grid,
  Box,
  Typography,
  Tooltip,
} from '@mui/material';

import { Report as ReportIcon } from '@mui/icons-material';

import { CrackJob } from 'models';

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
  job: CrackJob;
};

const InfoWrapper = styled(Box)`
  padding: ${props => props.theme.spacing(2)};
`;

export const CrackJobGeneralInfo: React.FC<PitbullInfoProps> = props => {
  const { job } = props;
  const { instance } = job;
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
            {!!job.keyword && (
              <Grid item xs={6} md={3}>
                <LabeledInfo
                  title='Keyword:'
                  value={job.keyword}
                  endAdornment={
                    job.customTokenlist ? (
                      <Tooltip title='This job has been run with custom tokenlist!'>
                        <ReportIcon fontSize='small' color='warning' />
                      </Tooltip>
                    ) : <></>
                  }
                />
              </Grid>
            )}
            {!!job.passlistUrl && (
              <Grid item xs={6} md={3}>
                <LabeledInfo
                  title="Passlist URL:"
                  value={job.passlistUrl}
                  toCopy={job.passlistUrl}
                />
              </Grid>
              )}
          <Grid item xs={6} md={3}>
            <LabeledInfo title='Last updated:'>
              <TimeAgo from={instance?.updatedAt} />
            </LabeledInfo>
          </Grid>

          <Grid item xs={6} md={3}>
            <LabeledInfo
              title='Started at:'
              value={instance.startedAt && toDateTimeString(new Date(instance.startedAt))}
            />
          </Grid>

          <Grid item xs={6} md={3}>
            <LabeledInfo
              title='Completed at:'
              value={instance.completedAt && toDateTimeString(new Date(instance.completedAt))}
            />
          </Grid>

          <Grid item xs={6} md={3}>
            <LabeledInfo title='Run for:'>
              {
                !!instance.startedAt
                  ? (<TimeFor from={instance.startedAt} to={instance.completedAt} />)
                  : '-'
              }
            </LabeledInfo>
          </Grid>

          <Grid item xs={6} md={3}>
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
