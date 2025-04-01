import React from 'react';
import { render, screen } from '@testing-library/react';
import Footer from '../../../components/layout/Footer';

// Mock the current year to ensure consistent testing
const mockYear = 2024;
jest.spyOn(global.Date.prototype, 'getFullYear').mockReturnValue(mockYear);

// Mock AuthContext
jest.mock('../../../components/auth/AuthContext', () => ({
  useAuth: () => ({
    openAuthModal: jest.fn(),
  }),
}));

describe('Footer', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders footer component', () => {
    render(<Footer />);
    const footer = screen.getByRole('contentinfo');
    expect(footer).toBeInTheDocument();
  });

  test('renders logo', () => {
    render(<Footer />);
    const logo = screen.getByAltText('SAGE.AI');
    expect(logo).toBeInTheDocument();
  });

  test('renders copyright text', () => {
    render(<Footer />);
    const copyright = screen.getByText(`© ${mockYear} SAGE.AI Inc. All Rights Reserved.`);
    expect(copyright).toBeInTheDocument();
  });

  test('renders social media links', () => {
    render(<Footer />);
    const links = screen.getAllByRole('link');
    expect(links.length).toBeGreaterThan(0);
  });

  test('renders CTA buttons', () => {
    render(<Footer />);
    const launchButton = screen.getByRole('button', { name: /launch app/i });
    const joinButton = screen.getByRole('button', { name: /join early access/i });
    expect(launchButton).toBeInTheDocument();
    expect(joinButton).toBeInTheDocument();
  });
}); 