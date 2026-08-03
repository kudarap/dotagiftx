import Box from '@mui/material/Box'
import Drawer from '@mui/material/Drawer'
import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'
import Divider from '@mui/material/Divider'
import Link from './Link'
import type { Profile } from '@/lib/types'

const convertToNav = (nav: string[]) => ({ label: nav[0], path: nav[1] })

const primaryLinks = [
  ['Home', '/'],
  ['Post item', '/post-item'],
].map(convertToNav)

const secondaryLinks = [
  ['Treasures', '/treasures'],
  ['Heroes', '/heroes'],
  ['Rules', '/rules'],
  ['Bans', '/bans'],
  ['Guides', '/guides'],
  ['FAQs', '/faqs'],
  ['Middleman', '/middleman'],
  ['Moderators', '/moderators'],
  ['Updates', '/updates'],
  ['Mobile', '/download'],
].map(convertToNav)

function MenuDrawer({
  profile,
  open,
  onClose,
}: {
  profile: Profile
  open: boolean
  onClose: () => void
}) {
  const links = [...primaryLinks]
  if (!profile.id) {
    links.splice(1, 0, convertToNav(['Login', '/login']))
  }

  return (
    <Drawer anchor="right" open={open} onClose={onClose}>
      <Box sx={{ width: 250 }} role="presentation" onClick={onClose} onKeyDown={onClose}>
        <List>
          {links.map(link => (
            <ListItemButton
              key={link.path}
              component={Link}
              nativeButton={false}
              href={link.path}
            >
              <ListItemText primary={link.label} />
            </ListItemButton>
          ))}
        </List>
        <Divider />
        <List>
          {secondaryLinks.map(link => (
            <ListItemButton
              key={link.path}
              component={Link}
              nativeButton={false}
              href={link.path}
            >
              <ListItemText primary={link.label} />
            </ListItemButton>
          ))}
          <ListItemButton
            key="discord"
            component={Link}
            nativeButton={false}
            href="https://discord.gg/UFt9Ny42kM"
            target="_blank"
            rel="noreferrer noopener"
          >
            <ListItemText primary="Discord" />
          </ListItemButton>
        </List>
        <Divider />
        <ListItemButton component={Link} nativeButton={false} href="/plus">
          <ListItemText
            primary={
              <span>
                Dotagift<span style={{ fontSize: 18 }}>+</span>
              </span>
            }
          />
        </ListItemButton>
      </Box>
    </Drawer>
  )
}

export default MenuDrawer
