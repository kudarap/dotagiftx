import React from 'react'
import PropTypes from 'prop-types'
import Box from '@mui/material/Box'
import Drawer from '@mui/material/Drawer'
import List from '@mui/material/List'
import ListItem from '@mui/material/ListItem'
import ListItemText from '@mui/material/ListItemText'
import Divider from '@mui/material/Divider'
import Link from './Link'

const convertToNav = ([label, path]) => ({ label, path })

const primaryLinks = [
  ['Home', '/'],
  ['Post item', '/post-item'],
].map(convertToNav)

const secondaryLinks = [
  ['Treasures', '/treasures'],
  ['Rules', '/rules'],
  ['Bans', '/bans'],
  ['Guides', '/guides'],
  ['FAQs', '/faqs'],
  ['Middleman', '/middleman'],
  ['Moderators', '/moderators'],
  ['Updates', '/updates'],
].map(convertToNav)

function MenuDrawer({ profile, open, onClose }) {
  let links = [...primaryLinks]
  if (!profile.id) {
    links.splice(1, 0, convertToNav(['Login', '/login']))
  }

  return (
    <>
      <Drawer anchor="right" open={open} onClose={onClose}>
        <Box sx={{ width: 250 }} role="presentation" onClick={onClose} onKeyDown={onClose}>
          <List>
            {links.map(link => (
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
