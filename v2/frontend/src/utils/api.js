
// src/utils/api.js
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
  
  // Ensure llmConfig is properly included
  const requestBody = {
    query,
    csvPath,
    useKnowledgeBase: options.useKnowledgeBase !== false,
    timestamp: options.timestamp || Date.now(),
    options: options.additionalOptions || {},
  };
  
  // Only add llmConfig if it's provided in options
  if (options.llmConfig) {
    requestBody.llmConfig = options.llmConfig;
  }
  
  // Add user ID if available
  let userId = "user123"; // Default
  try {
    const userInfo = JSON.parse(localStorage.getItem("userInfo"));
    if (userInfo && userInfo.id) {
      requestBody.userId = userInfo.id;
    }
  } catch (error) {
    console.error("Error getting user info:", error);
  }
  
  console.log("Request body:", requestBody);
  
  const response = await fetch(`${API_BASE_URL}/api/query`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    },
    body: JSON.stringify(requestBody),
    cache: "no-store"
  });

  return handleApiError(response);
};

// Training Data Operations
export const fetchTrainingData = async () => {
  console.log("Fetching training data");
  try {
    const response = await fetch(`${API_BASE_URL}/api/training/list`);
    
    // Check if the response is ok before processing
    if (!response.ok) {
      console.error(`Error response: ${response.status} ${response.statusText}`);
      const errorText = await response.text();
      console.error(`Error details: ${errorText}`);
      throw new Error(`Failed to fetch training data: ${response.status} ${response.statusText}`);
    }
    
    // Check if the response is empty
    const text = await response.text();
    if (!text || text.trim() === "") {
      console.log("Received empty response");
      return [];
    }
    
    // Try to parse as JSON
    try {
      const data = JSON.parse(text);
      console.log("Training data parsed successfully:", data);
      return data;
    } catch (parseError) {
      console.error("Failed to parse response as JSON:", parseError);
      console.error("Raw response:", text);
      throw new Error("Invalid JSON response from server");
    }
  } catch (error) {
    console.error("Error fetching training data:", error);
    // Return empty array on error to prevent UI crashes
    return [];
  }
};

export const uploadTrainingFile = async (formData) => {
  console.log("Uploading training file");
  
  // Add timeout to prevent long-running requests
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 30000); // 30 second timeout
  
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
    
    // Log the response status
    console.log(`Upload response status: ${response.status} ${response.statusText}`);
    
    // Get the response text first for debugging
    const responseText = await response.text();
    console.log("Response text:", responseText);
    
    // Try to parse as JSON if possible
    try {
      return responseText ? JSON.parse(responseText) : { success: true };
    } catch (err) {
      console.warn("Response not valid JSON:", err);
      return { success: true, warning: "Response was not JSON" };
    }
  } catch (error) {
    clearTimeout(timeoutId);
    
    // Handle abort error
    if (error.name === 'AbortError') {
      console.error("Request timed out after 30 seconds");
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
  const timeoutId = setTimeout(() => controller.abort(), 10000); // 10 second timeout
  
  try {
    // Make sure the URL is correctly formatted
    const url = `${API_BASE_URL}/api/training/delete/${id}`;
    console.log("Delete URL:", url);
    
    const response = await fetch(url, {
      method: 'DELETE',
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest',
        'Accept': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      signal: controller.signal
    });
    
    // Clear the timeout since request completed
    clearTimeout(timeoutId);
    
    console.log(`Delete response status: ${response.status} ${response.statusText}`);
    
    // If status is 204 No Content or 200 OK, return success object
    if (response.status === 204 || response.status === 200) {
      return { success: true };
    }
    
    // Try to parse response
    try {
      const text = await response.text();
      console.log("Raw response:", text);
      
      // Try to parse as JSON if possible
      const data = text ? JSON.parse(text) : { success: false };
      
      if (!response.ok) {
        throw new Error(data.message || data.error || `Server returned ${response.status}`);
      }
      
      return data;
    } catch (parseError) {
      console.error("Parse error:", parseError);
      // If can't parse as JSON but response is ok, consider it a success
      if (response.ok) {
        return { success: true };
      }
      
      throw new Error(`Delete failed: ${response.status} ${response.statusText}`);
    }
  } catch (error) {
    // Clear the timeout to prevent memory leaks
    clearTimeout(timeoutId);
    
    // Handle abort error more gracefully
    if (error.name === 'AbortError') {
      console.error("Request timed out after 10 seconds");
      // Return a partial success to prevent UI inconsistency
      return { 
        success: true, 
        warning: "Operation timed out but may have completed successfully",
        id: id
      };
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
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    },
    body: JSON.stringify(data),
  });

  return handleApiError(response);
};

