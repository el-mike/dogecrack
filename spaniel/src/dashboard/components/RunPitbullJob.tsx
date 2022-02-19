import { useState } from 'react';

import styled from 'styled-components';

import {
  Grid,
  Card,
  CardContent,
  CardActions,
  Typography,
  Divider,
  SelectChangeEvent,
  CircularProgress,
} from '@mui/material';

import { PlayArrow as PlayArrowIcon } from '@mui/icons-material';

import { RunPitbullJobPayload } from 'models';

import {
  SelectInput,
  TextInput,
  InputOption,
  CardHeader,
  Button,
  Spacer,
} from 'common/components';

import { usePitbullJobs } from '../pitbull-jobs.context';

export type RunPitbulLJobProps = {};

const CardFooter = styled(CardActions)`
  justify-content: flex-end;
`;

export const RunPitbullJob: React.FC<RunPitbulLJobProps> = () => {
  const [payload, setPayload] = useState<RunPitbullJobPayload>({ keyword: '' });

  const { run } = usePitbullJobs();

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) =>
    setPayload({
      keyword: event.target?.value || '',
    });

  const handleRun = () => {
    run(payload);
    setPayload({
      keyword: '',
    });
  };

  return (
    <Card>
      <CardHeader>
        <Typography variant='h5'>Run Pitbull</Typography>
      </CardHeader>
      
      <Divider />

      <CardContent>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={6} md={4}>
            <TextInput
              label='Keyword'
              value={payload.keyword}
              onChange={handleChange}
            />
          </Grid>
        </Grid>
      </CardContent>

      <Divider />

      <CardFooter>
      <Button
        size='large'
        variant='contained'
        endIcon={<PlayArrowIcon />}
        onClick={handleRun}
      >
        Run
      </Button>
      </CardFooter>
    </Card>
  );
};
