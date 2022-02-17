import styled from 'styled-components';

import { Typography } from '@mui/material';

import { ProgressInfo } from 'models';

import {
  Spacer,
  CircularProgressWithLabel,
} from 'common/components';

export type PitbullProgressProps = {
  progress: ProgressInfo;
}

const ProgressContainer = styled.div`
  display: flex;
  align-items: center;
`;

const ProgressInfoWrapper = styled.div``;

export const PitbullProgress: React.FC<PitbullProgressProps> = props => {
  const { progress } = props;

  const { checked, total } = progress;

  const percentage = !total
    ? 0
    : ((checked || 0) / total) * 100;

  return (
    <ProgressContainer>
      <CircularProgressWithLabel
        size={64}
        label={`${Math.round(percentage)}%`}
        value={percentage}
      />

      <Spacer mr={4} />

      <ProgressInfoWrapper>
        <Typography variant='caption'>Passwords checked:</Typography>
        <Typography variant='h5' fontWeight='bold'>{`${checked || 0} / ${total || '...'}`}</Typography>
      </ProgressInfoWrapper>

    </ProgressContainer>
  );
};

