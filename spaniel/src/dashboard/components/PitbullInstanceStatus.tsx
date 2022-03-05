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

const StatusContainer = styled.div`
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
    <StatusContainer>
      <CircularStatusIndicator
        waiting={status === statusEnum[PitbullInstanceStatusKey.WAITING_FOR_HOST]}
        pending={
          status === statusEnum[PitbullInstanceStatusKey.RUNNING]
          || status === statusEnum[PitbullInstanceStatusKey.HOST_STARTING]
        }
        error={status === statusEnum[PitbullInstanceStatusKey.FAILED]}
        finished={status === statusEnum[PitbullInstanceStatusKey.COMPLETED]}
        success={status === statusEnum[PitbullInstanceStatusKey.SUCCESS]}
      />

      <Spacer mr={4} />

      <StatusInfoWrapper>
        <Typography variant='caption'>Status:</Typography>
        <Typography variant='h5' fontWeight='bold'>{label || 'Waiting for host'}</Typography>
      </StatusInfoWrapper>

    </StatusContainer>
  );
};

