import React from 'react'
import filter from 'lodash/filter'
import TextField from '@material-ui/core/TextField'
import Autocomplete from '@material-ui/lab/Autocomplete'
import CircularProgress from '@material-ui/core/CircularProgress'
import { catalogSearch } from '@/service/api'

const catalogSearchFilter = { limit: 1000, sort: 'popular' }

export default function ItemAutoComplete({ onSelect }) {
  const [open, setOpen] = React.useState(false)
  const [options, setOptions] = React.useState([])
  const loading = open && options.length === 0

  React.useEffect(() => {
    let active = true

    if (!loading) {
      return undefined
    }

    ;(async () => {
      const catalogs = await catalogSearch(catalogSearchFilter)

      if (active) {
        setOptions(catalogs.data)
      }
    })()

    return () => {
      active = false
    }
  }, [loading])

  React.useEffect(() => {
    if (!open) {
      setOptions([])
    }
  }, [open])

  const handleInputChange = (e, text) => {
    const res = filter(options, { name: text })
    if (res.length === 0) {
      onSelect({})
      return
    }

    onSelect(res[0])
  }

  return (
    <Autocomplete
      id="asynchronous-item-search"
      fullWidth
      open={open}
      onOpen={() => {
        setOpen(true)
      }}
      onClose={() => {
        setOpen(false)
      }}
      onInputChange={handleInputChange}
      getOptionSelected={(option, value) => option.name === value.name}
      getOptionLabel={option => option.name}
      options={options}
      loading={loading}
      renderInput={params => (
        <TextField
          {...params}
          autoFocus
          color="secondary"
          label="Item name"
          helperText="Search item you want to post from your inventory."
          variant="outlined"
          InputProps={{
            ...params.InputProps,
            endAdornment: (
              <>
                {loading ? <CircularProgress color="inherit" size={20} /> : null}
                {params.InputProps.endAdornment}
              </>
            ),
          }}
        />
      )}
    />
  )
}
