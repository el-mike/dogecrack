import styled from 'styled-components';

import {
  Drawer,
  Box,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
} from '@mui/material';

import {
  Dashboard as DashboardIcon,
  ManageSearch as ManageSearchIcon,
} from '@mui/icons-material';

import { Link } from 'react-router-dom';

import { useNavigationContext } from '../navigation.context';

const NavigationContainer = styled(Box)`
  width: 200px;
`;

export const Navigation: React.FC = () => {
  const {
    isOpen,
    close,
  } = useNavigationContext();

  return (
    <Drawer open={isOpen} onClose={close}>
      <NavigationContainer onClick={close}>
        <List>
          <ListItem
            button
            component={Link}
            to='/dashboard'
          >
            <ListItemIcon>
              <DashboardIcon />
            </ListItemIcon>
            <ListItemText primary='Dashboard' />
          </ListItem>

          <ListItem
            button
            component={Link}
            to='/passchecks'
          >
            <ListItemIcon>
              <ManageSearchIcon />
            </ListItemIcon>
            <ListItemText primary='Passchecks' />
          </ListItem>
        </List>
      </NavigationContainer>
    </Drawer>
  );
};
