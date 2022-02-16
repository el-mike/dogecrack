import styled from 'styled-components';

import {
  Divider,
  Grid,
  IconButton,
} from '@mui/material';

import { ContentCopy as ContentCopyIcon } from '@mui/icons-material';

import { PitbullJob as PitbullJobModel } from 'models';

import {
  Card,
  CardContent,
  Typography,
} from '@mui/material';

import {
  Spacer,
  CollapsibleTerminal,
} from 'common/components';

import { toDateTimeString } from 'core/utils';

import { timeForPipe } from 'core/pipes';

import { JobStatus } from 'core/components';

import { PitbullProgress } from './PitbullProgress';
import { PitbullStatus } from './PitbullStatus';
import { LabeledInfo } from './LabeledInfo';
import { LastUpdated } from './LastUpdated';

export type PitbullJobsProps = {
  job: PitbullJobModel;
};

const CardHeader = styled(Grid)`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: ${props => props.theme.spacing(2)};
`;

const JobIdWrapper = styled.div`
  display: flex;
`;

export const PitbullJob: React.FC<PitbullJobsProps> = props => {
  const { job } = props;
  
  const { instance } = job;
  const { hostInstance } = instance;

  const lastFinishedAt = job.acknowledgedAt || job.rejectedAt;

  const handleCopyJobId = () => {
    navigator.clipboard.writeText(job.id);
  };

  return (
    <Card>
      <CardHeader>
        <JobIdWrapper>
          <Typography variant='subtitle1'>Job ID:</Typography>
          <Typography variant='subtitle1' fontWeight='bold'>&nbsp; {job.id}</Typography>
          
          <Spacer mr={1} />

          <IconButton onClick={handleCopyJobId} size='small'>
            <ContentCopyIcon fontSize='small' />
          </IconButton>
        </JobIdWrapper>

        <JobStatus status={job.status} />
      </CardHeader>

      <Divider />
      
      <CardContent>
      <Grid container>
          <Grid item xs={12}>
            <Typography variant='overline'>Pitbull info</Typography>
          </Grid>

          <Grid container>
            <Grid container spacing={2} item xs={12} md={6} lg={4}>
              <Grid item xs={12}>
                <PitbullStatus status={instance?.status || -1} />
              </Grid>
            
              <Spacer mb={2} />
            
              <Grid item xs={12}>
                <PitbullProgress progress={job.instance?.progress || {}} />
              </Grid>
            </Grid>

            <Grid container spacing={2} item xs={12} md={6} lg={8}>
              <Grid item xs={6} md={4}>
                <LastUpdated updatedAt={instance?.updatedAt} />
              </Grid>
        
              <Grid item xs={6} md={4}>
                <LabeledInfo
                  title='Provider:'
                  value={instance?.providerName || '-'}
                />
              </Grid>

              <Grid item xs={6} md={4}>
                <LabeledInfo
                  title='Host address:'
                  value={
                    hostInstance?.sshHost
                      ? `${hostInstance?.sshHost}:${hostInstance.sshPort}`
                      : '-'
                  }
                  toCopy={
                    hostInstance.sshHost
                      ? `ssh -p ${hostInstance.sshPort} root@${hostInstance.sshHost}`
                      : ''
                  }
                />
              </Grid>

              <Grid item xs={12}>
                <LabeledInfo
                  title='Passlist URL:'
                  value={instance?.passlistUrl}
                  toCopy={instance?.passlistUrl}
                />
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      
        <Grid container>
          <Grid item xs={12}>
            <CollapsibleTerminal
              title='Pitbull output'
              content={instance?.lastOutput || ''}
            />
          </Grid>
        </Grid>

        <Spacer mb={2} />
        <Divider />
        <Spacer mb={2} />

    
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
                value={
                  lastFinishedAt ? timeForPipe(job.lastScheduledAt, lastFinishedAt) : '-'
                }
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

      </CardContent>
    </Card>
  );
};
