import styled from 'styled-components';

import {
  CircularProgress,
  Typography,
} from '@mui/material';

import { Spacer } from 'common/components';

const ShellLoaderWrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
  margin: 0 auto;
  justify-content: center;
  align-items: center;
`;

export const ShellLoader: React.FC = () => {
  return (
    <ShellLoaderWrapper>
      <Typography variant='h1'>DOGECRACK</Typography>

      <Spacer mt={6} />
  
      <CircularProgress size={80} />
    </ShellLoaderWrapper>
  );
};
