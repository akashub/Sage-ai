// // import React, { createContext, useState, useEffect, useContext } from "react";
// // import { useNavigate } from "react-router-dom";

// // // Create Auth Context
// // export const AuthContext = createContext();

// // // Auth Provider Component
// // export const AuthProvider = ({ children }) => {
// //   const [user, setUser] = useState(null);
// //   const [loading, setLoading] = useState(true);
// //   const [error, setError] = useState(null);
// //   const [authModalOpen, setAuthModalOpen] = useState(false);
// //   const [authMode, setAuthMode] = useState("signin");
// //   const navigate = useNavigate();

// //   // Check if user is already logged in on initial load
// //   useEffect(() => {
// //     const checkAuthStatus = async () => {
// //       try {
// //         const response = await fetch("/api/auth/user", {
// //           method: "GET",
// //           credentials: "include",
// //         });

// //         if (response.ok) {
// //           const userData = await response.json();
// //           setUser(userData);
// //         }
// //       } catch (err) {
// //         console.error("Authentication check failed:", err);
// //       } finally {
// //         setLoading(false);
// //       }
// //     };

// //     checkAuthStatus();
// //   }, []);

// //   // Open auth modal
// //   const openAuthModal = () => {
// //     setAuthModalOpen(true);
// //   };

// //   // Close auth modal
// //   const closeAuthModal = () => {
// //     setAuthModalOpen(false);
// //   };

// //   // Switch between auth modes (signin, signup, forgot)
// //   const switchAuthMode = (mode) => {
// //     setAuthMode(mode);
// //   };

// //   // Sign in function
// // //   const signIn = async (email, password) => {
// // //     setError(null);
// // //     try {
// // //       const response = await fetch("/api/auth/signin", {
// // //         method: "POST",
// // //         headers: {
// // //           "Content-Type": "application/json",
// // //         },
// // //         credentials: "include",
// // //         body: JSON.stringify({ email, password }),
// // //       });

// // //       if (!response.ok) {
// // //         const errorData = await response.json();
// // //         throw new Error(errorData.message || "Sign in failed");
// // //       }

// // //       const data = await response.json();
// // //       setUser(data.user);
// // //       closeAuthModal();
// // //       navigate("/chat"); // Redirect to chat screen
// // //       return data;
// // //     } catch (err) {
// // //       setError(err.message);
// // //       throw err;
// // //     }
// // //   };
// // // const signIn = async (email, password) => {
// // //     setError(null);
// // //     try {
// // //       // First check if the API is available
// // //       try {
// // //         const checkResponse = await fetch("/api/health", { method: "GET" });
// // //         if (!checkResponse.ok) {
// // //           throw new Error("API server not available");
// // //         }
// // //       } catch (err) {
// // //         console.error("API server check failed:", err);
// // //         throw new Error("Backend API service is not running. Please start the server.");
// // //       }
      
// // //       const response = await fetch("/api/auth/signin", {
// // //         method: "POST",
// // //         headers: {
// // //           "Content-Type": "application/json",
// // //         },
// // //         credentials: "include",
// // //         body: JSON.stringify({ email, password }),
// // //       });
  
// // //       // Handle non-JSON responses
// // //       const contentType = response.headers.get("content-type");
// // //       if (!contentType || !contentType.includes("application/json")) {
// // //         console.error("Received non-JSON response:", await response.text());
// // //         throw new Error("Invalid response from server");
// // //       }
  
// // //       if (!response.ok) {
// // //         const errorData = await response.json();
// // //         throw new Error(errorData.message || "Sign in failed");
// // //       }
  
// // //       const data = await response.json();
// // //       setUser(data.user);
// // //       closeAuthModal();
// // //       navigate("/chat");
// // //       return data;
// // //     } catch (err) {
// // //       setError(err.message);
// // //       throw err;
// // //     }
// // //   };
// // const signIn = async (email, password) => {
// //   setError(null);
// //   try {
// //     console.log("Making sign-in request to:", "/api/auth/signin");
    
// //     const response = await fetch("/api/auth/signin", {
// //       method: "POST",
// //       headers: {
// //         "Content-Type": "application/json",
// //       },
// //       credentials: "include",
// //       body: JSON.stringify({ email, password }),
// //     });
    
// //     console.log("Received response:", response.status, response.statusText);
    
// //     // Check if response is JSON
// //     const contentType = response.headers.get("content-type");
// //     if (!contentType || !contentType.includes("application/json")) {
// //       // Log the HTML response for debugging
// //       const textResponse = await response.text();
// //       console.error("Received non-JSON response:", textResponse);
// //       throw new Error("Server returned an invalid response format");
// //     }

