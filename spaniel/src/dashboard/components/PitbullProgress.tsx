import styled from 'styled-components';

import { Typography } from '@mui/material';

import {
  Pitbull,
  ProgressInfo
} from 'models';

import {
  Spacer,
  CircularProgressWithLabel,
} from 'common/components';

export type PitbullProgressProps = {
  pitbull: Pitbull;
}

const ProgressContainer = styled.div`
  display: flex;
  align-items: center;
`;

const ProgressInfoWrapper = styled.div`
  display: flex;
  flex-direction: column;
`;

export const PitbullProgress: React.FC<PitbullProgressProps> = props => {
  const { pitbull } = props;
  const { progress, skipCount: baseSkipCount } = pitbull;

  const { checked: baseChecked, total: baseTotal } = progress || {};

  const skipCount = baseSkipCount || 0;
  const checked = baseChecked + skipCount;
  const total = baseTotal + skipCount;

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
        {!!skipCount && <Typography variant='caption'>Resumed at: {skipCount + 1}</Typography>}
      </ProgressInfoWrapper>
    </ProgressContainer>
  );
};
