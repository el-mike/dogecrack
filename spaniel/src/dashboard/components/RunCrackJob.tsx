import { useState } from 'react';

import styled from 'styled-components';

import {
  Box,
  Grid,
  Card,
  CardContent,
  CardActions,
  Typography,
  Divider,
  SelectChangeEvent,
  IconButton,
  Menu,
  MenuItem,
  ListItemIcon,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
} from '@mui/material';

import {
  PlayArrow as PlayArrowIcon,
  MoreVert as MoreVertIcon,
  Rule as RuleIcon,
  RemoveDone as RemoveDoneIcon,
} from '@mui/icons-material';

import {
  GetKeywordSuggestionsPayload,
  RunCrackJobPayload
} from 'models';

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

type RunCrackJobForm = RunCrackJobPayload & {
  tokenlistBaseKeyword: string;
};

const CardFooter = styled(CardActions)`
  justify-content: flex-end;
`;

const KeywordsActionsWrapper = styled(Box)`
  display: flex;
  justify-content: start;
  align-items: center;
  height: 100%;
`;

export const RunCrackJob: React.FC<RunCrackJobProps> = () => {
  const { run, getKeywordSuggestions, getUsedKeywords } = useCrackJobsContext();
  const { enums, latestTokenGeneratorVersion } = useGeneralContext();

  const [conflictingKeywords, setConflictingKeywords] = useState<string[]>([]);
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);

  const [payload, setPayload] = useState<RunCrackJobForm>({
    keywords: [],
    passlistUrl: '',
    name: '',
    tokenlist: '',
    tokenlistBaseKeyword: '',
    tokenGeneratorVersion: latestTokenGeneratorVersion,
  });

  const tokenGeneratorVersionOptions = [
    ...getEnumAsInputOptions(enums.tokenGeneratorVersion),
  ];

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>) => {
    setMenuAnchorEl(event.currentTarget);
  };

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

  const handleTokenlistBaseKeywordChange = (event: React.ChangeEvent<HTMLInputElement>) =>
    setPayload((prev) => ({
      ...prev,
      tokenlistBaseKeyword: event.target?.value || '',
    }));

  const handleTokenGeneratorVersionChange = (event: SelectChangeEvent<unknown>) =>
    setPayload((prev) => ({
      ...prev,
      tokenGeneratorVersion: event.target?.value as number || latestTokenGeneratorVersion,
    }));

  const handleRun = () => {
    /* When using custom tokenlist, set base keyword if available. */
    if (payload.tokenlist && payload.tokenlistBaseKeyword) {
      payload.keywords = [payload.tokenlistBaseKeyword];
    }

    run(payload);

    setPayload({
      keywords: [],
      passlistUrl: '',
      name: '',
      tokenlist: '',
      tokenlistBaseKeyword: '',
      tokenGeneratorVersion: latestTokenGeneratorVersion,
    });
  };

  const validateAndRun = async () => {
    if (!payload.keywords?.length || !payload.tokenGeneratorVersion) {
      handleRun();
    }

    try {
      const usedKeywords = await getUsedKeywords({ tokenGeneratorVersion: payload.tokenGeneratorVersion || 0});

      const conflictingKeywords = (payload.keywords || [])?.filter(keyword => usedKeywords.includes(keyword));

      if (!conflictingKeywords?.length) {
        handleRun();
        return;
      }

      setConflictingKeywords(conflictingKeywords);
    } catch {
      return;
    }
  };

  const handleGetKeywordSuggestions = async (payload: GetKeywordSuggestionsPayload) => {
    const suggestions = await getKeywordSuggestions(payload);

    setPayload((prev) => ({
      ...prev,
      keywords: suggestions,
    }));

    setMenuAnchorEl(null);
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
        </Grid>

        <Spacer mb={2} />

        <Grid container spacing={2}>
          <Grid item xs={6}>
            <TextInput
              label='Keywords (newline-separated)'
              multiline
              maxRows={10}
              value={payload.keywords?.join('\n')}
              onChange={handleKeywordChange}
              disabled={!!payload.passlistUrl || !!payload.tokenlist}
            />
          </Grid>

          <Grid item xs={6}>
            <KeywordsActionsWrapper>
              <IconButton onClick={handleMenuClick}>
                <MoreVertIcon />
              </IconButton>

              <Menu open={!!menuAnchorEl} anchorEl={menuAnchorEl} onClose={() => setMenuAnchorEl(null)}>
                <MenuItem key='cancel' onClick={() => handleGetKeywordSuggestions({ presetsOnly: true, tokenGeneratorVersion: payload.tokenGeneratorVersion })}>
                  <ListItemIcon>
                    <RemoveDoneIcon />
                  </ListItemIcon>

                  <Typography variant='body1'>
                    Get unchecked preset keywords
                  </Typography>
                </MenuItem>

                <MenuItem key='cancel' onClick={() => handleGetKeywordSuggestions({ presetsOnly: false, tokenGeneratorVersion: payload.tokenGeneratorVersion })}>
                  <ListItemIcon>
                    <RemoveDoneIcon />
                  </ListItemIcon>

                  <Typography variant='body1'>
                    Get unchecked known keywords
                  </Typography>
                </MenuItem>

                <MenuItem key='cancel' onClick={() => handleGetKeywordSuggestions({ presetsOnly: true })}>
                  <ListItemIcon>
                    <RuleIcon />
                  </ListItemIcon>

                  <Typography variant='body1'>
                    Get all preset keywords
                  </Typography>
                </MenuItem>

                <MenuItem key='cancel' onClick={() => handleGetKeywordSuggestions({ presetsOnly: false })}>
                  <ListItemIcon>
                    <RuleIcon />
                  </ListItemIcon>

                  <Typography variant='body1'>
                    Get all known keywords
                  </Typography>
                </MenuItem>
              </Menu>
            </KeywordsActionsWrapper>
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
          <Grid item xs={12} sm={6} md={4}>
            <TextInput
              label='Tokenlist base keyword'
              value={payload.tokenlistBaseKeyword}
              onChange={handleTokenlistBaseKeywordChange}
              multiline
              disabled={!!payload.keywords?.length || !!payload.passlistUrl}
            />
          </Grid>
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
          onClick={validateAndRun}
        >
          Run
        </Button>
      </CardFooter>

      <Dialog
        open={!!conflictingKeywords?.length}
        onClose={() => setConflictingKeywords([])}
      >
        <DialogTitle>Are you sure?</DialogTitle>

        <DialogContent>
            <Typography variant='body1'>
              Following keywords have already been used:
            </Typography>
          <Typography variant='body1' fontWeight='bold'>
            {(conflictingKeywords || [])?.join(', ')}
          </Typography>
          <Typography variant='body1'>
            Do you want to continue?
          </Typography>
        </DialogContent>

        <DialogActions>
          <Button onClick={() => setConflictingKeywords([])}>Cancel</Button>
          <Button variant='contained' onClick={() => {
            setConflictingKeywords([]);
            handleRun();
          }}>Confirm</Button>
        </DialogActions>
      </Dialog>
    </Card>
  );
};
