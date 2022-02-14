import { useState } from 'react';

import styled from 'styled-components';

import {
  Box,
  Typography,
} from '@mui/material';

import { UserCredentials } from 'models';

import { LoginForm } from '../components';

const LoginBox = styled(Box)`
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: ${props => props.theme.spacing(8)};
`;

export const LoginPage: React.FC = () => {
  const [loading, setLoading] = useState(false);

  const login = (creds: UserCredentials) => console.log(creds);

  return (
    <LoginBox>
      <Typography variant='h5'>Log in</Typography>
      <LoginForm loading={loading} login={login} />
    </LoginBox>
  );
};
