import styled from 'styled-components';

import {
  Typography,
  Box,
} from '@mui/material';

import { Spacer } from 'common/components';
import {
  CheckedIdeasProvider,
  CheckedIdeas,
} from '../components';

const CheckedIdeasWrapper = styled(Box)``;

export const CheckedIdeasPage: React.FC = () => {
  return (
    <CheckedIdeasProvider>
      <CheckedIdeasWrapper>
        <Typography variant='h4'>
          Passwords
        </Typography>

        <Spacer mb={4} />

        <CheckedIdeas />

      </CheckedIdeasWrapper>
    </CheckedIdeasProvider>
  );
};
