import styled from 'styled-components';

import {
  Typography,
  Box,
} from '@mui/material';

export type CollapsibleTerminalProps = {
  content: string;
}

const TerminalWrapper = styled(Box)`
  font-family: 'Roboto Mono';
`;

export const Terminal: React.FC<CollapsibleTerminalProps> = props => {
  const { content } = props;

  const lines = (content || '').split('\n');

  return (
    <TerminalWrapper>
      {(lines || []).map((line, i) => (
        <Typography
          key={i}
          fontFamily='inherit'
          fontSize='small'
        >
            {line}
          </Typography>
      ))}
    </TerminalWrapper>
  );
};

