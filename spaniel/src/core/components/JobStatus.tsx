import styled from 'styled-components';

import { Typography } from '@mui/material';

import { JobStatusKey } from 'models';

import { Spacer } from 'common/components';

import { useGeneralContext } from '../contexts';

import { getLabelForEnum } from '../utils';

import { JobStatusLight } from './JobStatusLight';

export type JobStatusProps = {
  status: number;
};

const StatusContainer = styled.div`
  display: flex;
  align-items: center;
`;

export const JobStatus: React.FC<JobStatusProps> = props => {
  const { status } = props;

  const { enums } = useGeneralContext();

  const { jobStatus: statusEnum } = enums;

  const label = getLabelForEnum(statusEnum, status);

  return (
    <StatusContainer>
      <Typography display='flex'>{label}</Typography>

      <Spacer mr={2} />

      <JobStatusLight
        $queued={status === statusEnum[JobStatusKey.SCHEDULED] || status === statusEnum[JobStatusKey.RESCHEDULED]}
        $processing={status === statusEnum[JobStatusKey.PROCESSING]}
        $rejected={status === statusEnum[JobStatusKey.REJECTED]}
        $acknowledged={status === statusEnum[JobStatusKey.ACKNOWLEDGED]}
      />
    </StatusContainer>
  );
};

