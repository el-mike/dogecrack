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

  const content = (job.tokens || []).join('\n');

  return (
    <Accordion title='Tokens' detailsComponent={TerminalWindow}>
      <Terminal content={content} />
    </Accordion>
  );
};
