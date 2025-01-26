import React, { useState } from 'react';
import styled from 'styled-components';

import {
  Box,
  CircularProgress,
  Tab,
  Tabs
} from '@mui/material';
import { Spacer } from 'common/components';
import { KeywordsTab } from './KeywordsTab';
import { PasslistsTab } from './PasslistsTab';
import { useCheckedIdeasContext } from '../checked-ideas.context';

const TabsWrapper = styled(Box)``;

const TabItem = styled(Tab)`
  font-weight: bold;
`;

export const CheckedIdeas: React.FC = () => {
  const { loading } = useCheckedIdeasContext();
  const [tab, setTab] = useState(0);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTab(newValue);
  };


  return (
    <Box>
      <TabsWrapper>
        <Tabs value={tab} onChange={handleTabChange}>
          <TabItem label="Keywords"></TabItem>
          <TabItem label="Passlists"></TabItem>
        </Tabs>
      </TabsWrapper>

      <Spacer mb={4} />

      <Box>
        {loading && (<CircularProgress />)}
        {!loading && (
          <>
            {tab === 0 && <KeywordsTab/>}
            {tab === 1 && <PasslistsTab/>}
          </>
        )}
      </Box>
    </Box>
  );
}
