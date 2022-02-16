import styled from 'styled-components';

import {
  Typography,
  Accordion,
  AccordionDetails,
  AccordionSummary,
} from '@mui/material';

import { ExpandMore as ExpandMoreIcon } from '@mui/icons-material';

export type CollapsibleTerminalProps = {
  title: string;
  content: string;
}

const Output = styled.div`
  font-family: 'Roboto Mono';
`;

const Terminal = styled(AccordionDetails)`
  background-color: ${props => props.theme.palette.background.default};
`;

export const CollapsibleTerminal: React.FC<CollapsibleTerminalProps> = props => {
  const { title, content: output } = props;

  const lines = (output || '').split('\n');

  return (
    <Output>
      <Accordion>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography variant='overline'>{title}</Typography>
        </AccordionSummary>

        <Terminal>
          {(lines || []).map((line, i) => (
            <Typography
              key={i}
              fontFamily='inherit'
              fontSize='small'
            >
                {line}
              </Typography>
          ))}
        </Terminal>
      </Accordion>
    </Output>
  );
};

