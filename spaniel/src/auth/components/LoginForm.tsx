import {
  Box,
  TextField,
} from '@mui/material';

import { UserCredentials } from 'models';

import {
  Spacer,
  Button,
} from 'common/components';

enum FormKeys {
  USERNAME = 'username',
  PASSWORD = 'password'
}

export type LoginFormProps = {
  loading: boolean;
  login: (creds: UserCredentials) => void;
};


export const LoginForm: React.FC<LoginFormProps> = props => {
  const { loading, login } = props;

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const data = new FormData(event.currentTarget);

    login({
      username: data.get('username'),
      password: data.get('password')
    } as UserCredentials);
  };

  return (
    <Box
      component='form'
      onSubmit={handleSubmit}
      noValidate
    >
      <TextField
        margin='normal'
        required
        fullWidth
        id={FormKeys.USERNAME}
        name={FormKeys.USERNAME}
        label='Username'
        autoFocus
      />

      <TextField
        margin='normal'
        required
        fullWidth
        id={FormKeys.PASSWORD}
        name={FormKeys.PASSWORD}
        label='Password'
        type='password'
        autoFocus
      />

      <Spacer mt={3} />

      <Button
        loading={loading}
        fullWidth
        type='submit'
        variant='contained'
      >
        Log in
      </Button>
    </Box>
  );
};
