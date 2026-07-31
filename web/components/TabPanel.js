import React from 'react'
import PropTypes from 'prop-types'

export default function TabPanel(props) {
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
TabPanel.propTypes = {
  children: PropTypes.node.isRequired,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired,
}