// //     if (!response.ok) {
// //       const errorData = await response.json();
// //       throw new Error(errorData.message || "Sign in failed");
// //     }

// //     const data = await response.json();
// //     setUser(data.user);
// //     closeAuthModal();
// //     navigate("/chat");
// //     return data;
// //   } catch (err) {
// //     console.error("Login error:", err);
// //     setError(err.message);
// //     throw err;
// //   }
// // };

// //   // Sign up function
// //   const signUp = async (email, password, name) => {
// //     setError(null);
// //     try {
// //       const response = await fetch("/api/auth/signup", {
// //         method: "POST",
// //         headers: {
// //           "Content-Type": "application/json",
// //         },
// //         credentials: "include",
// //         body: JSON.stringify({ email, password, name }),
// //       });

// //       if (!response.ok) {
// //         const errorData = await response.json();
// //         throw new Error(errorData.message || "Sign up failed");
// //       }

// //       const data = await response.json();
// //       setUser(data.user);
// //       closeAuthModal();
// //       navigate("/chat"); // Redirect to chat screen
// //       return data;
// //     } catch (err) {
// //       setError(err.message);
// //       throw err;
// //     }
// //   };

// //   // OAuth sign in function
// //   const oauthSignIn = async (provider, code, redirectUri) => {
// //     setError(null);
// //     try {
// //       const response = await fetch(`/api/auth/oauth/${provider}`, {
// //         method: "POST",
// //         headers: {
// //           "Content-Type": "application/json",
// //         },
// //         credentials: "include",
// //         body: JSON.stringify({ code, redirect_uri: redirectUri }),
// //       });

// //       if (!response.ok) {
// //         const errorData = await response.json();
// //         throw new Error(errorData.message || "OAuth sign in failed");
// //       }

// //       const data = await response.json();
// //       setUser(data.user);
// //       closeAuthModal();
// //       navigate("/chat"); // Redirect to chat screen
// //       return data;
// //     } catch (err) {
// //       setError(err.message);
// //       throw err;
// //     }
// //   };

// //   // Get OAuth URL function
// //   const getOAuthUrl = async (provider, redirectUri) => {
// //     try {
// //       const response = await fetch(`/api/auth/oauth/url/${provider}?redirect_uri=${encodeURIComponent(redirectUri)}`, {
// //         method: "GET",
// //       });

// //       if (!response.ok) {
// //         throw new Error("Failed to get OAuth URL");
// //       }

// //       const data = await response.json();
// //       return data.url;
// //     } catch (err) {
// //       setError(err.message);
// //       throw err;
// //     }
// //   };

// //   // Sign out function
// //   const signOut = async () => {
// //     try {
// //       await fetch("/api/auth/signout", {
// //         method: "POST",
// //         credentials: "include",
// //       });
// //       setUser(null);
// //       navigate("/");
// //     } catch (err) {
// //       console.error("Sign out failed:", err);
// //     }
// //   };

// //   return (
// //     <AuthContext.Provider
// //       value={{
// //         user,
// //         loading,
// //         error,
// //         signIn,
// //         signUp,
// //         signOut,
// //         oauthSignIn,
// //         getOAuthUrl,
// //         openAuthModal,
// //         closeAuthModal,
// //         authModalOpen,
// //         authMode,
// //         switchAuthMode
// //       }}
// //     >
// //       {children}
// //     </AuthContext.Provider>
// //   );
// // };

// // // Custom hook to use the auth context
// // export const useAuth = () => {
// //   const context = useContext(AuthContext);
// //   if (!context) {
// //     throw new Error("useAuth must be used within an AuthProvider");
// //   }
// //   return context;
// // };

// import React, { createContext, useState, useEffect, useContext } from "react";
// import { useNavigate } from "react-router-dom";

// // Create Auth Context
// export const AuthContext = createContext();

// // Auth Provider Component
// export const AuthProvider = ({ children }) => {
//   const [user, setUser] = useState(null);
//   const [loading, setLoading] = useState(true);
//   const [error, setError] = useState(null);
//   const [authModalOpen, setAuthModalOpen] = useState(false);
//   const [authMode, setAuthMode] = useState("signin");
//   const navigate = useNavigate();

//   // Check if user is already logged in on initial load
//   useEffect(() => {
//     const checkAuthStatus = async () => {
//       try {
//         const response = await fetch("/api/auth/user", {
//           method: "GET",
//           credentials: "include",
//         });

