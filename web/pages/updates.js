import React from 'react'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Divider from '@material-ui/core/Divider'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

export default function Updates() {
  const classes = useStyles()

  return (
    <div className="container">
      <Head>
        <title>{APP_NAME} :: Updates</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Survey and Feedback
            <Typography color="textSecondary">March 2, 2021 Update</Typography>
          </Typography>
          <Typography color="textSecondary">
            - Rework account pages
            <br />- Rework account pages
          </Typography>
          <br />
          <Divider />
          <br />
          <br />

          <Typography variant="h5" component="h1" gutterBottom>
            Febuary 26, 2021 Update
          </Typography>
          <Typography color="textSecondary">
            - Fixed the interaction of Earthshaker's Aghanim's Shard and a few other hero spells
            that blocked pathing (notably Nature's Prophet's Sprout). Earthshaker cannot walk
            through other blockers when he has a Shard, he can only walk through Fissure blockers.
          </Typography>
          <br />
          <Divider />
          <br />
          <br />

          <Typography variant="h5" component="h1" gutterBottom>
            Febuary 3, 2021 Update
          </Typography>
          <Typography color="textSecondary">
            - Fixed Psychic Headband destroying Homing Missile
            <br /> - Fixed Flaming Lasso doing more damage than intended
            <br />- Fixed Level 25 Nature's Prophet talent not giving enough health
            <br />- Fixed Stone Remnant duration being higher than intended
          </Typography>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
