import {
  FormControl,
  Select,
  SelectProps,
  InputLabel,
  MenuItem,
} from '@mui/material';

import { InputOption } from '../input-option';

export type SelectInputProps = SelectProps & {
  label: string;
  options: InputOption[];
};

const LABEL_ID = 'select-label';

export const SelectInput: React.FC<SelectInputProps> = props => {
  const { label, options, ...rest } = props;

  return (
    <FormControl fullWidth>
      <InputLabel id={LABEL_ID}>{label}</InputLabel>

      <Select
        labelId={LABEL_ID}
        label={label}
        displayEmpty
        {...rest}
      >
        {(options || []).map(option => (
          <MenuItem key={option.label} value={option.value}>
            {option.label}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};
