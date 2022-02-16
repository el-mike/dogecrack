import { Grid } from '@mui/material';

import { PitbullInstance } from 'models';

import {
  Typography,
  Box,
} from '@mui/material';

import {
  Spacer,
  CollapsibleTerminal,
  LabeledInfo
} from 'common/components';

import {
  getInstanceEstimatedCost,
  toDateTimeString,
} from 'core/utils';

import {
  TimeFor,
  TimeAgo,
} from 'core/components';

import { PitbullProgress } from './PitbullProgress';
import { PitbullStatus } from './PitbullStatus';

export type PitbullInfoProps = {
  instance: PitbullInstance;
};

export const PitbullInfo: React.FC<PitbullInfoProps> = props => {
  const { instance } = props;
  
  const { hostInstance } = instance;

  return (
    <Box>
      <Typography variant='overline'>Pitbull info</Typography>

      <Spacer mb={2} />

      <Grid container spacing={2}>
        <Grid container spacing={2} item xs={12}>
          <Grid container spacing={2} item xs={12} md={6} lg={4}>
            <Grid item xs={12}>
              <PitbullStatus status={instance?.status || -1} />
            </Grid>
          
            <Spacer mb={2} />
          
            <Grid item xs={12}>
              <PitbullProgress progress={instance?.progress || {}} />
            </Grid>
          </Grid>

          <Grid container spacing={2} item xs={12} md={6} lg={8}>
            <Grid item xs={6} md={4}>
              <LabeledInfo title='Last updated:'>
                <TimeAgo from={instance?.updatedAt} />
              </LabeledInfo>
            </Grid>
      
            <Grid item xs={6} md={4}>
              <LabeledInfo
                title='Started at:'
                value={instance.startedAt && toDateTimeString(new Date(instance.startedAt))}
              />
            </Grid>

            <Grid item xs={6} md={4}>
              <LabeledInfo
                title='Completed at:'
                value={instance.completedAt && toDateTimeString(new Date(instance.completedAt))}
              />
            </Grid>

            <Grid item xs={6} md={4}>
              <LabeledInfo title='Run for:'>
                {
                  !!instance.startedAt
                    ? (<TimeFor from={instance.startedAt} to={instance.completedAt} />)
                    : '-'
                }
              </LabeledInfo>
            </Grid>

            <Grid item xs={6} md={4}>
              <LabeledInfo
                title='Estimated cost:'
                value={`${getInstanceEstimatedCost(instance)} $`}
              />
          </Grid>
          </Grid>
        </Grid>

        <Grid container item xs={12}>
          <Grid item xs={6} md={3}>
            <LabeledInfo
                title='Provider:'
                value={instance?.providerName}
              />
          </Grid>

          <Grid item xs={6} md={4}>
              <LabeledInfo
                title='Host address:'
                value={
                  hostInstance?.sshHost
                    ? `${hostInstance?.sshHost}:${hostInstance.sshPort}`
                    : '-'
                }
                toCopy={
                  hostInstance?.sshHost
                    ? `ssh -p ${hostInstance?.sshPort} root@${hostInstance?.sshHost}`
                    : ''
                }
              />
            </Grid>
        </Grid>
    

        <Grid container item xs={12}>
          <Grid item xs={6} md={3}>
            <LabeledInfo
              title='GPU:'
              value={`${hostInstance?.gpuNum || 0} x ${hostInstance?.gpuName}`}
            />
          </Grid>

          <Grid item xs={6} md={3}>
            <LabeledInfo
              title='USD/Hour:'
              value={`${hostInstance?.dphTotal?.toFixed(3)} $`}
            />
          </Grid>

          <Grid item xs={6} md={3}>
            <LabeledInfo
              title='DLPerf:'
              value={`${hostInstance?.dlperf?.toFixed(3)}`}
            />
          </Grid>

          <Grid item xs={6} md={3}>
            <LabeledInfo
              title='DLPerf per cost:'
              value={`${hostInstance?.dlperfPerDphTotal?.toFixed(3)}`}
            />
          </Grid>
        </Grid>
    
        <Grid container item xs={12}>
          <Grid item xs={12}>
            <CollapsibleTerminal
              title='Pitbull output'
              content={instance?.lastOutput}
            />
          </Grid>
        </Grid>
      </Grid>
    </Box>
  );
};
