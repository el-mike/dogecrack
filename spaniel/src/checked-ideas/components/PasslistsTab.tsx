import {
  Box,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow
} from '@mui/material';
import { useCheckedIdeasContext } from '../checked-ideas.context';

export const PasslistsTab: React.FC = () => {
  const { checkedIdeas } = useCheckedIdeasContext();

  return (
    <Box>
      <TableContainer component={Paper}>
        <Table size='small'>
          <TableHead>
            <TableRow>
              <TableCell />
              <TableCell align='left'>Keyword</TableCell>
              <TableCell align='left'>Tested Generator Versions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(checkedIdeas?.checkedPasslists || []).map((item, i) => (
              <TableRow>
                <TableCell sx={{ paddingTop: 0, paddingBottom: 0 }}>
                  {item.name}
                </TableCell>
                <TableCell sx={{ paddingTop: 0, paddingBottom: 0 }}>
                  {item.passlistUrl}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
};
