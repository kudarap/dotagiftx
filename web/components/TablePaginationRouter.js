import React from 'react'
import PropTypes from 'prop-types'
import querystring from 'querystring'
import { makeStyles, useTheme } from '@material-ui/core/styles'
import IconButton from '@material-ui/core/IconButton'
import Typography from '@material-ui/core/Typography'
import FirstPageIcon from '@material-ui/icons/FirstPage'
import KeyboardArrowLeft from '@material-ui/icons/KeyboardArrowLeft'
import KeyboardArrowRight from '@material-ui/icons/KeyboardArrowRight'
import LastPageIcon from '@material-ui/icons/LastPage'
import Link from '@/components/Link'

const useStyles = makeStyles(theme => ({
  caption: {
    marginRight: theme.spacing(2.5),
  },
}))

function TablePaginationRouter({
  count,
  page: initPage,
  rowsPerPage,
  onChangePage,
  linkProps,
  ...other
}) {
  const classes = useStyles()
  const theme = useTheme()

  const page = Number(initPage)

  const handleFirstPageButtonClick = evt => {
    onChangePage(evt, 1)
  }

  const handleBackButtonClick = evt => {
    onChangePage(evt, page - 1)
  }

  const handleNextButtonClick = evt => {
    onChangePage(evt, page + 1)
  }

  const handleLastPageButtonClick = evt => {
    onChangePage(evt, Math.max(0, Math.ceil(count / rowsPerPage)))
  }

  const cPage = page === 0 ? 1 : page
  const resultMinCount = cPage * rowsPerPage - rowsPerPage + 1
  const resultMaxCount = cPage * rowsPerPage

  const getLinkProps = p => {
    const q = { ...linkProps.query, page: p }
    if (!linkProps.as) {
      return { href: `${linkProps.href}?${querystring.stringify(q)}` }
    }

    return { href: linkProps.href, as: `${linkProps.as}?${querystring.stringify(q)}` }
  }

  if (count === 0) {
    return null
  }

  return (
    <div {...other}>
      <Typography
        className={classes.caption}
        component="span"
        variant="body2"
        color="textSecondary">
        {resultMinCount}-{resultMaxCount >= count ? count : resultMaxCount} of {count}
      </Typography>
      <IconButton
        component={Link}
        {...getLinkProps(1)}
        onClick={handleFirstPageButtonClick}
        disabled={page === 1}
        aria-label="First Page">
        {theme.direction === 'rtl' ? <LastPageIcon /> : <FirstPageIcon />}
      </IconButton>
      <IconButton
        component={Link}
        {...getLinkProps(page - 1)}
        onClick={handleBackButtonClick}
        disabled={page === 1}
        aria-label="Previous Page">
        {theme.direction === 'rtl' ? <KeyboardArrowRight /> : <KeyboardArrowLeft />}
      </IconButton>
      <IconButton
        component={Link}
        {...getLinkProps(page + 1)}
        onClick={handleNextButtonClick}
        disabled={page >= Math.ceil(count / rowsPerPage)}
        aria-label="Next Page">
        {theme.direction === 'rtl' ? <KeyboardArrowLeft /> : <KeyboardArrowRight />}
      </IconButton>
      <IconButton
        component={Link}
        {...getLinkProps(Math.ceil(count / rowsPerPage))}
        onClick={handleLastPageButtonClick}
        disabled={page >= Math.ceil(count / rowsPerPage)}
        aria-label="Last Page">
        {theme.direction === 'rtl' ? <FirstPageIcon /> : <LastPageIcon />}
      </IconButton>
    </div>
  )
}
TablePaginationRouter.propTypes = {
  linkProps: PropTypes.object.isRequired,
  count: PropTypes.number.isRequired,
  onChangePage: PropTypes.func,
  page: PropTypes.number.isRequired,
  rowsPerPage: PropTypes.number,
}
TablePaginationRouter.defaultProps = {
  rowsPerPage: 10,
  onChangePage: () => {},
}

export default TablePaginationRouter
