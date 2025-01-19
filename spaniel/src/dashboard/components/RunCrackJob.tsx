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
} from '@mui/material';

import { PlayArrow as PlayArrowIcon } from '@mui/icons-material';

import { RunCrackJobPayload } from 'models';

import {
  TextInput,
  CardHeader,
  Button,
  Spacer,
  SelectInput,
} from 'common/components';
import { useGeneralContext } from 'core/contexts';
import { getEnumAsInputOptions } from 'core/utils';

import { useCrackJobsContext } from '../crack-jobs.context';

export type RunCrackJobProps = {};

const CardFooter = styled(CardActions)`
  justify-content: flex-end;
`;

export const RunCrackJob: React.FC<RunCrackJobProps> = () => {
  const { run } = useCrackJobsContext();
  const { enums, latestTokenGeneratorVersion } = useGeneralContext();

  const [payload, setPayload] = useState<RunCrackJobPayload>({
    keywords: [],
    passlistUrl: '',
    name: '',
    tokenlist: '',
    tokenGeneratorVersion: latestTokenGeneratorVersion,
  });

  const tokenGeneratorVersionOptions = [
    ...getEnumAsInputOptions(enums.tokenGeneratorVersion),
  ];

  const handleKeywordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target?.value;

    setPayload((prev) => ({
      ...prev,
      keywords: value ? value.split('\n') : [],
    }));
  };


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

  const handleTokenlistChange = (event: React.ChangeEvent<HTMLInputElement>) =>
    setPayload((prev) => ({
      ...prev,
      tokenlist: event.target?.value || '',
    }));

  const handleTokenGeneratorVersionChange = (event: SelectChangeEvent<unknown>) =>
    setPayload((prev) => ({
      ...prev,
      tokenGeneratorVersion: event.target?.value as number || latestTokenGeneratorVersion,
    }));

  const handleRun = () => {
    run(payload);
    setPayload({
      keywords: [],
      passlistUrl: '',
      name: '',
      tokenlist:'',
      tokenGeneratorVersion: latestTokenGeneratorVersion,
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
              label='Name'
              value={payload.name}
              onChange={handleNameChange}
            />
          </Grid>
        </Grid>

        <Spacer mb={2} />
        <Divider />
        <Spacer mb={2} />

        <Grid container spacing={2}>
          <Grid item xs={12} sm={6} md={4}>
            <SelectInput
              label='Token generator version'
              options={tokenGeneratorVersionOptions}
              value={payload.tokenGeneratorVersion}
              onChange={handleTokenGeneratorVersionChange}
            />
          </Grid>
          <Grid item xs={12}>
            <TextInput
              label='Keywords (newline-separated)'
              multiline
              value={payload.keywords?.join('\n')}
              onChange={handleKeywordChange}
              disabled={!!payload.passlistUrl || !!payload.tokenlist}
            />
          </Grid>
        </Grid>

        <Spacer mb={2} />
        <Divider />
        <Spacer mb={2} />

        <Grid container spacing={2}>
          <Grid item xs={12}>
            <TextInput
              label='Passlist URL'
              value={payload.passlistUrl}
              onChange={handlePasslistUrlChange}
              disabled={!!payload.keywords?.length || !!payload.tokenlist}
            />
          </Grid>
        </Grid>

        <Spacer mb={2} />
        <Divider />
        <Spacer mb={2} />

        <Grid container spacing={2}>
          <Grid item xs={12}>
            <TextInput
              label='Tokenlist'
              value={payload.tokenlist}
              onChange={handleTokenlistChange}
              multiline
              disabled={!!payload.keywords?.length || !!payload.passlistUrl}
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
