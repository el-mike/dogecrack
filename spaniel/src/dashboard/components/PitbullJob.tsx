import styled from 'styled-components';

import { PitbullJob as PitbullJobModel } from "models";

import {
  Card,
  CardContent,
  Typography,
} from '@mui/material';

import { Spacer } from 'common/components';

export type PitbullJobsProps = {
  job: PitbullJobModel;
};

const CardHeader = styled.div`
  display: flex;
  justify-content: space-between;
  padding: ${props => props.theme.spacing(2)};
`;

export const PitbullJob: React.FC<PitbullJobsProps> = props => {
  const { job } = props;

  return (
    <Card>
      <CardHeader>
        <Typography variant='subtitle1' fontWeight=''>&nbsp; {job.id}</Typography>
      </CardHeader>
      
      <CardContent>
        {job.passlistUrl}
      </CardContent>
    </Card>
  );
};
