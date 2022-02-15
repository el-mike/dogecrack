import { Outlet } from 'react-router-dom';

import { Container } from '@mui/material';

export const UnauthenticatedLayout: React.FC = () => {
  return (
    <Container component='main' maxWidth='xs'>
      <Outlet />
    </Container>
  );
};
