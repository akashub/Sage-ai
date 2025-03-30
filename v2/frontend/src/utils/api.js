// const API_BASE = process.env.NODE_ENV === 'production' ? '' : 'http://localhost:8080';

// export const uploadFile = async (file) => {
//   const formData = new FormData();
//   formData.append("file", file);
  
//   try {
//     const response = await fetch(`${API_BASE}/api/upload`, {
//       method: "POST",
//       body: formData,
//     });
    
//     if (!response.ok) {
//       throw new Error(`Upload failed: ${response.status}`);
//     }
    
//     return response.json();
//   } catch (error) {
//     console.error("API Error (uploadFile):", error);
//     throw error;
//   }
// };

// export const queryData = async (query, csvPath) => {
//   try {
//     const response = await fetch(`${API_BASE}/api/query`, {
//       method: "POST",
//       headers: {
//         "Content-Type": "application/json",
//       },
//       body: JSON.stringify({
//         query,
//         csvPath,
//       }),
//     });
    
//     if (!response.ok) {
//       throw new Error(`Query failed: ${response.status}`);
//     }
    
//     return response.json();
//   } catch (error) {
//     console.error("API Error (queryData):", error);
//     throw error;
//   }
// };

// src/utils/api.js

// src/utils/api.js

// Base API URL - adjust to match your deployment
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// Error handler for API responses
const handleApiError = (response) => {
  return response.json().then(data => {
    if (!response.ok) {
      const error = new Error(data.message || "API error");
      error.status = response.status;
      error.data = data;
      throw error;
    }
    return data;
  });
};

// CSV Data Operations

export const uploadFile = async (file) => {
  console.log("Uploading file:", file.name);
  const formData = new FormData();
  formData.append('file', file);

  try {
    const response = await fetch(`${API_BASE_URL}/api/upload`, {
      method: 'POST',
      body: formData,
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || `Upload failed with status: ${response.status}`);
    }
    
    const responseData = await response.json();
    console.log("Upload response data:", responseData);
    return responseData;
  } catch (error) {
    console.error("API Error (uploadFile):", error);
    throw error;
  }
};

export const queryData = async (query, csvPath, options = {}) => {
  console.log("Sending query:", query, csvPath, options);
  
  const response = await fetch(`${API_BASE_URL}/api/query`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query,
      csvPath,
      useKnowledgeBase: options.useKnowledgeBase !== false, // Default to true
      timestamp: options.timestamp || Date.now(), // Prevent caching
      options: options.additionalOptions || {},
    }),
    // Disable caching
    cache: "no-store"
  });

  return handleApiError(response);
};

// Training Data Operations

export const fetchTrainingData = async () => {
  console.log("Fetching training data");
  try {
    const response = await fetch(`${API_BASE_URL}/api/training/list`);
    
    if (!response.ok) {
      throw new Error(`Failed to fetch training data: ${response.status}`);
    }
    
    const data = await response.json();
    console.log("Fetched training data:", data);
    return data;
  } catch (error) {
    console.error("Error fetching training data:", error);
    return [];
  }
};

export const uploadTrainingFile = async (formData) => {
  console.log("Uploading training file");
  const response = await fetch(`${API_BASE_URL}/api/training/upload`, {
    method: 'POST',
    body: formData,
  });

  return handleApiError(response);
};

export const addTrainingData = async (data) => {
  console.log("Adding training data:", data);
  const response = await fetch(`${API_BASE_URL}/api/training/add`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  return handleApiError(response);
};

export const deleteTrainingData = async (id) => {
  console.log("Deleting training data:", id);
  const response = await fetch(`${API_BASE_URL}/api/training/delete/${id}`, {
    method: 'DELETE',
  });

  return handleApiError(response);
};

export const viewTrainingData = async (id) => {
  console.log("Viewing training data:", id);
  const response = await fetch(`${API_BASE_URL}/api/training/view/${id}`);
  
  return handleApiError(response);
};

// Authentication Operations

export const loginUser = async (credentials) => {
  const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(credentials),
  });

  return handleApiError(response);
};

export const registerUser = async (userData) => {
  const response = await fetch(`${API_BASE_URL}/api/auth/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  return handleApiError(response);
};

export const getUserProfile = async (token) => {
  const response = await fetch(`${API_BASE_URL}/api/auth/profile`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  return handleApiError(response);
};

export const getOAuthUrl = async (provider) => {
  const response = await fetch(`${API_BASE_URL}/api/auth/oauth/${provider}/url`);
  return handleApiError(response);
};

export const handleOAuthCallback = async (provider, code) => {
  const response = await fetch(`${API_BASE_URL}/api/auth/oauth/${provider}/callback`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ code }),
  });

  return handleApiError(response);
};