export const viewTrainingData = async (id) => {
  console.log("Viewing training data:", id);
  
  try {
    const response = await fetch(`${API_BASE_URL}/api/training/view/${id}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    // Check if the response is ok
    if (!response.ok) {
      console.error(`Error response: ${response.status} ${response.statusText}`);
      const errorText = await response.text();
      console.error(`Error details: ${errorText}`);
      throw new Error(`Failed to view training data: ${response.status} ${response.statusText}`);
    }
    
    // Check if the response is empty
    const text = await response.text();
    if (!text || text.trim() === "") {
      console.log("Received empty response");
      return {
        id: id,
        type: "unknown",
        description: "Unknown item",
        content: "No content available",
        date_added: new Date().toISOString()
      };
    }
    
    // Try to parse as JSON
    try {
      const data = JSON.parse(text);
      console.log("Training data item parsed successfully:", data);
      return data;
    } catch (parseError) {
      console.error("Failed to parse response as JSON:", parseError);
      console.error("Raw response:", text);
      throw new Error("Invalid JSON response from server");
    }
  } catch (error) {
    console.error("Error fetching training data item:", error);
    // Return a valid but error-indicating object
    return {
      id: id,
      type: "error",
      description: "Error fetching item",
      content: `Error: ${error.message || "Unknown error"}`,
      date_added: new Date().toISOString()
    };
  }
};

// Chat History Operations
export const fetchChatHistory = async () => {
  console.log("Fetching chat history");
  
  // Get current user ID
  let userId = "user123"; // Default
  try {
    const userInfo = JSON.parse(localStorage.getItem("userInfo"));
    if (userInfo && userInfo.id) {
      userId = userInfo.id;
    }
  } catch (error) {
    console.error("Error getting user info:", error);
  }
  
  const response = await fetch(`${API_BASE_URL}/api/chats?userId=${userId}`, {
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  });
  
  return handleApiError(response);
};

export const fetchChatById = async (chatId) => {
  console.log("Fetching chat by ID:", chatId);
  const response = await fetch(`${API_BASE_URL}/api/chats/${chatId}`, {
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  });
  return handleApiError(response);
};

export const createChat = async (data = {}) => {
  console.log("Creating new chat");
  
  // Get current user ID
  let userId = "user123"; // Default
  try {
    const userInfo = JSON.parse(localStorage.getItem("userInfo"));
    if (userInfo && userInfo.id) {
      userId = userInfo.id;
    }
  } catch (error) {
    console.error("Error getting user info:", error);
  }
  
  // Add user ID to request
  const chatData = {
    ...data,
    userId,
    title: data.title || "New Chat",
    timestamp: data.timestamp || new Date().toISOString()
  };
  
  // Add timeout to prevent long-running requests
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 5000); // 5 second timeout
  
  try {
    const response = await fetch(`${API_BASE_URL}/api/chats`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(chatData),
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
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    },
    body: JSON.stringify(data),
  });

  return handleApiError(response);
};

export const deleteChat = async (chatId) => {
  console.log("Deleting chat:", chatId);
  
  // Add timeout to prevent long-running requests
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 10000); // 10 second timeout
  
  try {
    const response = await fetch(`${API_BASE_URL}/api/chats/${chatId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      signal: controller.signal,
      mode: 'cors' // Explicitly set CORS mode
    });
    
    // Clear the timeout since request completed
    clearTimeout(timeoutId);
    
    // If status is 204 No Content, return success object
    if (response.status === 204) {
      console.log(`Chat ${chatId} deleted successfully with status 204`);
      return { success: true };
    }
    
    console.log(`Delete response status: ${response.status} ${response.statusText}`);
    
    // For other status codes, try to parse response
    try {
      const data = await response.json();
      console.log("Delete response data:", data);
      
      if (!response.ok) {
        throw new Error(data.message || data.error || `Server returned ${response.status}`);
      }
      
      return data;
    } catch (parseError) {
      // If can't parse as JSON, check if response is ok
      if (response.ok) {
        return { success: true };
      }
      
      throw new Error(`Delete failed: ${response.status} ${response.statusText}`);
    }
  } catch (error) {
    // Clear the timeout to prevent memory leaks
    clearTimeout(timeoutId);
    
    // Handle abort error more gracefully
    if (error.name === 'AbortError') {
      console.error("Delete request timed out after 10 seconds");
      return { 
        success: false, 
        error: "Operation timed out",
        chatId: chatId
      };
    }
    
    console.error(`Error deleting chat ${chatId}:`, error);
    throw error;
  }
};

