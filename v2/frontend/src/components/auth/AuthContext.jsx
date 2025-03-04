// "use client"

// import { createContext, useContext, useState, useEffect } from "react"

// const AuthContext = createContext()

// export const AuthProvider = ({ children }) => {
//     const [user, setUser] = useState(null)
//     const [authModalOpen, setAuthModalOpen] = useState(false)
//     const [authMode, setAuthMode] = useState('signin') // 'signin', 'signup', 'forgot'
//     const [loading, setLoading] = useState(true)

//     useEffect(() => {
//         // Check if user is already logged in (e.g., token in localStorage)
//         const checkAuth = async () => {
//             try {
//                 const token = localStorage.getItem('authToken')
//                 if (token) {
//                     // Here you would validate the token with your backend
//                     // For now, just simulate a user
//                     setUser({
//                         email: localStorage.getItem('userEmail') || 'user@example.com',
//                         id: '123'
//                     })
//                 }
//             } catch (error) {
//                 console.error('Auth check failed:', error)
//                 localStorage.removeItem('authToken')
//             } finally {
//                 setLoading(false)
//             }
//         }

//         checkAuth()
//     }, [])

//     const openAuthModal = (mode = 'signin') => {
//         setAuthMode(mode)
//         setAuthModalOpen(true)
//     }

//     const closeAuthModal = () => {
//         setAuthModalOpen(false)
//         // Reset to signin mode after closing for next time
//         setTimeout(() => setAuthMode('signin'), 300)
//     }

//     const switchAuthMode = (mode) => {
//         setAuthMode(mode)
//     }

//     const signIn = async (email, password) => {
//         try {
//             // Here you would call your API to sign in
//             console.log('Signing in with:', email, password)
            
//             // Simulate successful login
//             localStorage.setItem('authToken', 'sample-token-123')
//             localStorage.setItem('userEmail', email)
            
//             setUser({
//                 email,
//                 id: '123'
//             })
            
//             closeAuthModal()
//             return { success: true }
//         } catch (error) {
//             console.error('Sign in failed:', error)
//             return { 
//                 success: false, 
//                 error: error.message || 'Sign in failed. Please try again.'
//             }
//         }
//     }

//     const signUp = async (email, password) => {
//         try {
//             // Here you would call your API to register
//             console.log('Signing up with:', email, password)
            
//             // Simulate successful registration
//             localStorage.setItem('authToken', 'sample-token-123')
//             localStorage.setItem('userEmail', email)
            
//             setUser({
//                 email,
//                 id: '123'
//             })
            
//             closeAuthModal()
//             return { success: true }
//         } catch (error) {
//             console.error('Sign up failed:', error)
//             return { 
//                 success: false, 
//                 error: error.message || 'Registration failed. Please try again.'
//             }
//         }
//     }

//     const signOut = () => {
//         localStorage.removeItem('authToken')
//         localStorage.removeItem('userEmail')
//         setUser(null)
//     }

//     const resetPassword = async (email) => {
//         try {
//             // Here you would call your API to send password reset email
//             console.log('Password reset for:', email)
            
//             // Simulate successful password reset request
//             return { 
//                 success: true,
//                 message: 'Password reset link sent! Check your email.'
//             }
//         } catch (error) {
//             console.error('Password reset failed:', error)
//             return { 
//                 success: false, 
//                 error: error.message || 'Password reset failed. Please try again.'
//             }
//         }
//     }

//     const value = {
//         user,
//         setUser,
//         loading,
//         authModalOpen,
//         authMode,
//         openAuthModal,
//         closeAuthModal,
//         switchAuthMode,
//         signIn,
//         signUp,
//         signOut,
//         resetPassword
//     }

//     return (
//         <AuthContext.Provider value={value}>
//             {children}
//         </AuthContext.Provider>
//     )
// }

// export const useAuth = () => {
//     const context = useContext(AuthContext)
//     if (!context) {
//         throw new Error("useAuth must be used within an AuthProvider")
//     }
//     return context
// }

// frontend/src/components/auth/AuthContext.jsx
// components/auth/AuthContext.jsx
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
  const navigate = useNavigate();

  // Check if user is already logged in on initial load
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
  const openAuthModal = () => {
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
//   const signIn = async (email, password) => {
//     setError(null);
//     try {
//       const response = await fetch("/api/auth/signin", {
//         method: "POST",
//         headers: {
//           "Content-Type": "application/json",
//         },
//         credentials: "include",
//         body: JSON.stringify({ email, password }),
//       });

//       if (!response.ok) {
//         const errorData = await response.json();
//         throw new Error(errorData.message || "Sign in failed");
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
const signIn = async (email, password) => {
    setError(null);
    try {
      // First check if the API is available
      try {
        const checkResponse = await fetch("/api/health", { method: "GET" });
        if (!checkResponse.ok) {
          throw new Error("API server not available");
        }
      } catch (err) {
        console.error("API server check failed:", err);
        throw new Error("Backend API service is not running. Please start the server.");
      }
      
      const response = await fetch("/api/auth/signin", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ email, password }),
      });
  
      // Handle non-JSON responses
      const contentType = response.headers.get("content-type");
      if (!contentType || !contentType.includes("application/json")) {
        console.error("Received non-JSON response:", await response.text());
        throw new Error("Invalid response from server");
      }
  
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Sign in failed");
      }
  
      const data = await response.json();
      setUser(data.user);
      closeAuthModal();
      navigate("/chat");
      return data;
    } catch (err) {
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
      setUser(data.user);
      closeAuthModal();
      navigate("/chat"); // Redirect to chat screen
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
      const response = await fetch(`/api/auth/oauth/${provider}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ code, redirect_uri: redirectUri }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "OAuth sign in failed");
      }

      const data = await response.json();
      setUser(data.user);
      closeAuthModal();
      navigate("/chat"); // Redirect to chat screen
      return data;
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  // Get OAuth URL function
  const getOAuthUrl = async (provider, redirectUri) => {
    try {
      const response = await fetch(`/api/auth/oauth/url/${provider}?redirect_uri=${encodeURIComponent(redirectUri)}`, {
        method: "GET",
      });

      if (!response.ok) {
        throw new Error("Failed to get OAuth URL");
      }

      const data = await response.json();
      return data.url;
    } catch (err) {
      setError(err.message);
      throw err;
    }
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
        switchAuthMode
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