import styled from 'styled-components';

import { AccordionDetails } from '@mui/material';

import {
  Accordion,
  Terminal,
} from 'common/components';

export type PitbullOutputProps = {
  output: string;
}

const TerminalWindow = styled(AccordionDetails)`
  background-color: ${props => props.theme.palette.background.default};
`;

export const PitbullOutput: React.FC<PitbullOutputProps> = props => {
  const { output } = props;

  return (
    <Accordion title='Pitbull output' detailsComponent={TerminalWindow}>
      <Terminal content={output} />
    </Accordion>
  );
};

