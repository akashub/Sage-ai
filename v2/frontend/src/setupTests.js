// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
import "@testing-library/jest-dom"

// Mock the framer-motion
jest.mock("framer-motion", () => {
  const actual = jest.requireActual("framer-motion")
  return {
    ...actual,
    motion: {
      div: ({ children, ...props }) => <div {...props}>{children}</div>,
      span: ({ children, ...props }) => <span {...props}>{children}</span>,
    },
    AnimatePresence: ({ children }) => <>{children}</>,
  }
})

// Mock the Material-UI components that might cause issues
jest.mock("@mui/material", () => {
  const actual = jest.requireActual("@mui/material")
  return {
    ...actual,
    useTheme: () => ({
      palette: {
        mode: "dark",
        primary: { main: "#5865F2" },
        secondary: { main: "#EB459E" },
        background: { default: "#1a1c20", paper: "#2F3136" },
        text: { primary: "#FFFFFF", secondary: "rgba(255, 255, 255, 0.7)" },
      },
      breakpoints: {
        down: () => false,
        up: () => true,
      },
    }),
  }
})

