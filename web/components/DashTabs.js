import React from 'react'
import Tabs from '@mui/material/Tabs'
import { styled } from '@mui/material/styles'

const StyledTabs = styled(props => <Tabs {...props} TabIndicatorProps={{ children: <span /> }} />)(
  ({ theme }) => ({
    indicator: {
      display: 'flex',
      justifyContent: 'center',
      backgroundColor: 'transparent',
      '& > span': {
        width: '100%',
        backgroundColor: theme.palette.grey[400],
      },
    },
  })
)

export default StyledTabs
