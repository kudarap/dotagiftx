import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'

const useStyles = makeStyles()(theme => ({
  root: {
    marginTop: theme.spacing(1),
  },
  text: {
    color: theme.palette.info.light,
  },
}))

export default function MarketNotes({ text }: { text?: string }) {
  const { classes } = useStyles()

  return (
    <div className={classes.root}>
      <Typography color="textSecondary" component="span" variant="body2">
        {`Notes: `}
        <span className={classes.text}>{text}</span>
      </Typography>
    </div>
  )
}
