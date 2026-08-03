import React from 'react'
import clsx from 'clsx'
import { useRouter } from 'next/router'
import NextLink from 'next/link'
import type { LinkProps as NextLinkProps } from 'next/link'
import MuiLink from '@mui/material/Link'
import type { LinkProps as MuiLinkProps } from '@mui/material/Link'

type NextComposedProps = Omit<NextLinkProps, 'as' | 'href'> &
  Omit<React.AnchorHTMLAttributes<HTMLAnchorElement>, 'href'> & {
    as?: NextLinkProps['as']
    href?: NextLinkProps['href']
    disabled?: boolean
  }

const NextComposed = React.forwardRef<HTMLAnchorElement, NextComposedProps>(
  function NextComposed(props, ref) {
    const { as, href, disabled, ...other } = props

    return (
      <NextLink
        href={disabled ? '#' : (href as NonNullable<NextLinkProps['href']>)}
        as={as}
        ref={ref}
        {...(other as object)}
      />
    )
  }
)

export type LinkProps = Omit<NextLinkProps, 'as' | 'href'> &
  Omit<MuiLinkProps, 'href'> & {
    activeClassName?: string
    as?: NextLinkProps['as']
    className?: string
    href?: NextLinkProps['href']
    innerRef?: React.Ref<HTMLAnchorElement>
    naked?: boolean
    disableUnderline?: boolean
    disabled?: boolean
  }

// A styled version of the Next.js Link component:
// https://nextjs.org/docs/#with-link
function Link(props: LinkProps) {
  const {
    href = '',
    activeClassName = 'active',
    className: classNameProps,
    innerRef,
    naked,
    disableUnderline,
    disabled,
    ...other
  } = props

  const router = useRouter()
  const pathname = typeof href === 'string' ? href : href.pathname
  const className = clsx(classNameProps, {
    [activeClassName]: router.pathname === pathname && activeClassName,
  })

  if (naked) {
    return (
      <NextComposed
        className={className}
        ref={innerRef}
        href={href}
        disabled={disabled}
        {...(other as object)}
      />
    )
  }

  return (
    <MuiLink
      component={NextComposed}
      className={className}
      ref={innerRef}
      href={href as string}
      color="textPrimary"
      style={disableUnderline ? { textDecoration: 'none' } : undefined}
      {...other}
    />
  )
}

const linkRef = (props: LinkProps, ref: React.ForwardedRef<HTMLAnchorElement>) => (
  <Link {...props} innerRef={ref} />
)

export default React.forwardRef<HTMLAnchorElement, LinkProps>(linkRef)
