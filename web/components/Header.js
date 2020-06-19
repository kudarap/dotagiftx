import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import Typography from '@material-ui/core/Typography'
import Container from '@/components/Container'

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    textShadow: '0px 0px 16px #C79123',
    fontWeight: 'bold',
    background: 'linear-gradient(#F8E8B9 10%, #fff 90%)',
    '-webkit-background-clip': 'text',
    '-webkit-text-fill-color': 'transparent',
    filter: 'drop-shadow(0px 0px 10px black)',
    letterSpacing: 2,
  },
}))

export default function () {
  const classes = useStyles()

  return (
    <header>
      <AppBar position="static" variant="outlined">
        <Container disableMinHeight>
          <Toolbar variant="dense" disableGutters>
            <Typography variant="h6" className={classes.title}>
              Dota 2 Giftables
            </Typography>
          </Toolbar>
        </Container>
      </AppBar>
    </header>
  )
}
