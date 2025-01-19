import styled from 'styled-components';
import { AccordionDetails } from '@mui/material';

import { CrackJob } from 'models';

import { Accordion, Terminal } from 'common/components';

export type CrackJobTokensProps = {
  job: CrackJob;
};

const TerminalWindow = styled(AccordionDetails)`
  background-color: ${props => props.theme.palette.background.default};
`;

export const CrackJobTokens: React.FC<CrackJobTokensProps> = props => {
  const { job } = props;

  return (
    <Accordion title='Tokenlist' detailsComponent={TerminalWindow}>
      <Terminal content={job.tokenlist || ''} />
    </Accordion>
  );
};
