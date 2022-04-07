import React, { useContext, useState } from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import InputBase from '@mui/material/InputBase'
import Typography from '@mui/material/Typography'
import Dialog from '@mui/material/Dialog'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Grid from '@mui/material/Grid'
import Box from '@mui/material/Box'
import Drawer from '@mui/material/Drawer'
import Button from '@mui/material/Button'
import List from '@mui/material/List'
import ListItem from '@mui/material/ListItem'
import ListItemIcon from '@mui/material/ListItemIcon'
import ListItemText from '@mui/material/ListItemText'
import InboxIcon from '@mui/icons-material/MoveToInbox'
import MailIcon from '@mui/icons-material/Mail'
import Divider from '@mui/material/Divider'
import SearchIcon from '@mui/icons-material/Search'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'
import Link from './Link'

const primaryLinks = [
  ['Home', '/'],
  ['Treasures', '/treasures'],
  ['Rules', '/rules'],
  ['Bans', '/banned-users'],
  ['Guides', '/guides'],
].map(n => ({ label: n[0], path: n[1] }))

const secondaryLinks = [
  ['FAQs', '/faqs'],
  ['Updates', '/updates'],
  ['Middleman', '/middlemen'],
].map(n => ({ label: n[0], path: n[1] }))

function MenuDrawer({ profile, open, onClose }) {
  return (
    <>
      <Drawer anchor="right" open={open} onClose={onClose}>
        <Box sx={{ width: 250 }} role="presentation" onClick={onClose} onKeyDown={onClose}>
          <List>
            {primaryLinks.map(link => (
              <ListItem button key={link.path} component={Link} href={link.path}>
                <ListItemText primary={link.label} />
              </ListItem>
            ))}
          </List>
          <Divider />
          <List>
            {secondaryLinks.map(link => (
              <ListItem button key={link.path} component={Link} href={link.path}>
                <ListItemText primary={link.label} />
              </ListItem>
            ))}
            <ListItem
              button
              key="discord"
              component={Link}
              href="https://discord.gg/UFt9Ny42kM"
              target="_blank"
              rel="noreferrer noopener">
              <ListItemText primary="Discord" />
            </ListItem>
          </List>
          <Divider />
          <ListItem button component={Link} href="/plus">
            <ListItemText
              primary={
                <span>
                  Dotagift<span style={{ fontSize: 18 }}>+</span>
                </span>
              }
            />
          </ListItem>
        </Box>
      </Drawer>
    </>
  )
}
MenuDrawer.propTypes = {
  profile: PropTypes.object,
  open: PropTypes.bool,
  onClose: PropTypes.func,
}
MenuDrawer.defaultProps = {
  profile: {},
  open: false,
  onClose: () => {},
}

export default MenuDrawer
