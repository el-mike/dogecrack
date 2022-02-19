import styled from 'styled-components';

import { AccordionDetails } from '@mui/material';

import {
  Accordion,
  Terminal,
} from 'common/components';

export type PitbullJobErrorLogProps = {
  errorLog: string;
}

const TerminalWindow = styled(AccordionDetails)`
  background-color: ${props => props.theme.palette.background.default};
`;

export const PitbullJobErrorLog: React.FC<PitbullJobErrorLogProps> = props => {
  const { errorLog } = props;

  return (
    <Accordion title='Job error log' detailsComponent={TerminalWindow}>
      <Terminal content={errorLog} />
    </Accordion>
  );
};

