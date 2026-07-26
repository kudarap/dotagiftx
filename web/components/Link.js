import React from 'react'
import PropTypes from 'prop-types'
import clsx from 'clsx'
import { useRouter } from 'next/router'
import NextLink from 'next/link'
import MuiLink from '@mui/material/Link'

const nextLinkRef = (props, ref) => {
  const { as, href, disabled, ...other } = props

  return <NextLink href={disabled ? '#' : href} as={as} ref={ref} {...other} />
}

const NextComposed = React.forwardRef(nextLinkRef)

NextComposed.propTypes = {
  as: PropTypes.oneOfType([PropTypes.string, PropTypes.object]),
  href: PropTypes.oneOfType([PropTypes.string, PropTypes.object]),
  prefetch: PropTypes.bool,
  disabled: PropTypes.bool,
}

// A styled version of the Next.js Link component:
// https://nextjs.org/docs/#with-link
function Link(props) {
  const {
    href,
    activeClassName = 'active',
    className: classNameProps,
    innerRef,
    naked,
    disableUnderline,
    ...other
  } = props

  const router = useRouter()
  const pathname = typeof href === 'string' ? href : href.pathname
  const className = clsx(classNameProps, {
    [activeClassName]: router.pathname === pathname && activeClassName,
  })

  if (naked) {
    return <NextComposed className={className} ref={innerRef} href={href} {...other} />
  }

  return (
    <MuiLink
      component={NextComposed}
      className={className}
      ref={innerRef}
      href={href}
      color="textPrimary"
      style={disableUnderline ? { textDecoration: 'none' } : null}
      {...other}
    />
  )
}

Link.propTypes = {
  activeClassName: PropTypes.string,
  as: PropTypes.oneOfType([PropTypes.string, PropTypes.object]),
  className: PropTypes.string,
  href: PropTypes.oneOfType([PropTypes.string, PropTypes.object]),
  innerRef: PropTypes.oneOfType([PropTypes.func, PropTypes.object]),
  naked: PropTypes.bool,
  onClick: PropTypes.func,
  prefetch: PropTypes.bool,
  disableUnderline: PropTypes.bool,
  disabled: PropTypes.bool,
}

Link.defaultProps = {
  activeClassName: '',
  as: '',
  className: '',
  href: '',
  innerRef: '',
  naked: false,
  onClick: () => {},
  prefetch: null,
  disableUnderline: false,
  disabled: false,
}

const linkRef = (props, ref) => <Link {...props} innerRef={ref} />

export default React.forwardRef(linkRef)
