import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { AuthProvider, useAuth } from "../../../components/auth/AuthContext";

// Mock useNavigate
const mockNavigate = jest.fn();
jest.mock("react-router-dom", () => ({
  ...jest.requireActual("react-router-dom"),
  useNavigate: () => mockNavigate,
}));

// Simple test component
const TestComponent = () => {
  const { authModalOpen, openAuthModal, closeAuthModal, authMode, switchAuthMode } = useAuth();
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
  const renderWithAuth = () => {
    return render(
      <BrowserRouter>
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      </BrowserRouter>
    );
  };

  test("renders without crashing", () => {
    renderWithAuth();
    expect(screen.getByText("Open Modal")).toBeInTheDocument();
  });

  // test("handles modal open/close", () => {
  //   renderWithAuth();
    
  //   // Modal should not be visible initially
  //   expect(screen.queryByTestId("auth-modal")).not.toBeInTheDocument();
    
  //   // Open modal
  //   fireEvent.click(screen.getByText("Open Modal"));
  //   expect(screen.getByTestId("auth-modal")).toBeInTheDocument();
    
  //   // Close modal
  //   fireEvent.click(screen.getByText("Close Modal"));
  //   expect(screen.queryByTestId("auth-modal")).not.toBeInTheDocument();
  // });

  test("handles auth mode switch", () => {
    renderWithAuth();
    
    // Should start with signin mode
    expect(screen.getByTestId("auth-mode")).toHaveTextContent("signin");
    
    // Switch to signup mode
    fireEvent.click(screen.getByText("Switch to Signup"));
    expect(screen.getByTestId("auth-mode")).toHaveTextContent("signup");
  });
});