// Chat Training Data Operations
export const updateChatTrainingData = async (chatId, trainingDataIds) => {
  console.log("Updating training data for chat:", chatId);
  if (!chatId) {
    return false;
  }
  
  const response = await fetch(`${API_BASE_URL}/api/chats/${chatId}/training`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    },
    body: JSON.stringify({
      trainingDataIds: trainingDataIds
    })
  });
  
  return handleApiError(response);
};

export const getChatTrainingData = async (chatId) => {
  console.log("Getting training data for chat:", chatId);
  if (!chatId) {
    return [];
  }
  
  const response = await fetch(`${API_BASE_URL}/api/chats/${chatId}/training`, {
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  });
  
  return handleApiError(response);
};

// Authentication Operations
export const loginUser = async (credentials) => {
  const response = await fetch(`${API_BASE_URL}/api/auth/signin`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(credentials),
  });

  const data = await handleApiError(response);
  
  // Save user info in localStorage
  if (data.user) {
    localStorage.setItem("userInfo", JSON.stringify(data.user));
    localStorage.setItem("token", data.accessToken);
  }
  
  return data;
};

export const registerUser = async (userData) => {
  const response = await fetch(`${API_BASE_URL}/api/auth/signup`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  const data = await handleApiError(response);
  
  // Save user info in localStorage
  if (data.user) {
    localStorage.setItem("userInfo", JSON.stringify(data.user));
    localStorage.setItem("token", data.accessToken);
  }
  
  return data;
};

export const logoutUser = async () => {
  // Clear user data from localStorage
  localStorage.removeItem("userInfo");
  localStorage.removeItem("token");
  
  // Call the signout endpoint
  try {
    const response = await fetch(`${API_BASE_URL}/api/auth/signout`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    return handleApiError(response);
  } catch (error) {
    console.error("Error in logout:", error);
    // Even if the server request fails, we've cleared local storage
    return { success: true };
  }
};

export const getUserProfile = async () => {
  const response = await fetch(`${API_BASE_URL}/api/profile`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  });
  
  return handleApiError(response);
};

export const updateUserProfile = async (updates) => {
  const response = await fetch(`${API_BASE_URL}/api/profile`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    },
    body: JSON.stringify(updates)
  });
  
  return handleApiError(response);
};

// API Key Management
export const getAPIKeys = async () => {
  const response = await fetch(`${API_BASE_URL}/api/apikeys`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  });
  
  return handleApiError(response);
};

export const saveAPIKey = async (keyData) => {
  const response = await fetch(`${API_BASE_URL}/api/apikeys`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    },
    body: JSON.stringify(keyData)
  });
  
  return handleApiError(response);
};

export const deleteAPIKey = async (keyId) => {
  const response = await fetch(`${API_BASE_URL}/api/apikeys/${keyId}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  });
  
  if (response.status === 204) {
    return { success: true };
  }
  
  return handleApiError(response);
};

export const setDefaultAPIKey = async (keyId) => {
  const response = await fetch(`${API_BASE_URL}/api/apikeys/${keyId}/default`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  });
  
  if (response.status === 204) {
    return { success: true };
  }
  
  return handleApiError(response);
};

// OAuth Helper Functions
export const getOAuthUrl = async (provider) => {
  const redirectUri = `${window.location.origin}/oauth-callback`;
  const response = await fetch(`${API_BASE_URL}/api/auth/oauth/url/${provider}?redirect_uri=${encodeURIComponent(redirectUri)}`);
  return handleApiError(response);
};

export const handleOAuthCallback = async (provider, code) => {
  const redirectUri = `${window.location.origin}/oauth-callback`;
  const response = await fetch(`${API_BASE_URL}/api/auth/oauth/${provider}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ 
      code,
      redirect_uri: redirectUri
    }),
  });

  const data = await handleApiError(response);
  
  // Save user info in localStorage
  if (data.user) {
    localStorage.setItem("userInfo", JSON.stringify(data.user));
    localStorage.setItem("token", data.accessToken);
  }
  
  return data;
};

// LLM Configuration Validation
export const validateApiKey = async (provider, apiKey) => {
  const response = await fetch(`${API_BASE_URL}/api/validate-api-key`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      provider: provider,
      api_key: apiKey
    }),
  });
  
  return handleApiError(response);
};