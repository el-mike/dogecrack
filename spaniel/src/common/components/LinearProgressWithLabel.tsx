import styled from 'styled-components';

import {
  Typography,
  Box,
  LinearProgress,
  LinearProgressProps,
  useTheme,
} from '@mui/material';

import { roundDecimals } from '../utils';

import { Spacer } from './Spacer';


const ProgressContainer = styled(Box)`
  display: flex;
  align-items: center;
  width: 100%;
`;

const ProgressWrapper = styled(Box)`
  flex: 1 1 auto;
`;

const LabelWrapper = styled(Box)`
  flex: 35px 0 auto;
`;


export type LinearProgressWithLabelProps = LinearProgressProps & {
  value: number;
};

export const LinearProgressWithLabel: React.FC<LinearProgressWithLabelProps> = props => {
  const { value, ...rest} = props;

  const color = (value >= 75)
    ? 'success'
    : (value >= 50)
      ? 'warning'
      : 'error';

  const theme = useTheme();

  return (
    <ProgressContainer>
      <ProgressWrapper>
        <LinearProgress
          /**
           * LinearProgress cannot be styled with styled-components for some reason,
           * therefore we us sx prop.
           */
          sx={{ height: theme.spacing(1.5) }}
          color={color}
          variant='determinate'
          value={value}
          {...rest}
        />
      </ProgressWrapper>

      <Spacer mr={1} />

      <LabelWrapper>
        <Typography variant='subtitle1'>{`${roundDecimals(value)}%`}</Typography>
      </LabelWrapper>
    </ProgressContainer>
  );
};
