import styled from 'styled-components';

import { Outlet } from 'react-router-dom';

import { Container } from '@mui/material';

import { NavigationProvider } from './NavigationProvider';

import { AppBar } from './AppBar';
import { Navigation } from './Navigation';

const AppBarOffset = styled.div`
  ${props => ({ ...props.theme.mixins.toolbar } as any)};
`;

const ContentContainer = styled(Container)`
  padding: ${props => props.theme.spacing(3)};
  padding-top: ${props => props.theme.spacing(4)};
`;

export const AuthenticatedLayout: React.FC = props => {
  return (
  <NavigationProvider>
    <AppBar />

    <AppBarOffset />
    
    <Navigation />

    <ContentContainer maxWidth='xl'>
      <Outlet />
    </ContentContainer>
    
  </NavigationProvider>
  );
};