//         if (response.ok) {
//           const userData = await response.json();
//           setUser(userData);
//         }
//       } catch (err) {
//         console.error("Authentication check failed:", err);
//       } finally {
//         setLoading(false);
//       }
//     };

//     checkAuthStatus();
//   }, []);

//   // Open auth modal
//   const openAuthModal = (mode = "signin") => {
//     setAuthMode(mode);
//     setAuthModalOpen(true);
//   };

//   // Close auth modal
//   const closeAuthModal = () => {
//     setAuthModalOpen(false);
//   };

//   // Switch between auth modes (signin, signup, forgot)
//   const switchAuthMode = (mode) => {
//     setAuthMode(mode);
//   };

//   // Sign in function
//   const signIn = async (email, password) => {
//     setError(null);
//     try {
//       console.log("Making sign-in request to:", "/api/auth/signin");
      
//       const response = await fetch("/api/auth/signin", {
//         method: "POST",
//         headers: {
//           "Content-Type": "application/json",
//         },
//         credentials: "include",
//         body: JSON.stringify({ email, password }),
//       });
      
//       console.log("Received response:", response.status, response.statusText);
      
//       // Check if response is JSON
//       const contentType = response.headers.get("content-type");
//       if (!contentType || !contentType.includes("application/json")) {
//         // Log the HTML response for debugging
//         const textResponse = await response.text();
//         console.error("Received non-JSON response:", textResponse);
//         throw new Error("Server returned an invalid response format");
//       }

//       if (!response.ok) {
//         const errorData = await response.json();
//         throw new Error(errorData.message || "Sign in failed");
//       }

//       const data = await response.json();
//       setUser(data.user);
//       closeAuthModal();
//       navigate("/chat");
//       return data;
//     } catch (err) {
//       console.error("Login error:", err);
//       setError(err.message);
//       throw err;
//     }
//   };

//   // Sign up function
//   const signUp = async (email, password, name) => {
//     setError(null);
//     try {
//       const response = await fetch("/api/auth/signup", {
//         method: "POST",
//         headers: {
//           "Content-Type": "application/json",
//         },
//         credentials: "include",
//         body: JSON.stringify({ email, password, name }),
//       });

//       if (!response.ok) {
//         const errorData = await response.json();
//         throw new Error(errorData.message || "Sign up failed");
//       }

//       const data = await response.json();
//       setUser(data.user);
//       closeAuthModal();
//       navigate("/chat"); // Redirect to chat screen
//       return data;
//     } catch (err) {
//       setError(err.message);
//       throw err;
//     }
//   };

//   // OAuth sign in function
//   const oauthSignIn = async (provider, code, redirectUri) => {
//     setError(null);
//     try {
//       console.log(`Making OAuth sign-in request for provider: ${provider} with code length: ${code.length}`);
      
//       const response = await fetch(`/api/auth/oauth/${provider}`, {
//         method: "POST",
//         headers: {
//           "Content-Type": "application/json",
//           "Accept": "application/json"
//         },
//         credentials: "include",
//         body: JSON.stringify({ code, redirect_uri: redirectUri }),
//       });

//       console.log(`OAuth response status: ${response.status}`);
      
//       // Check content type
//       const contentType = response.headers.get("content-type");
//       if (!contentType || !contentType.includes("application/json")) {
//         const textResponse = await response.text();
//         console.error("OAuth non-JSON response:", textResponse);
//         throw new Error("Server returned an invalid response format");
//       }

//       if (!response.ok) {
//         const errorData = await response.json();
//         throw new Error(errorData.message || "OAuth sign in failed");
//       }

//       const data = await response.json();
//       setUser(data.user);
//       closeAuthModal();
//       navigate("/chat");
//       return data;
//     } catch (err) {
//       console.error("OAuth sign-in error:", err);
//       setError(err.message);
//       throw err;
//     }
//   };

//   // Get OAuth URL function - improved for better error handling
//   const getOAuthUrl = async (provider, redirectUri) => {
//     try {
//       const encodedRedirectUri = encodeURIComponent(redirectUri);
//       const url = `/api/auth/oauth/url/${provider}?redirect_uri=${encodedRedirectUri}`;
//       console.log("Making OAuth URL request to:", url);
      
//       const response = await fetch(url, {
//         method: "GET",
//         headers: {
//           "Accept": "application/json",
//         },
//       });

//       console.log("OAuth URL response status:", response.status);

