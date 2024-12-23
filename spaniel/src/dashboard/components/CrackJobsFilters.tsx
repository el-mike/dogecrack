import {
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
  SelectChangeEvent,
  CircularProgress,
} from '@mui/material';

import { isNullish } from 'common/utils';

import {
  SelectInput,
  TextInput,
  InputOption,
  CardHeader,
  Button,
  Spacer,
} from 'common/components';

import { getEnumAsInputOptions } from 'core/utils';

import { useDebouncedInput } from 'core/hooks';

import { useGeneralContext } from 'core/contexts';

import { useCrackJobsContext } from '../crack-jobs.context';
import { Refresh } from '@mui/icons-material';

const ActionsContainer = styled.div`
  display: flex;
  align-items: center;
`;

/**
 * Empty string causes rendering issues with MUI select, therefore we use
 * defined, but incorrect value.
 */
  const ALL_VALUE = -1;

/**
 * Checks if given status is defined, accommodating for 0 and -1 values.
 */
const isStatusValid = (status: number | undefined) =>
  status !== -1 && !isNullish(status, { skipZero: true });

const ALL_OPTION = {
  label: 'All',
  value: ALL_VALUE,
} as InputOption;

export const CrackJobsFilters: React.FC = () => {
  const { enums } = useGeneralContext();

  const {
    filters,
    filter,
    loading,
    reload,
    resetFilters,
  } = useCrackJobsContext();

  const { jobStatus: statusEnum } = enums;

  const jobStatusOptions = [
    ...getEnumAsInputOptions(statusEnum),
    ALL_OPTION,
  ];

  const handleStatusChange = (event: SelectChangeEvent<unknown>) => {
    const status = event.target?.value as number;
    filter({
      ...filters,
      statuses: (isStatusValid(status)  && [status]) || [],
    });
  };

  const debouncedHandleKeywordChange = useDebouncedInput(
    (event: ChangeEvent<HTMLInputElement>) => {
      filter({
        ...filters,
        keyword: event.target.value,
      })
    },
    300,
    [filters],
  );


  const debouncedHandleJobIdChange = useDebouncedInput(
    (event: ChangeEvent<HTMLInputElement>) => {
      filter({
        ...filters,
        jobId: event.target.value,
      })
    },
    300,
    [filters],
  );

  const debouncedHandleNameChange = useDebouncedInput(
    (event: ChangeEvent<HTMLInputElement>) => {
      filter({
        ...filters,
        name: event.target.value,
      })
    },
    300,
    [filters],
  );

  useEffect(() => {
    return () => {
      debouncedHandleKeywordChange.cancel();
      debouncedHandleJobIdChange.cancel();
      debouncedHandleNameChange.cancel();
    };
  }, [debouncedHandleKeywordChange, debouncedHandleJobIdChange, debouncedHandleNameChange]);

  const status = filters?.statuses?.[0];

  return (
    <Card>
      <CardHeader>
        <Typography variant='h5'>Filters</Typography>
        <ActionsContainer>
          {loading && (
            <>
              <CircularProgress size={20} />
              <Spacer mr={2} />
            </>
          )}

          <Button
            size='small'
            variant='contained'
            endIcon={<Refresh />}
            onClick={reload}
            disabled={loading}
          >
            Refresh
          </Button>
        </ActionsContainer>
      </CardHeader>

      <Divider />

      <CardContent>
        <Grid container spacing={2}>
          <Grid item xs={6} md={3}>
            <SelectInput
              label='Job status'
              options={jobStatusOptions}
              value={isStatusValid(status) ? status : ALL_VALUE}
              onChange={handleStatusChange}
            />
          </Grid>

          <Grid item xs={6} md={3}>
            <TextInput
              label='Keyword'
              value={filters.keyword || ''}
              defaultValue={filters.keyword || ''}
              onChange={debouncedHandleKeywordChange}
            />
          </Grid>

          <Grid item xs={6} md={3}>
            <TextInput
              label='Job ID'
              defaultValue={filters.jobId || ''}
              onChange={debouncedHandleJobIdChange}
            />
          </Grid>
          <Grid item xs={6} md={3}>
            <TextInput
              label='Job Name'
              defaultValue={filters.name || ''}
              onChange={debouncedHandleNameChange}
            />
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  );
};
