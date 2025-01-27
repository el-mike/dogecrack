import { useState } from 'react';
import {
  Box,
  IconButton,
  Menu,
  MenuItem
} from '@mui/material';
import { MoreVert as MoreVertIcon } from '@mui/icons-material';

import { CrackJob, JobStatusKey } from 'models';
import { useGeneralContext } from '../../core/contexts';
import { useCrackJobsContext } from '../crack-jobs.context';

export type CrackJobActionsMenuProps = {
  job: CrackJob;
}

export const CrackJobActionsMenu = (props: CrackJobActionsMenuProps) => {
  const { job } = props;

  const { enums } = useGeneralContext();
  const { cancel, recreate } = useCrackJobsContext();
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);

  const { jobStatus: statusEnum } = enums;

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>) => {
    setMenuAnchorEl(event.currentTarget);
  };

  const handleCancel = () => {
    cancel(job.id);
  };

  const handleRecreate = () => {
    recreate(job.id);
  };

  return (
    <Box>
      <IconButton onClick={handleMenuClick} size='small'>
        <MoreVertIcon />
      </IconButton>

      <Menu open={!!menuAnchorEl} anchorEl={menuAnchorEl} onClose={() => setMenuAnchorEl(null)}>
        <MenuItem key='cancel' disabled={job.status !== statusEnum[JobStatusKey.PROCESSING]} onClick={handleCancel}>
          Cancel
        </MenuItem>
        <MenuItem key='recreate' disabled={job.status === statusEnum[JobStatusKey.PROCESSING]} onClick={handleRecreate}>
          Recreate
        </MenuItem>
      </Menu>
    </Box>
  );
}
