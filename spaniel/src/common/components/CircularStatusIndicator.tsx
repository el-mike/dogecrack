import styled from 'styled-components';

import {
  Box,
  BoxProps,
  SvgIcon,
} from '@mui/material';

import {
  Done as DoneIcon,
  DoneAll as DoneAllIcon,
  MoreHoriz as MoreHorizIcon,
  HourglassEmpty as HourglassEmptyIcon,
  PriorityHigh as PriorityHighIcon,
  QuestionMark as QuestionMarkIcon,
} from '@mui/icons-material';

export type CircularStatusIndicatorProps = {
  waiting?: boolean;
  pending?: boolean;
  error?: boolean;
  finished?: boolean;
  success?: boolean;
}

type OuterBoxProps = BoxProps & {
  $waiting?: boolean;
  $pending?: boolean;
  $error?: boolean;
  $finished?: boolean;
  $success?: boolean;
  $unknown?: boolean;
};

const OuterBox = styled(Box)<OuterBoxProps>`
  width: ${props => props.theme.spacing(8)};
  height: ${props => props.theme.spacing(8)};
  position: relative;
  display: inline-flex;
  // It's set to 0.7 to match CircularLoadingIndicator from MUI.
  border-width: ${props => props.theme.spacing(0.7)};
  border-style: solid;
  border-radius: 50%;

  ${props => props.$waiting && `
     color: ${props.theme.palette.warning.main};
    border-color: ${props.theme.palette.warning.main};
  `}

  ${props => props.$pending && `
    color: ${props.theme.palette.info.main};
    border-color: ${props.theme.palette.info.main};
  `}

  ${props => props.$error && `
    color: ${props.theme.palette.error.main};
    border-color: ${props.theme.palette.error.main};
  `}

  ${props => props.$finished && `
    color: ${props.theme.palette.common.white};
    border-color: ${props.theme.palette.common.white};
  `}

  ${props => props.$success && `
    color: ${props.theme.palette.success.main};
    border-color: ${props.theme.palette.success.main};
  `}

${props => props.$unknown && `
    color: ${props.theme.palette.common.white};
    border-color: ${props.theme.palette.common.white};
  `}
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
export const CircularStatusIndicator: React.FC<CircularStatusIndicatorProps> = props => {
  const {
    waiting,
    pending,
    error,
    finished,
    success,
  } = props;

  const unknown = !waiting
    && !pending
    && !error
    && !finished
    && !success;

  return (
    <OuterBox
      $waiting={waiting}
      $pending={pending}
      $error={error}
      $finished={finished}
      $success={success}
      $unknown={unknown}
    >
      <InnerBox>
        {!!waiting && <HourglassEmptyIcon />}
        {!!pending && <MoreHorizIcon />}
        {!!error && <PriorityHighIcon />}
        {!!finished && <DoneIcon />}
        {!!success && <DoneAllIcon />}
        {!!unknown && <QuestionMarkIcon />}
      </InnerBox>
    </OuterBox>
  );
};
