import {
  Box,
  BoxProps,
} from '@mui/material';

export type SpacerProps = BoxProps['sx'];

export const Spacer: React.FC<SpacerProps> = props => {
  return (
    <Box sx={props}></Box>
  );
};
