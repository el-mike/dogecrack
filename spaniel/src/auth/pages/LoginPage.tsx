import styled from 'styled-components';

import {
  Box,
  Typography,
} from '@mui/material';

import { LoginForm } from '../components';

import { useAuthContext } from '../auth.context';

const LoginBox = styled(Box)`
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: ${props => props.theme.spacing(8)};
`;

export const LoginPage: React.FC = () => {
  const {
    loginLoading,
    login
  } = useAuthContext();

  return (
    <LoginBox>
      <Typography variant='h5'>Log in</Typography>
      <LoginForm loading={loginLoading} login={login} />
    </LoginBox>
  );
};
