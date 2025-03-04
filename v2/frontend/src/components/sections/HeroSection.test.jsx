import { render, screen, fireEvent } from "@testing-library/react"
import "@testing-library/jest-dom"
import HeroSection from "./HeroSection"

describe("HeroSection", () => {
  it("renders the main heading", () => {
    render(<HeroSection />)
    expect(screen.getByText("Generate")).toBeInTheDocument()
    expect(screen.getByText("SQL")).toBeInTheDocument()
  })

  it("renders the subheading", () => {
    render(<HeroSection />)
    expect(screen.getByText("the easiest way to create database queries with AI")).toBeInTheDocument()
  })

  it("renders the Try Generator button", () => {
    render(<HeroSection />)
    expect(screen.getByText("Try Generator")).toBeInTheDocument()
  })

  it("renders the Watch Demo button", () => {
    render(<HeroSection />)
    expect(screen.getByText("Watch Demo")).toBeInTheDocument()
  })

  it("calls onAuthOpen when Try Generator button is clicked", () => {
    const onAuthOpenMock = jest.fn()
    render(<HeroSection onAuthOpen={onAuthOpenMock} />)

    // Click the Try Generator button
    fireEvent.click(screen.getByText("Try Generator"))

    // Check if onAuthOpen was called
    expect(onAuthOpenMock).toHaveBeenCalled()
  })

  it("renders the floating SQL snippets", () => {
    render(<HeroSection />)

    // Check if the SQL snippets are rendered
    expect(screen.getByText(/SELECT \* FROM users WHERE created_at/)).toBeInTheDocument()
    expect(screen.getByText(/INSERT INTO orders/)).toBeInTheDocument()
    expect(screen.getByText(/UPDATE products SET stock/)).toBeInTheDocument()
    expect(screen.getByText(/DELETE FROM cart/)).toBeInTheDocument()
  })
})

