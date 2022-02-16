import styled from 'styled-components';

import {
  Box,
  Typography,
  CircularProgress,
  CircularProgressProps,
} from '@mui/material';

const OuterBox = styled(Box)`
  position: relative;
  display: inline-flex;
`;

const InnerBox = styled(Box)`
  top: 0;
  left: 0;
  bottom: 0;
  right: 0;
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
`;

export type CircularProgressWithLabelProps = CircularProgressProps & {
  label: string;
};

export const CircularProgressWithLabel: React.FC<CircularProgressWithLabelProps> = props => {
  const {
    label,
    ...rest
  } = props;

  return (
    <OuterBox>
      <CircularProgress variant='determinate' {...rest} />
      <InnerBox>
        <Typography variant='caption' component='div' color='textSecondary'>
          {label}
        </Typography>
      </InnerBox>
    </OuterBox>
  );
};
