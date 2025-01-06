import {
  useCallback,
  useState
} from 'react';
import {
  NotificationOptions,
  NotificationVariant,
  SnackbarContext,
  snackbarContext,
} from '../contexts';
import {
  Alert,
  Box,
  Snackbar
} from '@mui/material';

type SnackbarState = {
  open: boolean;
  message: string;
  variant: NotificationVariant;
};

export const SnackbarProvider: React.FC = props => {
  const [snackbarState, setSnackbarState] = useState<SnackbarState | null>(null);

  const notify = useCallback((options: NotificationOptions) => {
    setSnackbarState({
      open: true,
      message: options.message,
      variant: options.variant,
    });
  }, [setSnackbarState]);

  const hide = useCallback(() => {
    setSnackbarState(null);
  }, [setSnackbarState]);

  const value = {
    notify,
  } as SnackbarContext;

  return (
    <snackbarContext.Provider value={value}>
      <Box>
        {props.children}
        <Snackbar open={snackbarState?.open} autoHideDuration={5000} onClose={hide}>
          <Box>
            {snackbarState?.variant === NotificationVariant.SUCCESS && (
              <Alert severity='success' variant='filled' onClose={hide}>
                {snackbarState?.message}
              </Alert>
            )}
            {snackbarState?.variant === NotificationVariant.ERROR && (
              <Alert severity='error' variant='filled' onClose={hide}>
                {snackbarState?.message}
              </Alert>
            )}
          </Box>
        </Snackbar>
      </Box>
    </snackbarContext.Provider>
  );
};
