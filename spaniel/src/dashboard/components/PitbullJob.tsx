import styled from 'styled-components';

import {
  Divider,
  Grid,
  IconButton,
  Card,
  CardContent,
  Typography,
} from '@mui/material';

import { ContentCopy as ContentCopyIcon } from '@mui/icons-material';

import { PitbullJob as PitbullJobModel } from 'models';

import { Spacer } from 'common/components';

import { JobStatus } from 'core/components';

import { PitbullInfo } from './PitbullInfo';
import { PitbullJobInfo } from './PitbullJobInfo';

export type PitbullJobProps = {
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

export const PitbullJob: React.FC<PitbullJobProps> = props => {
  const { job } = props;
  
  const { instance } = job;

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
        <PitbullInfo instance={instance} />

        <Spacer mb={2} />
        <Divider />
        <Spacer mb={2} />

    
        <PitbullJobInfo job={job} />

      </CardContent>
    </Card>
  );
};
