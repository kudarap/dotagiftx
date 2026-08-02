import React from 'react'
import type { ReactNode } from 'react'

interface TabPanelProps extends React.HTMLAttributes<HTMLDivElement> {
  children: ReactNode
  value: unknown
  index: unknown
}

export default function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props
  const [visited, setVisited] = React.useState(value === index)

  if (value === index && !visited) {
    setVisited(true)
  }

  if (value !== index && !visited) {
    return null
  }

  return (
    <div hidden={value !== index} {...other}>
      {children}
    </div>
  )
}
