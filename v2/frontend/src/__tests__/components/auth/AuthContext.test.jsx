import React from "react";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { AuthProvider, useAuth } from "../../../components/auth/AuthContext";

// Mock useNavigate
const mockNavigate = jest.fn();
jest.mock("react-router-dom", () => ({
  ...jest.requireActual("react-router-dom"),
  useNavigate: () => mockNavigate,
}));

// Simple test component using AuthContext
const TestComponent = () => {
  const {
    authModalOpen,
    authMode,
    openAuthModal,
    closeAuthModal,
    switchAuthMode
  } = useAuth();

  return (
    <div>
      <button onClick={openAuthModal}>Open Modal</button>
      <button onClick={closeAuthModal}>Close Modal</button>
      <button onClick={() => switchAuthMode("signup")}>Switch to Signup</button>

      <div data-testid="auth-mode">{authMode}</div>
      {authModalOpen && <div data-testid="auth-modal">Modal Content</div>}
    </div>
  );
};

describe("AuthContext", () => {
  const renderWithAuth = () =>
    render(
      <BrowserRouter>
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      </BrowserRouter>
    );

  test("renders without crashing", () => {
    renderWithAuth();
    expect(screen.getByText("Open Modal")).toBeInTheDocument();
  });

  test("handles modal open/close", async () => {
    renderWithAuth();

    // Modal should not be visible initially
    expect(screen.queryByTestId("auth-modal")).not.toBeInTheDocument();

    // Open modal
    fireEvent.click(screen.getByText("Open Modal"));

    await waitFor(() => {
      expect(screen.getByTestId("auth-modal")).toBeInTheDocument();
    });

    // Close modal
    fireEvent.click(screen.getByText("Close Modal"));

    await waitFor(() => {
      expect(screen.queryByTestId("auth-modal")).not.toBeInTheDocument();
    });
  });

  test("switches auth mode correctly", async () => {
    renderWithAuth();

    // Should start with signin
    expect(screen.getByTestId("auth-mode")).toHaveTextContent("signin");

    // Switch to signup
    fireEvent.click(screen.getByText("Switch to Signup"));

    await waitFor(() => {
      expect(screen.getByTestId("auth-mode")).toHaveTextContent("signup");
    });
  });
});
