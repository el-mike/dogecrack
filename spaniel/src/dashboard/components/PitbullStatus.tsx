import styled from 'styled-components';

import { Typography } from '@mui/material';

import { PitbullInstanceStatusKey } from 'models';

import {
  Spacer,
  CircularStatusIndicator,
} from 'common/components';

import { getLabelForEnum } from 'core/utils';

import { useGeneralContext } from 'core/contexts';

export type PitbullStatusProps = {
  status: number;
}

const ProgressContainer = styled.div`
  display: flex;
  align-items: center;
`;

const StatusInfoWrapper = styled.div``;

export const PitbullInstanceStatus: React.FC<PitbullStatusProps> = props => {
  const { status } = props;

  const { enums } = useGeneralContext();

  const { pitbullInstanceStatus: statusEnum } = enums;

  const label = getLabelForEnum(statusEnum, status);

  return (
    <ProgressContainer>
      <CircularStatusIndicator
        waiting={
          status === statusEnum[PitbullInstanceStatusKey.WAITING_FOR_HOST]
          || status === statusEnum[PitbullInstanceStatusKey.HOST_STARTING]
        }
        pending={status === statusEnum[PitbullInstanceStatusKey.RUNNING]}
        error={
          status === statusEnum[PitbullInstanceStatusKey.INTERRUPTED]
          || status === statusEnum[PitbullInstanceStatusKey.FAILED]
        }
        finished={status === statusEnum[PitbullInstanceStatusKey.COMPLETED]}
        success={status === statusEnum[PitbullInstanceStatusKey.SUCCESS]}
      />

      <Spacer mr={4} />

      <StatusInfoWrapper>
        <Typography variant='caption'>Status:</Typography>
        <Typography variant='h5' fontWeight='bold'>{label || 'Starting Pitbull'}</Typography>
      </StatusInfoWrapper>

    </ProgressContainer>
  );
};

