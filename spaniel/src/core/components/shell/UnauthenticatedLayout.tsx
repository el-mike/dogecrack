import styled from 'styled-components';

import { Container } from '@mui/material';

export const UnauthenticatedLayout: React.FC = props => {
  return (
    <Container component='main' maxWidth='xs'>
      {props.children}
    </Container>
  );
};
