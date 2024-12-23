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
  Spacer,
} from 'common/components';

import { useCrackJobsContext } from '../crack-jobs.context';

export type RunCrackJobProps = {};

const CardFooter = styled(CardActions)`
  justify-content: flex-end;
`;

export const RunCrackJob: React.FC<RunCrackJobProps> = () => {
  const [payload, setPayload] = useState<RunCrackJobPayload>({ keyword: '', passlistUrl: '', name: '' });

  const { run } = useCrackJobsContext();

  const handleKeywordChange = (event: React.ChangeEvent<HTMLInputElement>) =>
    setPayload((prev) => ({
      ...prev,
      keyword: event.target?.value || '',
    }));

  const handlePasslistUrlChange = (event: React.ChangeEvent<HTMLInputElement>) =>
    setPayload((prev) => ({
      ...prev,
      passlistUrl: event.target?.value || '',
    }));

  const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) =>
    setPayload((prev) => ({
      ...prev,
      name: event.target?.value || '',
    }));

  const handleRun = () => {
    run(payload);
    setPayload({
      keyword: '',
      passlistUrl: '',
      name: '',
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
              onChange={handleKeywordChange}
              disabled={!!payload.passlistUrl}
            />
          </Grid>

            <Grid item xs={12} sm={6} md={4}>
              <TextInput
                label='Passlist URL'
                value={payload.passlistUrl}
                onChange={handlePasslistUrlChange}
                disabled={!!payload.keyword}
              />
          </Grid>
        </Grid>

        <Spacer mb={2} />
        <Divider />
        <Spacer mb={2} />

        <Grid container spacing={2}>
          <Grid item xs={12} sm={6} md={4}>
            <TextInput
              label='Name'
              value={payload.name}
              onChange={handleNameChange}
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
