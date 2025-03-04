const API_BASE = process.env.NODE_ENV === 'production' ? '' : 'http://localhost:8080';

export const uploadFile = async (file) => {
  const formData = new FormData();
  formData.append("file", file);
  
  try {
    const response = await fetch(`${API_BASE}/api/upload`, {
      method: "POST",
      body: formData,
    });
    
    if (!response.ok) {
      throw new Error(`Upload failed: ${response.status}`);
    }
    
    return response.json();
  } catch (error) {
    console.error("API Error (uploadFile):", error);
    throw error;
  }
};

export const queryData = async (query, csvPath) => {
  try {
    const response = await fetch(`${API_BASE}/api/query`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        query,
        csvPath,
      }),
    });
    
    if (!response.ok) {
      throw new Error(`Query failed: ${response.status}`);
    }
    
    return response.json();
  } catch (error) {
    console.error("API Error (queryData):", error);
    throw error;
  }
};