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
// src/utils/api.js

// Base API URL - adjust to match your deployment
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// Error handler for API responses
const handleApiError = (response) => {
  return response.json().then(data => {
    if (!response.ok) {
      const error = new Error(data.message || data.error || "API error");
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
    
    return handleApiError(response);
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
  const response = await fetch(`${API_BASE_URL}/api/training/list`);
  return handleApiError(response);
};

// export const uploadTrainingFile = async (formData) => {
//   console.log("Uploading training file");
//   const response = await fetch(`${API_BASE_URL}/api/training/upload`, {
//     method: 'POST',
//     body: formData,
//   });

//   return handleApiError(response);
// };
export const uploadTrainingFile = async (formData) => {
  console.log("Uploading training file");
  
  // Add timeout to prevent long-running requests
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 10000); // 10 second timeout
  
  try {
    const response = await fetch(`${API_BASE_URL}/api/training/upload`, {
      method: 'POST',
      body: formData,
      mode: 'cors',
      credentials: 'same-origin',
      signal: controller.signal,
      // Don't set Content-Type for FormData - browser will do it with boundary
    });
    
    clearTimeout(timeoutId);
    return handleApiError(response);
  } catch (error) {
    clearTimeout(timeoutId);
    
    // Handle abort error
    if (error.name === 'AbortError') {
      console.error("Request timed out after 10 seconds");
      throw new Error("Upload timed out. Server may be busy or unreachable.");
    }
    
    console.error("Error in uploadTrainingFile:", error);
    throw error;
  }
};

export const deleteTrainingData = async (id) => {
  console.log("Deleting training data:", id);
  
  // Add timeout to prevent long-running requests
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 5000); // 5 second timeout
  
  try {
    const response = await fetch(`${API_BASE_URL}/api/training/delete/${id}`, {
      method: 'DELETE',
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest'
      },
      signal: controller.signal
    });
    
    // Clear the timeout since request completed
    clearTimeout(timeoutId);
    
    // If status is 204 No Content, return success object
    if (response.status === 204) {
      return { success: true };
    }
    
    return handleApiError(response);
  } catch (error) {
    // Clear the timeout to prevent memory leaks
    clearTimeout(timeoutId);
    
    // Handle abort error more gracefully
    if (error.name === 'AbortError') {
      console.error("Request timed out after 5 seconds");
      return { success: true, warning: "Operation timed out but may have completed successfully" };
    }
    
    console.error("Error in deleteTrainingData:", error);
    throw error;
  }
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

// export const deleteTrainingData = async (id) => {
//   console.log("Deleting training data:", id);
//   const response = await fetch(`${API_BASE_URL}/api/training/delete/${id}`, {
//     method: 'DELETE',
//   });

//   return handleApiError(response);
// };

export const viewTrainingData = async (id) => {
  console.log("Viewing training data:", id);
  const response = await fetch(`${API_BASE_URL}/api/training/view/${id}`);
  return handleApiError(response);
};

// Chat History Operations
export const fetchChatHistory = async () => {
  console.log("Fetching chat history");
  const response = await fetch(`${API_BASE_URL}/api/chats`);
  return handleApiError(response);
};

export const fetchChatById = async (chatId) => {
  console.log("Fetching chat by ID:", chatId);
  const response = await fetch(`${API_BASE_URL}/api/chats/${chatId}`);
  return handleApiError(response);
};

export const createChat = async (data = {}) => {
  console.log("Creating new chat");
  
  // Add timeout to prevent long-running requests
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 5000); // 5 second timeout
  
  try {
    const response = await fetch(`${API_BASE_URL}/api/chats`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest'
      },
      body: JSON.stringify({
        title: data.title || "New Chat",
        timestamp: data.timestamp || new Date().toISOString(),
        ...data
      }),
      signal: controller.signal,
      mode: 'cors' // Explicitly set CORS mode
    });
    
    // Clear the timeout since request completed
    clearTimeout(timeoutId);
    
    return handleApiError(response);
  } catch (error) {
    // Clear the timeout to prevent memory leaks
    clearTimeout(timeoutId);
    
    // Handle abort error more gracefully
    if (error.name === 'AbortError') {
      console.error("Request timed out after 5 seconds");
      // Return a fallback object to prevent further errors
      return { 
        id: `temp_${Date.now()}`,
        title: data.title || "New Chat",
        timestamp: new Date().toISOString(),
        ...data,
        _warning: "Created locally due to server timeout"
      };
    }
    
    console.error("Error in createChat:", error);
    throw error;
  }
};

export const updateChat = async (chatId, data) => {
  console.log("Updating chat:", chatId, data);
  const response = await fetch(`${API_BASE_URL}/api/chats/${chatId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  return handleApiError(response);
};

export const deleteChat = async (chatId) => {
  console.log("Deleting chat:", chatId);
  const response = await fetch(`${API_BASE_URL}/api/chats/${chatId}`, {
    method: 'DELETE',
  });

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