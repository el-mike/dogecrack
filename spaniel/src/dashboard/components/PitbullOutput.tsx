import styled from 'styled-components';

import {
  Typography,
  Accordion,
  AccordionDetails,
  AccordionSummary,
} from '@mui/material';

import { ExpandMore as ExpandMoreIcon } from '@mui/icons-material';

export type PitbullOutputProps = {
  output: string;
}

const Output = styled.div`
  font-family: 'Roboto Mono';
`;

const Terminal = styled(AccordionDetails)`
  background-color: ${props => props.theme.palette.background.default};
`;

export const PitbullOutput: React.FC<PitbullOutputProps> = props => {
  const { output } = props;

  const lines = (output || '').split('\n');

  return (
    <Output>
      <Accordion>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography variant='overline'>Pitbull output</Typography>
        </AccordionSummary>

        <Terminal>
          {(lines || []).map(line => (
            <Typography fontFamily='inherit' fontSize='small'>{line}</Typography>
          ))}
        </Terminal>
      </Accordion>
    </Output>
  );
};

