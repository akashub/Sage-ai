import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import About from '../pages/About';
import '@testing-library/jest-dom';

// Mock the Navigation component
jest.mock('../components/layout/Navigation', () => () => <div data-testid="mock-navigation" />);

describe('About Page', () => {
  beforeEach(() => {
    render(
      <BrowserRouter>
        <About />
      </BrowserRouter>
    );
  });

  it('renders the navigation component', () => {
    expect(screen.getByTestId('mock-navigation')).toBeInTheDocument();
  });

  it('displays the main heading', () => {
    expect(screen.getByText('About Sage.AI')).toBeInTheDocument();
  });

  it('displays the mission statement', () => {
    expect(screen.getByText('Our Mission')).toBeInTheDocument();
    expect(screen.getByText(/democratize data analysis/i)).toBeInTheDocument();
  });

  it('displays all feature cards', () => {
    const features = [
      'Advanced AI',
      'Enterprise Security',
      'Lightning Fast'
    ];

    features.forEach(feature => {
      expect(screen.getByText(feature)).toBeInTheDocument();
    });
  });

  it('displays all team members', () => {
    const teamMembers = [
      'Akash Singh',
      'Nitin Reddy',
      'Sudiksha Rajavaram',
      'Yash Kishore'
    ];

    teamMembers.forEach(member => {
      expect(screen.getByText(member)).toBeInTheDocument();
    });
  });

  it('displays the call to action section', () => {
    expect(screen.getByText('Ready to Transform Your Data Analysis?')).toBeInTheDocument();
    expect(screen.getByText('Get Started')).toBeInTheDocument();
  });

  it('renders all feature icons', () => {
    const icons = [
      'CodeIcon',
      'SecurityIcon',
      'SpeedIcon'
    ];

    icons.forEach(icon => {
      expect(screen.getByTestId(icon)).toBeInTheDocument();
    });
  });

  it('renders team member avatars', () => {
    const avatars = screen.getAllByRole('img');
    expect(avatars.length).toBe(4); // 4 team members
  });

  it('displays feature descriptions', () => {
    const descriptions = [
      /state-of-the-art language models/i,
      /bank-grade encryption/i,
      /minimal latency/i
    ];

    descriptions.forEach(description => {
      expect(screen.getByText(description)).toBeInTheDocument();
    });
  });


}); 