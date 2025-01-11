import styled from 'styled-components';

import {
  Divider,
  IconButton,
  Paper,
  Typography,
} from '@mui/material';

import { ContentCopy as ContentCopyIcon } from '@mui/icons-material';

import { CrackJob as CrackJobModel } from 'models';

import {
  Spacer,
  CardHeader,
} from 'common/components';

import { JobStatus } from 'core/components';

import { CrackJobGeneralInfo } from './CrackJobGeneralInfo';
import { PitbullHostInfo } from './PitbullHostInfo';
import { PitbullOutput } from './PitbullOutput';
import { CrackJobInfo } from './CrackJobInfo';
import { CrackJobErrorLog } from './CrackJobErrorLog';
import { CrackJobTokens } from './CrackJobTokens';
import { CrackJobActionsMenu } from './CrackJobActionsMenu';

export type CrackJobProps = {
  job: CrackJobModel;
};

const JobIdWrapper = styled.div`
  display: flex;
`;

const StatusAndMenuWrapper = styled.div`
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

          {!!job.name && (
            <>
              <Spacer mr={2} />
              <Typography variant='subtitle1'>Name:</Typography>
              <Typography variant='subtitle1' fontWeight='bold'>&nbsp; {job.name}</Typography>
            </>
          )}
        </JobIdWrapper>

        <StatusAndMenuWrapper>
          <JobStatus status={job.status} />
          <Spacer mr={2} />
          <CrackJobActionsMenu job={job} />
        </StatusAndMenuWrapper>
      </CardHeader>

      <Divider />

      <CrackJobGeneralInfo job={job} />
      <CrackJobInfo job={job} />
      <PitbullHostInfo instance={instance} />
      <PitbullOutput output={pitbull?.lastOutput || ''} />
      {!!job.tokens?.length && <CrackJobTokens job={job} />}
      <CrackJobErrorLog errorLog={job?.errorLog || ''} />
    </Paper>
  );
};
