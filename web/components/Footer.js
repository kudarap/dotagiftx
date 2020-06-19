import React from 'react'
import MuiLink from '@material-ui/core/Link'
import Divider from '@material-ui/core/Divider'
import Container from '@/components/Container'
import Link from '@/components/Link'

export default function () {
  return (
    <footer style={{ marginTop: 20 }}>
      <Divider />
      <br />
      <Container disableMinHeight>
        <ul>
          <li>
            <Link href="/about">About</Link>
          </li>
          <li>
            <Link href="/about">FAQ</Link>
          </li>
          <li>
            <Link href="/privacy">Privacy</Link>
          </li>
          <li>
            <MuiLink href="http://vercel.com" target="_blank">
              Powered by Vercel
            </MuiLink>
          </li>
          <li>
            <MuiLink href="http://chiligarlic.com" target="_blank">
              A chiliGarlic project
            </MuiLink>
          </li>
        </ul>
      </Container>
    </footer>
  )
}
