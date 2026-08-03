import React from 'react'
import querystring from 'querystring'
import { useTheme } from '@mui/material/styles'
import { makeStyles } from 'tss-react/mui'
import IconButton from '@mui/material/IconButton'
import Typography from '@mui/material/Typography'
import FirstPageIcon from '@mui/icons-material/FirstPage'
import KeyboardArrowLeft from '@mui/icons-material/KeyboardArrowLeft'
import KeyboardArrowRight from '@mui/icons-material/KeyboardArrowRight'
import LastPageIcon from '@mui/icons-material/LastPage'
import Link from '@/components/Link'

const useStyles = makeStyles()(theme => ({
  caption: {
    marginRight: theme.spacing(2.5),
  },
}))

interface LinkProps {
  href: string
  as?: string
  query?: Record<string, string | number | undefined>
}

interface TablePaginationRouterProps {
  count: number
  page: number
  rowsPerPage?: number
  colSpan?: number
  style?: React.CSSProperties
  onPageChange?: (event: unknown, page: number) => void
  linkProps: LinkProps
}

function TablePaginationRouter({
  count,
  page: initPage,
  rowsPerPage = 10,
  onPageChange = () => {},
  linkProps,
  ...other
}: TablePaginationRouterProps) {
  const { classes } = useStyles()
  const theme = useTheme()

  const page = Number(initPage)

  const handleFirstPageButtonClick = (evt: unknown) => {
    onPageChange(evt, 1)
  }

  const handleBackButtonClick = (evt: unknown) => {
    onPageChange(evt, page - 1)
  }

  const handleNextButtonClick = (evt: unknown) => {
    onPageChange(evt, page + 1)
  }

  const handleLastPageButtonClick = (evt: unknown) => {
    onPageChange(evt, Math.max(0, Math.ceil(count / rowsPerPage)))
  }

  const cPage = page === 0 ? 1 : page
  const resultMinCount = cPage * rowsPerPage - rowsPerPage + 1
  const resultMaxCount = cPage * rowsPerPage

  const getLinkProps = (p: number) => {
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
        color="textSecondary"
      >
        {resultMinCount}-{resultMaxCount >= count ? count : resultMaxCount} of {count}
      </Typography>
      <IconButton
        component={Link}
        nativeButton={false}
        {...getLinkProps(1)}
        onClick={handleFirstPageButtonClick}
        disabled={page === 1}
        aria-label="First Page"
        size="large"
      >
        {theme.direction === 'rtl' ? <LastPageIcon /> : <FirstPageIcon />}
      </IconButton>
      <IconButton
        component={Link}
        nativeButton={false}
        {...getLinkProps(page - 1)}
        onClick={handleBackButtonClick}
        disabled={page === 1}
        aria-label="Previous Page"
        size="large"
      >
        {theme.direction === 'rtl' ? <KeyboardArrowRight /> : <KeyboardArrowLeft />}
      </IconButton>
      <IconButton
        component={Link}
        nativeButton={false}
        {...getLinkProps(page + 1)}
        onClick={handleNextButtonClick}
        disabled={page >= Math.ceil(count / rowsPerPage)}
        aria-label="Next Page"
        size="large"
      >
        {theme.direction === 'rtl' ? <KeyboardArrowLeft /> : <KeyboardArrowRight />}
      </IconButton>
      <IconButton
        component={Link}
        nativeButton={false}
        {...getLinkProps(Math.ceil(count / rowsPerPage))}
        onClick={handleLastPageButtonClick}
        disabled={page >= Math.ceil(count / rowsPerPage)}
        aria-label="Last Page"
        size="large"
      >
        {theme.direction === 'rtl' ? <FirstPageIcon /> : <LastPageIcon />}
      </IconButton>
    </div>
  )
}

export default TablePaginationRouter
