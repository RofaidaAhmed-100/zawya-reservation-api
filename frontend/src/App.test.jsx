import { render, screen } from '@testing-library/react'
import App from './App'

describe('App', () => {
  it('renders the Zawya brand name', () => {
    render(<App />)
    expect(screen.getByText('Zawya')).toBeInTheDocument()
  })

  it('has the dark background', () => {
    const { container } = render(<App />)
    expect(container.firstChild).toHaveClass('bg-slate-900')
  })
})
