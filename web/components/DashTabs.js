import React from 'react'
import withStyles from '@mui/styles/withStyles'
import Tabs from '@mui/material/Tabs'

const StyledTabs = withStyles(theme => ({
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
}))(props => <Tabs {...props} TabIndicatorProps={{ children: <span /> }} />)

export default StyledTabs
