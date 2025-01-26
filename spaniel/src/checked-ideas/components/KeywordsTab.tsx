import { useState } from 'react';
import styled from 'styled-components';
import {
  Box,
  Collapse,
  Divider,
  IconButton,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow
} from '@mui/material';

import { KeyboardArrowUp as KeyboardArrowUpIcon, KeyboardArrowDown as KeyboardArrowDownIcon } from '@mui/icons-material';

import { CheckedKeyword } from 'models';
import { getLabelForEnum } from 'core/utils';
import { useGeneralContext } from 'core/contexts';
import { Terminal } from 'common/components';

import { useCheckedIdeasContext } from '../checked-ideas.context';

type KeywordsTableRowProps = {
  checkedKeyword: CheckedKeyword;
};

const TokenlistsWrapper = styled.div`
    padding-top: ${props => props.theme.spacing(1)};
    padding-bottom: ${props => props.theme.spacing(1)};
`;

const TerminalWrapper = styled.div`
    padding: ${props => props.theme.spacing(1)};
    background-color: ${props => props.theme.palette.background.default};  
`;

const KeywordsTableRow: React.FC<KeywordsTableRowProps> = props => {
  const { checkedKeyword } = props;

  const { enums } = useGeneralContext();
  const [open, setOpen] = useState(false);

  const generatorVersions = checkedKeyword.generatorVersions
    .map(version => getLabelForEnum(enums.tokenGeneratorVersion, version))
    .join(",");

  return (
    <>
      <TableRow>
        <TableCell sx={{ padding: 0 }}>
          <IconButton size='small' onClick={() => setOpen(!open)}>
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </TableCell>

        <TableCell sx={{ paddingTop: 0, paddingBottom: 0 }}>
          {checkedKeyword.keyword}
        </TableCell>

        <TableCell sx={{ paddingTop: 0, paddingBottom: 0 }}>
          {checkedKeyword.runsCount}
        </TableCell>

        <TableCell sx={{ paddingTop: 0, paddingBottom: 0 }}>
          {generatorVersions}
        </TableCell>
      </TableRow>

      <TableRow>
        <TableCell sx={{ padding: 0 }} colSpan={4}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <TokenlistsWrapper>
              {(checkedKeyword?.tokenlists || []).map((tokenlist, i) => (
                <>
                  { i > 0 && <Divider />}
                  <TerminalWrapper key={i}>
                    <Terminal content={tokenlist} />
                  </TerminalWrapper>
                </>
              ))}
            </TokenlistsWrapper>
          </Collapse>
        </TableCell>
      </TableRow>
    </>
  );
};

export const KeywordsTab: React.FC = () => {
  const { checkedIdeas } = useCheckedIdeasContext();

  return (
    <Box>
      <TableContainer component={Paper}>
        <Table size='small'>
          <TableHead>
            <TableRow>
              <TableCell />
              <TableCell align='left'>Keyword</TableCell>
              <TableCell align='left'>Runs Count</TableCell>
              <TableCell align='left'>Tested Generator Versions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(checkedIdeas?.checkedKeywords || []).map((item, i) => (
              <KeywordsTableRow key={i} checkedKeyword={item} />
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
}
