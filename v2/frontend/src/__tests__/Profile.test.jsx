import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import Profile from '../pages/Profile';
import { AuthProvider } from '../components/auth/AuthContext';

// Mock the API functions
jest.mock('../utils/api', () => ({
  getUserProfile: jest.fn(),
  updateUserProfile: jest.fn(),
  getAPIKeys: jest.fn(),
  saveAPIKey: jest.fn(),
  deleteAPIKey: jest.fn(),
  setDefaultAPIKey: jest.fn()
}));

// Mock the Navigation component
jest.mock('../components/layout/Navigation', () => () => <div data-testid="navigation">Navigation</div>);

describe('Profile Page', () => {
  const mockUser = {
    name: 'Test User',
    email: 'test@example.com',
    picture: 'test-picture-url'
  };

  const mockProfile = {
    name: 'Test User',
    email: 'test@example.com',
    plan: 'Free',
    createdAt: '2024-01-01',
    stats: {
      totalChats: 10,
      totalQueries: 50
    }
  };

  const mockApiKeys = [
    {
      id: '1',
      name: 'Test Key',
      provider: 'gemini',
      maskedKey: '*****1234',
      isDefault: true,
      lastUsed: '2024-01-01'
    }
  ];

  beforeEach(() => {
    // Reset all mocks before each test
    jest.clearAllMocks();
    
    // Setup default mock implementations
    require('../utils/api').getUserProfile.mockResolvedValue(mockProfile);
    require('../utils/api').getAPIKeys.mockResolvedValue(mockApiKeys);
  });

  const renderProfile = () => {
    return render(
      <BrowserRouter>
        <AuthProvider value={{ user: mockUser }}>
          <Profile />
        </AuthProvider>
      </BrowserRouter>
    );
  };

  it('renders loading state initially', () => {
    require('../utils/api').getUserProfile.mockImplementation(() => new Promise(() => {}));
    renderProfile();
    expect(screen.getByText('Loading profile...')).toBeInTheDocument();
  });

  it('displays user profile information correctly', async () => {
    renderProfile();
    
    await waitFor(() => {
      expect(screen.getByText('Test User')).toBeInTheDocument();
      expect(screen.getByText('test@example.com')).toBeInTheDocument();
      expect(screen.getByText('Plan: Free')).toBeInTheDocument();
    });
  });

  it('displays usage statistics', async () => {
    renderProfile();
    
    await waitFor(() => {
      expect(screen.getByText('10')).toBeInTheDocument(); // Total chats
      expect(screen.getByText('50')).toBeInTheDocument(); // Total queries
      expect(screen.getByText('1')).toBeInTheDocument(); // API credentials
    });
  });

//   it('allows editing profile name', async () => {
//     renderProfile();
    
//     await waitFor(() => {
//       const editButton = screen.getByText('Edit Profile');
//       fireEvent.click(editButton);
//     });

//     const nameInput = screen.getByLabelText('Name');
//     fireEvent.change(nameInput, { target: { value: 'New Name' } });

//     const saveButton = screen.getByText('Save');
//     fireEvent.click(saveButton);

//     await waitFor(() => {
//       expect(require('../utils/api').updateUserProfile).toHaveBeenCalledWith({
//         name: 'New Name',
//         profilePicUrl: ''
//       });
//     });
//   });

  it('handles API key management', async () => {
    renderProfile();
    
    await waitFor(() => {
      const addKeyButton = screen.getByText('Add Key');
      fireEvent.click(addKeyButton);
    });

    // Fill in the API key form
    fireEvent.change(screen.getByLabelText('Key Name'), { target: { value: 'New API Key' } });
    fireEvent.change(screen.getByLabelText('API Key'), { target: { value: 'test-api-key' } });

    const saveButton = screen.getByText('Save');
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(require('../utils/api').saveAPIKey).toHaveBeenCalledWith({
        provider: 'gemini',
        apiKey: 'test-api-key',
        name: 'New API Key',
        isDefault: false
      });
    });
  });

  it('displays error message when profile loading fails', async () => {
    require('../utils/api').getUserProfile.mockRejectedValue(new Error('Failed to load profile'));
    renderProfile();
    
    await waitFor(() => {
      expect(screen.getByText('Failed to load profile data. Please refresh the page.')).toBeInTheDocument();
    });
  });

  it('handles empty API keys list', async () => {
    require('../utils/api').getAPIKeys.mockResolvedValue([]);
    renderProfile();
    
    await waitFor(() => {
      expect(screen.getByText('No API keys saved')).toBeInTheDocument();
      expect(screen.getByText('Add an API key to use with your chats')).toBeInTheDocument();
    });
  });

  it('displays default API key indicator', async () => {
    renderProfile();
    
    await waitFor(() => {
      expect(screen.getByText('Default')).toBeInTheDocument();
    });
  });
}); 