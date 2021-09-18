import React from 'react'
import PropTypes from 'prop-types'
import querystring from 'querystring'
import { useTheme } from '@mui/material/styles'
import makeStyles from '@mui/styles/makeStyles'
import IconButton from '@mui/material/IconButton'
import Typography from '@mui/material/Typography'
import FirstPageIcon from '@mui/icons-material/FirstPage'
import KeyboardArrowLeft from '@mui/icons-material/KeyboardArrowLeft'
import KeyboardArrowRight from '@mui/icons-material/KeyboardArrowRight'
import LastPageIcon from '@mui/icons-material/LastPage'
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
  onPageChange,
  linkProps,
  ...other
}) {
  const classes = useStyles()
  const theme = useTheme()

  const page = Number(initPage)

  const handleFirstPageButtonClick = evt => {
    onPageChange(evt, 1)
  }

  const handleBackButtonClick = evt => {
    onPageChange(evt, page - 1)
  }

  const handleNextButtonClick = evt => {
    onPageChange(evt, page + 1)
  }

  const handleLastPageButtonClick = evt => {
    onPageChange(evt, Math.max(0, Math.ceil(count / rowsPerPage)))
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
        aria-label="First Page"
        size="large">
        {theme.direction === 'rtl' ? <LastPageIcon /> : <FirstPageIcon />}
      </IconButton>
      <IconButton
        component={Link}
        {...getLinkProps(page - 1)}
        onClick={handleBackButtonClick}
        disabled={page === 1}
        aria-label="Previous Page"
        size="large">
        {theme.direction === 'rtl' ? <KeyboardArrowRight /> : <KeyboardArrowLeft />}
      </IconButton>
      <IconButton
        component={Link}
        {...getLinkProps(page + 1)}
        onClick={handleNextButtonClick}
        disabled={page >= Math.ceil(count / rowsPerPage)}
        aria-label="Next Page"
        size="large">
        {theme.direction === 'rtl' ? <KeyboardArrowLeft /> : <KeyboardArrowRight />}
      </IconButton>
      <IconButton
        component={Link}
        {...getLinkProps(Math.ceil(count / rowsPerPage))}
        onClick={handleLastPageButtonClick}
        disabled={page >= Math.ceil(count / rowsPerPage)}
        aria-label="Last Page"
        size="large">
        {theme.direction === 'rtl' ? <FirstPageIcon /> : <LastPageIcon />}
      </IconButton>
    </div>
  )
}
TablePaginationRouter.propTypes = {
  linkProps: PropTypes.object.isRequired,
  count: PropTypes.number.isRequired,
  onPageChange: PropTypes.func,
  page: PropTypes.number.isRequired,
  rowsPerPage: PropTypes.number,
}
TablePaginationRouter.defaultProps = {
  rowsPerPage: 10,
  onPageChange: () => {},
}

export default TablePaginationRouter
