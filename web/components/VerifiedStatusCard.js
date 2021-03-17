import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Card from '@material-ui/core/Card'
import CardActions from '@material-ui/core/CardActions'
import CardContent from '@material-ui/core/CardContent'
import Button from '@/components/Button'
import { LightTheme } from '@/components/Theme'

const useStyles = makeStyles(theme => ({
  root: {
    minWidth: 275,
  },
  link: {},
  image: {},
}))

export default function VerifiedStatusCard() {
  const classes = useStyles()

  return (
    <CardX className={classes.root}>
      <CardContent>
        <Typography className={classes.title} color="textSecondary" gutterBottom>
          Word of the Day
        </Typography>
        <Typography variant="h5" component="h2">
          Word of the Day Word of the Day Word of the Day
        </Typography>
        <Typography className={classes.pos} color="textSecondary">
          adjective
        </Typography>
        <Typography variant="body2" component="p">
          well meaning and kindly.
          <br />a benevolent smile
        </Typography>
      </CardContent>
      <CardActions>
        <Button size="small">Learn More</Button>
      </CardActions>
    </CardX>
  )
}

VerifiedStatusCard.propTypes = {}
VerifiedStatusCard.defaultProps = {}

function CardX(props) {
  return (
    <LightTheme>
      <Card {...props} />
    </LightTheme>
  )
}
