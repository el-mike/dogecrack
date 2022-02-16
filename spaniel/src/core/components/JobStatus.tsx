import styled from 'styled-components';

import { Typography } from '@mui/material';

import { JobStatusKey } from 'models';

import { Spacer } from 'common/components';

import { useGeneralContext } from '../contexts';

import { getLabelForEnum } from '../utils';

export type JobStatusProps = {
  status: number;
};

type StatusLightProps = {
  $queued: boolean;
  $processing: boolean;
  $rejected: boolean;
  $acknowledged: boolean;
};

const StatusContainer = styled.div`
  display: flex;
  align-items: center;
`;

const StatusLight = styled.span<StatusLightProps>`
  display: flex;
  width: ${props => props.theme.spacing(2)};
  height: ${props => props.theme.spacing(2)};
  border-radius: 50%;

  ${props => props.$queued && `
    background-color: ${props.theme.palette.warning.light};
  `}

  ${props => props.$processing && `
    background-color: ${props.theme.palette.info.light};
  `}

  ${props => props.$rejected && `
    background-color: ${props.theme.palette.error.light};
  `}

  ${props => props.$acknowledged && `
    background-color: ${props.theme.palette.success.light};
  `}
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

      <StatusLight
        $queued={status === statusEnum[JobStatusKey.JOB_SCHEDULED] || status === statusEnum[JobStatusKey.JOB_RESCHEDULED]}
        $processing={status === statusEnum[JobStatusKey.JOB_PROCESSING]}
        $rejected={status === statusEnum[JobStatusKey.JOB_REJECTED]}
        $acknowledged={status === statusEnum[JobStatusKey.JOB_ACKNOWLEDGED]}
      />
    </StatusContainer>
  );
};

