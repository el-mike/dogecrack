import styled from 'styled-components';

import {
  Divider,
  IconButton,
  Paper,
  Typography,
} from '@mui/material';

import { ContentCopy as ContentCopyIcon } from '@mui/icons-material';

import { CrackJob as PitbullJobModel } from 'models';

import {
  Spacer,
  CardHeader,
} from 'common/components';

import { JobStatus } from 'core/components';

import { PitbullInstanceInfo } from './PitbullInstanceInfo';
import { PitbullHostInfo } from './PitbullHostInfo';
import { PitbullOutput } from './PitbullOutput';
import { CrackJobInfo } from './CrackJobInfo';
import { CrackJobErrorLog } from './CrackJobErrorLog';
import { CrackJobTokens } from './CrackJobTokens';

export type CrackJobProps = {
  job: PitbullJobModel;
};
const JobIdWrapper = styled.div`
  display: flex;
`;

export const CrackJob: React.FC<CrackJobProps> = props => {
  const { job } = props;

  const { instance } = job;
  const { pitbull } = instance;

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

      <PitbullInstanceInfo instance={instance} />
      <PitbullHostInfo instance={instance} />
      <PitbullOutput output={pitbull?.lastOutput || ''} />
      <CrackJobInfo job={job} />
      {!!job.tokens?.length && <CrackJobTokens job={job} />}
      <CrackJobErrorLog errorLog={job?.errorLog || ''} />
    </Paper>
  );
};
