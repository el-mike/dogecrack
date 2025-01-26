import styled from 'styled-components';

import {
  Typography,
  Box,
} from '@mui/material';

import { Spacer } from 'common/components';

import {
  SettingsProvider,
  Settings,
} from '../components';

const SettingsWrapper = styled(Box)``;

export const SettingsPage: React.FC = () => {
  return (
    <SettingsProvider>
      <SettingsWrapper>
        <Typography variant='h4'>
          Settings
        </Typography>

        <Spacer mb={4} />

        <Settings />

      </SettingsWrapper>
    </SettingsProvider>
  );
};
