import React from 'react'
import { withStyles } from 'tss-react/mui'
import Tabs from '@mui/material/Tabs'

const StyledTabs = withStyles(
  props => <Tabs {...props} TabIndicatorProps={{ children: <span /> }} />,
  theme => ({
    indicator: {
      display: 'flex',
      justifyContent: 'center',
      backgroundColor: 'transparent',
      '& > span': {
        // maxWidth: 40,
        width: '100%',
        // backgroundColor: theme.palette.accent.main,
        backgroundColor: theme.palette.grey[400],
      },
    },
  })
)

export default StyledTabs
