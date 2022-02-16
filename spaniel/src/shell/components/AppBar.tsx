import styled from 'styled-components';

import {
  AppBar as MuiAppBar,
  Toolbar,
  Typography,
  IconButton,
} from '@mui/material';

import {
  Logout as LogoutIcon,
  Menu as MenuIcon
} from '@mui/icons-material';

import { Spacer } from 'common/components';

import { useAuthContext } from 'auth';

import { useNavigationContext } from '../navigation.context';

const ActionsContainer = styled.div``;

export const AppBar: React.FC = () => {
  const { logout } = useAuthContext();
  const { open } = useNavigationContext();

  const handleNavigationOpen = () => open();
  const handleLogout = () => logout();

  return (
    <MuiAppBar>
      <Toolbar>
        <IconButton
          size='large'
          edge='start'
          color='inherit'
          aria-label='menu'
          onClick={handleNavigationOpen}
        >
          <MenuIcon />
        </IconButton>

        <Spacer mr={2} />

        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          DOGECRACK
        </Typography>

        <ActionsContainer>
        <IconButton
          size="large"
          aria-label="logout"
          aria-controls="menu-appbar"
          aria-haspopup="true"
          onClick={handleLogout}
          color="inherit"
        >
          <LogoutIcon />
        </IconButton>
        </ActionsContainer>
      </Toolbar>
    </MuiAppBar>
  );
};
