import styled from 'styled-components';

import { Box } from '@mui/material';

export const CardHeader = styled(Box)`
display: flex;
justify-content: space-between;
align-items: center;
padding: ${props => props.theme.spacing(2)};
`;