//       if (!response.ok) {
//         console.error(`OAuth URL error: ${response.status}`);
//         const errorText = await response.text();
//         console.error("Error response text:", errorText);
//         throw new Error(`Failed to get OAuth URL: ${response.status}`);
//       }

//       const contentType = response.headers.get("content-type");
//       if (!contentType || !contentType.includes("application/json")) {
//         const text = await response.text();
//         console.error("Non-JSON OAuth URL response:", text);
//         throw new Error("Invalid response format from OAuth URL endpoint");
//       }

//       const data = await response.json();
//       console.log("OAuth URL response data:", data);
      
//       if (!data.url) {
//         throw new Error("OAuth URL not provided in response");
//       }
      
//       return data.url;
//     } catch (err) {
//       console.error("getOAuthUrl error:", err);
//       setError("Failed to get authorization URL. Please try again later.");
//       throw err;
//     }
//   };

//   // Sign out function
//   const signOut = async () => {
//     try {
//       await fetch("/api/auth/signout", {
//         method: "POST",
//         credentials: "include",
//       });
//       setUser(null);
//       navigate("/");
//     } catch (err) {
//       console.error("Sign out failed:", err);
//     }
//   };

//   return (
//     <AuthContext.Provider
//       value={{
//         user,
//         loading,
//         error,
//         signIn,
//         signUp,
//         signOut,
//         oauthSignIn,
//         getOAuthUrl,
//         openAuthModal,
//         closeAuthModal,
//         authModalOpen,
//         authMode,
//         switchAuthMode
//       }}
//     >
//       {children}
//     </AuthContext.Provider>
//   );
// };

// // Custom hook to use the auth context
// export const useAuth = () => {
//   const context = useContext(AuthContext);
//   if (!context) {
//     throw new Error("useAuth must be used within an AuthProvider");
//   }
//   return context;
// };

import React, { createContext, useState, useEffect, useContext } from "react";
import { useNavigate } from "react-router-dom";

// Create Auth Context
export const AuthContext = createContext();

