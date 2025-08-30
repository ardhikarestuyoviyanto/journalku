import React, { useEffect, useState } from 'react'
import Select from 'react-select'

// Dark Theme
const customDarkTheme = (theme) => ({
  ...theme,
  colors: {
    ...theme.colors,
    primary25: '#2a2d32',
    primary: '#20a8d8',
    neutral0: '#212631',
    neutral80: '#fff',
    neutral20: '#555',
    neutral30: '#777',
    neutral50: '#aaa',
  },
})

const customDarkStyles = {
  control: (provided) => ({
    ...provided,
    backgroundColor: '##212631',
    borderColor: '#555',
    color: '#fff',
  }),
  menu: (provided) => ({
    ...provided,
    backgroundColor: '#2a2d32',
  }),
  singleValue: (provided) => ({
    ...provided,
    color: '#fff',
  }),
  input: (provided) => ({
    ...provided,
    color: '#fff',
  }),
  option: (provided, state) => ({
    ...provided,
    backgroundColor: state.isFocused ? '#343a40' : '#2a2d32',
    color: '#fff',
  }),
}

// Light Theme (default)
const customLightTheme = (theme) => ({
  ...theme,
})

const customLightStyles = {} // Default light styles

export default function Select2CustomTheme({ options, placeholder, onChange, ...rest }) {
  const [isDarkMode, setIsDarkMode] = useState(false)

  useEffect(() => {
    const checkTheme = () => {
      const theme = document.documentElement.getAttribute('data-coreui-theme')
      setIsDarkMode(theme === 'dark')
    }

    checkTheme()

    const observer = new MutationObserver(checkTheme)
    observer.observe(document.body, {
      attributes: true,
      attributeFilter: ['data-coreui-theme'],
    })

    return () => observer.disconnect()
  }, [])

  return (
    <Select
      options={options}
      placeholder={placeholder}
      onChange={onChange}
      theme={isDarkMode ? customDarkTheme : customLightTheme}
      styles={isDarkMode ? customDarkStyles : customLightStyles}
      {...rest}
    />
  )
}
