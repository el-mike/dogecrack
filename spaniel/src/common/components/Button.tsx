import {
  Button as MatButton,
  ButtonProps as MatButtonProps,
  CircularProgress,
} from '@mui/material';

export type ButtonProps = MatButtonProps & {
  loading?: boolean;
}

/**
 * Size is arbitrary - 26px is a size that does not require default button
 * to shrink or grow.
 * 
 * @TODO:
 * Find more reliable way to force CircularProgress fit Button component.
 */
const CIRCULAR_PROGRESS_SIZE = 26;

export const Button: React.FC<ButtonProps> = props => {
  const { loading, disabled, ...rest } = props;

  return (
    <MatButton {...rest} disabled={loading || disabled}>
      {
        loading
          ? <CircularProgress size={CIRCULAR_PROGRESS_SIZE} />
          : props.children
      }
    </MatButton>
  );
};

