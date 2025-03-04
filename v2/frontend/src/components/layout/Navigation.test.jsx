import { render, screen, fireEvent } from "@testing-library/react"
import "@testing-library/jest-dom"
import { BrowserRouter } from "react-router-dom"
import Navigation from "./Navigation"
import { AuthProvider } from "../auth/AuthContext"

// Mock the useMediaQuery hook
jest.mock("@mui/material", () => {
  const actual = jest.requireActual("@mui/material")
  return {
    ...actual,
    useMediaQuery: () => false, // Always return desktop view
  }
})

describe("Navigation", () => {
  const renderNavigation = () => {
    return render(
      <BrowserRouter>
        <AuthProvider>
          <Navigation />
        </AuthProvider>
      </BrowserRouter>,
    )
  }

  it("renders the logo", () => {
    renderNavigation()
    const logo = screen.getByAltText("SQL AI")
    expect(logo).toBeInTheDocument()
  })

  it("renders navigation links", () => {
    renderNavigation()
    expect(screen.getByText("Features")).toBeInTheDocument()
    expect(screen.getByText("Docs")).toBeInTheDocument()
    expect(screen.getByText("About")).toBeInTheDocument()
  })

  it("renders the Launch App button", () => {
    renderNavigation()
    const launchAppButton = screen.getByText("Launch App")
    expect(launchAppButton).toBeInTheDocument()
    expect(launchAppButton.closest("a")).toHaveAttribute("href", "/chat")
  })

  it("renders the Join Early Access button", () => {
    renderNavigation()
    expect(screen.getByText("Join Early Access")).toBeInTheDocument()
  })

  it("opens auth modal when Join Early Access is clicked", () => {
    const openAuthModalMock = jest.fn()

    // Mock the useAuth hook
    jest.spyOn(require("../auth/AuthContext"), "useAuth").mockImplementation(() => ({
      openAuthModal: openAuthModalMock,
    }))

    renderNavigation()

    // Click the Join Early Access button
    fireEvent.click(screen.getByText("Join Early Access"))

    // Check if openAuthModal was called
    expect(openAuthModalMock).toHaveBeenCalled()
  })
})

