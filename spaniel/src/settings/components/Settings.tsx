import React, {
  useState,
  useEffect,
  ChangeEvent,
} from 'react';

import styled from 'styled-components';

import {
  Grid,
  Card,
  CardContent,
  Typography,
  Divider,
  CircularProgress,
  CardActions,
  Box,
  Chip,
} from '@mui/material';

import { Settings as SettingsModel } from 'models';

import {
  CardHeader,
  TextInput,
  Button,
  Spacer,
} from 'common/components';

import { useSettingsContext } from '../settings.context';

const CardFooter = styled(CardActions)`
  justify-content: flex-end;
`;

const ChipsWrapper = styled(Box)`
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
`;

const ChipWrapper = styled(Box)`
  display: flex;
  padding: ${props => props.theme.spacing(0.5)};
`;

export const Settings: React.FC = () => {
  const {
    settings,
    loading,
    update,
  } = useSettingsContext();

  const [settingsForm, setSettingsForm] = useState<SettingsModel>({} as SettingsModel);
  const [presetsToAdd, setPresetsToAdd] = useState<string>('');

  useEffect(() => {
    setSettingsForm(settings);
  }, [settings]);

  const getChangeHandler = (field: keyof SettingsModel, asNumber = true) =>
    (event: ChangeEvent<HTMLInputElement>) =>
      setSettingsForm({
        ...settingsForm,
        [field]: asNumber ? event.target.valueAsNumber : event.target.value,
      });

  const handleUpdate = () => update(settingsForm);

  const onAddPreset = () => {
    const currentPresets = settingsForm.keywordPresets || [];

    if (!presetsToAdd) {
      return;
    }

    const newPresets = presetsToAdd
      .split(',')
      .map(preset => preset.trim())
      .filter(preset => !currentPresets.includes(preset));


    if (!newPresets.length) {
      return;
    }

    setSettingsForm({
      ...settingsForm,
      keywordPresets: [...currentPresets, ...newPresets],
    });

    setPresetsToAdd('');
  };

  const onDeletePreset = (preset: string) => {
    const currentPresets = settingsForm.keywordPresets || [];

    if (!preset || !currentPresets.includes(preset)) {
      return;
    }

    setSettingsForm({
      ...settingsForm,
      keywordPresets: currentPresets.filter(p => p !== preset),
    });
  };

  const presetsSorted = (settingsForm.keywordPresets || []).sort();

  return (
    <Card>
      <CardHeader>
      <Typography variant='h5'>Settings</Typography>
      {loading && (<CircularProgress />)}
      </CardHeader>

      <Divider />

      {!loading && (
        <>
        <CardContent>
          <Grid container spacing={2}>
            <Grid item xs={4} md={2}>
              <TextInput
                type='number'
                label='Start host attempts limit'
                defaultValue={settingsForm.startHostAttemptsLimit}
                onChange={getChangeHandler('startHostAttemptsLimit')}
              />
            </Grid>

            <Grid item xs={4} md={2}>
              <TextInput
                type='number'
                label='Check status retry limit'
                defaultValue={settingsForm.checkStatusRetryLimit}
                onChange={getChangeHandler('checkStatusRetryLimit')}
              />
            </Grid>

            <Grid item xs={4} md={2}>
              <TextInput
                type='number'
                label='Stalled progress limit'
                defaultValue={settingsForm.stalledProgressLimit}
                onChange={getChangeHandler('stalledProgressLimit')}
              />
            </Grid>

            <Grid item xs={4} md={2}>
              <TextInput
                type='number'
                label='Reschedule limit'
                defaultValue={settingsForm.rescheduleLimit}
                onChange={getChangeHandler('rescheduleLimit')}
              />
            </Grid>

            <Grid item xs={4} md={2}>
              <TextInput
                type='number'
                label='Running instances limit'
                defaultValue={settingsForm.runningInstancesLimit}
                onChange={getChangeHandler('runningInstancesLimit')}
              />
            </Grid>

            {/* Pushing items to the next line. */}
            <Grid item xs={12} />

            <Grid item xs={4} md={2}>
              <TextInput
                type='number'
                label='Check host interval'
                defaultValue={settingsForm.checkHostInterval}
                onChange={getChangeHandler('checkHostInterval')}
              />
            </Grid>

            <Grid item xs={4} md={2}>
              <TextInput
                type='number'
                label='Check Pitbull interval'
                defaultValue={settingsForm.checkPitbullInterval}
                onChange={getChangeHandler('checkPitbullInterval')}
              />
            </Grid>
          </Grid>

            <Spacer mb={2} />
            <Divider />
            <Spacer mb={2} />

          <Grid container spacing={2}>
            <Grid item xs={12}>
              <TextInput
                type='text'
                label='Vast search criteria'
                defaultValue={settingsForm.vastSearchCriteria}
                onChange={getChangeHandler('vastSearchCriteria', false)}
              />
            </Grid>
          </Grid>

          <Spacer mb={2} />
          <Divider />
          <Spacer mb={2} />

          <Grid container spacing={2}>
            <Grid item xs={4} md={2}>
              <TextInput
                type='number'
                label='Min password length'
                defaultValue={settingsForm.minPasswordLength}
                onChange={getChangeHandler('minPasswordLength')}
              />
            </Grid>

            <Grid item xs={4} md={2}>
              <TextInput
                type='number'
                label='Max password length'
                defaultValue={settingsForm.maxPasswordLength}
                onChange={getChangeHandler('maxPasswordLength')}
              />
            </Grid>
          </Grid>

          <Spacer mb={2} />
          <Divider />
          <Spacer mb={2} />

          <Grid container spacing={2}>
            <Grid item xs={12}>
              <Typography variant='h6'>Keyword presets</Typography>
              <Spacer mb={2} />

              <TextInput
                type='text'
                label='Add keyword presets (comma-separated)'
                value={presetsToAdd}
                onChange={(event) => setPresetsToAdd(event.target?.value)}
                onKeyPress={(event) => event.key === 'Enter' && onAddPreset()}
              />

              <Spacer mb={2} />

              <ChipsWrapper>
                {presetsSorted.map(item => (
                  <ChipWrapper>
                    <Chip color='primary' label={item} key={item} onDelete={() => onDeletePreset(item)} />
                  </ChipWrapper>
                ))}
              </ChipsWrapper>
            </Grid>
          </Grid>
        </CardContent>

        <Divider />

        <CardFooter>
          <Button
            size='large'
            variant='contained'
            onClick={handleUpdate}
          >
            Update
          </Button>
        </CardFooter>
        </>
      )}
    </Card>
  );
};
