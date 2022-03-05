import { Grid } from '@mui/material';

import { PitbullInstance } from 'models';

import {
  LabeledInfo,
  Accordion,
} from 'common/components';

export type PitbullHostInfoProps = {
  instance: PitbullInstance;
};

export const PitbullHostInfo: React.FC<PitbullHostInfoProps> = props => {
  const { instance } = props;
  
  const { hostInstance } = instance;

  return (
    <Accordion title='Host machine'>
      <Grid container spacing={2}>
        <Grid container item xs={12}>
          <Grid item xs={6} md={3}>
            <LabeledInfo
                title='Provider:'
                value={instance?.providerName}
              />
          </Grid>

          <Grid item xs={6} md={3}>
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

          <Grid item xs={6} md={3}>
            <LabeledInfo
                title='Host ID:'
                value={hostInstance?.id}
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
              value={`${hostInstance?.dphTotal?.toFixed(2)} $`}
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
      </Grid>
    </Accordion>
  );
};