// Auth Provider Component
export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [authModalOpen, setAuthModalOpen] = useState(false);
  const [authMode, setAuthMode] = useState("signin");
  const [welcomeDialogOpen, setWelcomeDialogOpen] = useState(false);
  const navigate = useNavigate();

  // Check if user is already logged in on initial load - without showing welcome dialog
  useEffect(() => {
    const checkAuthStatus = async () => {
      try {
        const response = await fetch("/api/auth/user", {
          method: "GET",
          credentials: "include",
        });

        if (response.ok) {
          const userData = await response.json();
          setUser(userData);
          // Don't show welcome dialog on initial page load
          // We'll only show it on fresh logins
        }
      } catch (err) {
        console.error("Authentication check failed:", err);
      } finally {
        setLoading(false);
      }
    };

    checkAuthStatus();
  }, []);

  // Open auth modal
  const openAuthModal = (mode = "signin") => {
    setAuthMode(mode);
    setAuthModalOpen(true);
  };

  // Close auth modal
  const closeAuthModal = () => {
    setAuthModalOpen(false);
  };

  // Switch between auth modes (signin, signup, forgot)
  const switchAuthMode = (mode) => {
    setAuthMode(mode);
  };

  // Sign in function
  const signIn = async (email, password) => {
    setError(null);
    try {
      console.log("Making sign-in request to:", "/api/auth/signin");
      
      const response = await fetch("/api/auth/signin", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ email, password }),
      });
      
      console.log("Received response:", response.status, response.statusText);
      
      // Check if response is JSON
      const contentType = response.headers.get("content-type");
      if (!contentType || !contentType.includes("application/json")) {
        // Log the HTML response for debugging
        const textResponse = await response.text();
        console.error("Received non-JSON response:", textResponse);
        throw new Error("Server returned an invalid response format");
      }

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Sign in failed");
      }

      const data = await response.json();
      
      // Before setting user, check if we should customize the user's name
      const shouldShowWelcomeDialog = 
        data.user.name === "Aakash Singh" || 
        data.user.name === "GitHub User" || 
        data.user.name === "Google User" || 
        data.user.name === "OAuth User" ||
        !data.user.name;
      
      // Set the user data
      setUser(data.user);
      closeAuthModal();
      
      // Show welcome dialog only on fresh login if needed
      if (shouldShowWelcomeDialog) {
        setWelcomeDialogOpen(true);
      }
      
      navigate("/chat");
      return data;
    } catch (err) {
      console.error("Login error:", err);
      setError(err.message);
      throw err;
    }
  };

  // Sign up function
  const signUp = async (email, password, name) => {
    setError(null);
    try {
      const response = await fetch("/api/auth/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ email, password, name }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Sign up failed");
      }

      const data = await response.json();
      
      // Set user data from the response
      setUser(data.user);
      closeAuthModal();
      
      // Always show the welcome dialog for new sign-ups
      setWelcomeDialogOpen(true);
      
      navigate("/chat");
      return data;
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  // OAuth sign in function
  const oauthSignIn = async (provider, code, redirectUri) => {
    setError(null);
    try {
      console.log(`Making OAuth sign-in request for provider: ${provider} with code length: ${code.length}`);
      
      const response = await fetch(`/api/auth/oauth/${provider}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Accept": "application/json"
        },
        credentials: "include",
        body: JSON.stringify({ code, redirect_uri: redirectUri }),
      });

      console.log(`OAuth response status: ${response.status}`);
      
      // Check content type
      const contentType = response.headers.get("content-type");
      if (!contentType || !contentType.includes("application/json")) {
        const textResponse = await response.text();
        console.error("OAuth non-JSON response:", textResponse);
        throw new Error("Server returned an invalid response format");
      }

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "OAuth sign in failed");
      }

      const data = await response.json();
      
      // Check if we need to show the welcome dialog (for OAuth users with generic names)
      const shouldShowWelcomeDialog = 
        data.user.name === "Test User" || 
        data.user.name === "GitHub User" || 
        data.user.name === "Google User" || 
        data.user.name === "OAuth User" || 
        !data.user.name;
      
      setUser(data.user);
      closeAuthModal();
      
      // Show welcome dialog only on fresh OAuth login if needed
      if (shouldShowWelcomeDialog) {
        setWelcomeDialogOpen(true);
      }
      
      navigate("/chat");
      return data;
    } catch (err) {
      console.error("OAuth sign-in error:", err);
      setError(err.message);
      throw err;
    }
  };

  // Update user name
  const updateUserName = async (name) => {
    try {
      // Make a real API call to update the user's name
      // Use the existing endpoint from your backend (/api/profile with PUT method)
      const response = await fetch("/api/profile", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ name }),
      });

      if (!response.ok) {
        throw new Error("Failed to update profile");
      }

      const updatedUser = await response.json();
      
      // Update the local user state with the response data
      setUser(prevUser => ({
        ...prevUser,
        name: updatedUser.name
      }));
      
      return updatedUser;
    } catch (err) {
      console.error("Error updating user name:", err);
      throw err;
    }
  };

  // Get OAuth URL function
  const getOAuthUrl = async (provider, redirectUri) => {
    try {
      const encodedRedirectUri = encodeURIComponent(redirectUri);
      const url = `/api/auth/oauth/url/${provider}?redirect_uri=${encodedRedirectUri}`;
      console.log("Making OAuth URL request to:", url);
      
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Accept": "application/json",
        },
      });

      console.log("OAuth URL response status:", response.status);

      if (!response.ok) {
        console.error(`OAuth URL error: ${response.status}`);
        const errorText = await response.text();
        console.error("Error response text:", errorText);
        throw new Error(`Failed to get OAuth URL: ${response.status}`);
      }

      const contentType = response.headers.get("content-type");
      if (!contentType || !contentType.includes("application/json")) {
        const text = await response.text();
        console.error("Non-JSON OAuth URL response:", text);
        throw new Error("Invalid response format from OAuth URL endpoint");
      }

      const data = await response.json();
      console.log("OAuth URL response data:", data);
      
      if (!data.url) {
        throw new Error("OAuth URL not provided in response");
      }
      
      return data.url;
    } catch (err) {
      console.error("getOAuthUrl error:", err);
      setError("Failed to get authorization URL. Please try again later.");
      throw err;
    }
  };

  // Close welcome dialog
  const closeWelcomeDialog = () => {
    setWelcomeDialogOpen(false);
  };

  // Sign out function
  const signOut = async () => {
    try {
      await fetch("/api/auth/signout", {
        method: "POST",
        credentials: "include",
      });
      setUser(null);
      navigate("/");
    } catch (err) {
      console.error("Sign out failed:", err);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        error,
        signIn,
        signUp,
        signOut,
        oauthSignIn,
        getOAuthUrl,
        openAuthModal,
        closeAuthModal,
        authModalOpen,
        authMode,
        switchAuthMode,
        updateUserName,
        welcomeDialogOpen,
        closeWelcomeDialog
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

// Custom hook to use the auth context
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};