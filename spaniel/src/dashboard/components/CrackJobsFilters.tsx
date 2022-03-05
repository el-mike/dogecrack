import { useEffect } from 'react';

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
} from 'common/components';

import { getEnumAsInputOptions } from 'core/utils';

import { useDebouncedInput } from 'core/hooks';

import { useGeneralContext } from 'core/contexts';

import { useCrackJobsContext } from '../crack-jobs.context';

  /**
   * Empty string causes rendering issues with MUI select, therefore we use
   * defined, but uncorrect value. 
   */
   const ALL_VALUE = -1;

/**
 * Checks if given status is defined, accomodating for 0 and -1 values.
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
    (event: React.ChangeEvent<HTMLInputElement>) => {
      filter({
        ...filters,
        keyword: event.target.value,
      })
    },
    300,
    [filters],
  );

  useEffect(() => {
    return () => debouncedHandleKeywordChange.cancel();
  }, [debouncedHandleKeywordChange]);

  const status = filters?.statuses?.[0];

  return (
    <Card>
      <CardHeader>
        <Typography variant='h5'>Filters</Typography>
        {!!loading && (<CircularProgress />)}
      </CardHeader>

      <Divider />

      <CardContent>
        <Grid container spacing={2}>
          <Grid item xs={6} md={4}>
            <SelectInput
              label='Job status'
              options={jobStatusOptions}
              value={isStatusValid(status) ? status : ALL_VALUE}
              onChange={handleStatusChange}
            />
          </Grid>

          <Grid item xs={6} md={4}>
            <TextInput
              label='Keyword'
              defaultValue={filters.keyword || ''}
              onChange={debouncedHandleKeywordChange}
            />
          </Grid>
        </Grid>        
      </CardContent>
    </Card>
  );
};
