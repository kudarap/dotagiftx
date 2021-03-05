import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import Tabs from '@material-ui/core/Tabs'

const StyledTabs = withStyles(theme => ({
  indicator: {
    display: 'flex',
    justifyContent: 'center',
    backgroundColor: 'transparent',
    '& > span': {
      // maxWidth: 40,
      width: '100%',
      backgroundColor: theme.palette.accent.main,
    },
  },
}))(props => <Tabs {...props} TabIndicatorProps={{ children: <span /> }} />)

export default StyledTabs
