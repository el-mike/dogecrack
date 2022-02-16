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

import { Spacer } from 'common/components';

import { toDateTimeString } from 'core/utils';

import { JobStatus } from 'core/components';

import { PitbullProgress } from './PitbullProgress';
import { PitbullStatus } from './PitbullStatus';
import { PitbullOutput } from './PitbullOutput';
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

  const firstScheduledAt = toDateTimeString(new Date(job.firstScheduledAt));
  const jobUpdatedAt = toDateTimeString(new Date(job.updatedAt));

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
                value={firstScheduledAt}
              />
            </Grid>

            <Grid item xs={6} md={3}>
              <LabeledInfo
                title='updated at:'
                value={jobUpdatedAt}
              />
            </Grid>

            <Grid item xs={6} md={3}>
              <LabeledInfo
                title='Rescheduled count:'
                value={job.rescheduleCount}
              />
            </Grid>
          </Grid>
        </Grid>

        <Spacer mb={2} />
        <Divider />
        <Spacer mb={2} />

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
                  value={instance?.providerName}
                />
              </Grid>

              <Grid item xs={6} md={4}>
                <LabeledInfo
                  title='Host address:'
                  value={
                    hostInstance?.sshHost
                    ? `${hostInstance?.sshHost}:${hostInstance.sshPort}`
                    : ''
                  }
                  allowCopy={true}
                />
              </Grid>

              <Grid item xs={12}>
                <LabeledInfo
                  title='Passlist URL:'
                  value={instance?.passlistUrl}
                  allowCopy={true}
                />
              </Grid>

            </Grid>
          </Grid>
    
        </Grid>

        <Grid container>
          <Grid item xs={12}>
            <PitbullOutput output={instance?.lastOutput || ''} />
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  );
};
