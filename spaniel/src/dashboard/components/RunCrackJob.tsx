import { useState } from 'react';

import styled from 'styled-components';

import {
  Grid,
  Card,
  CardContent,
  CardActions,
  Typography,
  Divider,
} from '@mui/material';

import { PlayArrow as PlayArrowIcon } from '@mui/icons-material';

import { RunCrackJobPayload } from 'models';

import {
  TextInput,
  CardHeader,
  Button,
} from 'common/components';

import { useCrackJobs } from '../crack-jobs.context';

export type RunCrackJobProps = {};

const CardFooter = styled(CardActions)`
  justify-content: flex-end;
`;

export const RunCrackJob: React.FC<RunCrackJobProps> = () => {
  const [payload, setPayload] = useState<RunCrackJobPayload>({ keyword: '' });

  const { run } = useCrackJobs();

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
        <Typography variant='h5'>Run Crack Job</Typography>
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
