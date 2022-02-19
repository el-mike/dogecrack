import styled from 'styled-components';

import {
  Divider,
  IconButton,
  Card,
  Paper,
  CardContent,
  Typography,
} from '@mui/material';

import { ContentCopy as ContentCopyIcon } from '@mui/icons-material';

import { PitbullJob as PitbullJobModel } from 'models';

import {
  Spacer,
  CardHeader,
} from 'common/components';

import { JobStatus } from 'core/components';

import { PitbullInfo } from './PitbullInfo';
import { PitbullHostInfo } from './PitbullHostInfo';
import { PitbullOutput } from './PitbullOutput';
import { PitbullJobInfo } from './PitbullJobInfo';
import { PitbullJobErrorLog } from './PitbullJobErrorLog';

export type PitbullJobProps = {
  job: PitbullJobModel;
};
const JobIdWrapper = styled.div`
  display: flex;
`;

export const PitbullJob: React.FC<PitbullJobProps> = props => {
  const { job } = props;
  
  const { instance } = job;

  const handleCopyJobId = () => {
    navigator.clipboard.writeText(job.id);
  };

  return (
    <Paper>
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
      
      <PitbullInfo instance={instance} />
      <PitbullHostInfo instance={instance} />
      <PitbullOutput output={instance?.lastOutput || ''} />
      <PitbullJobInfo job={job} />
      <PitbullJobErrorLog errorLog={job?.errorLog || ''} />
    </Paper>
  );
};
