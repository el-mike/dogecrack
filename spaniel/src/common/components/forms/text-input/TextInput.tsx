import {
  TextField,
  TextFieldProps,
} from '@mui/material';

export type TextInputProps = TextFieldProps & {}

export const TextInput: React.FC<TextInputProps> = props => {
  return (
    <TextField {...props} fullWidth />
  );
};
