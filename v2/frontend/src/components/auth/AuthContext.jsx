"use client"

import { createContext, useContext, useState, useEffect } from "react"

const AuthContext = createContext()

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null)
    const [authModalOpen, setAuthModalOpen] = useState(false)
    const [authMode, setAuthMode] = useState('signin') // 'signin', 'signup', 'forgot'
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        // Check if user is already logged in (e.g., token in localStorage)
        const checkAuth = async () => {
            try {
                const token = localStorage.getItem('authToken')
                if (token) {
                    // Here you would validate the token with your backend
                    // For now, just simulate a user
                    setUser({
                        email: localStorage.getItem('userEmail') || 'user@example.com',
                        id: '123'
                    })
                }
            } catch (error) {
                console.error('Auth check failed:', error)
                localStorage.removeItem('authToken')
            } finally {
                setLoading(false)
            }
        }

        checkAuth()
    }, [])

    const openAuthModal = (mode = 'signin') => {
        setAuthMode(mode)
        setAuthModalOpen(true)
    }

    const closeAuthModal = () => {
        setAuthModalOpen(false)
        // Reset to signin mode after closing for next time
        setTimeout(() => setAuthMode('signin'), 300)
    }

    const switchAuthMode = (mode) => {
        setAuthMode(mode)
    }

    const signIn = async (email, password) => {
        try {
            // Here you would call your API to sign in
            console.log('Signing in with:', email, password)
            
            // Simulate successful login
            localStorage.setItem('authToken', 'sample-token-123')
            localStorage.setItem('userEmail', email)
            
            setUser({
                email,
                id: '123'
            })
            
            closeAuthModal()
            return { success: true }
        } catch (error) {
            console.error('Sign in failed:', error)
            return { 
                success: false, 
                error: error.message || 'Sign in failed. Please try again.'
            }
        }
    }

    const signUp = async (email, password) => {
        try {
            // Here you would call your API to register
            console.log('Signing up with:', email, password)
            
            // Simulate successful registration
            localStorage.setItem('authToken', 'sample-token-123')
            localStorage.setItem('userEmail', email)
            
            setUser({
                email,
                id: '123'
            })
            
            closeAuthModal()
            return { success: true }
        } catch (error) {
            console.error('Sign up failed:', error)
            return { 
                success: false, 
                error: error.message || 'Registration failed. Please try again.'
            }
        }
    }

    const signOut = () => {
        localStorage.removeItem('authToken')
        localStorage.removeItem('userEmail')
        setUser(null)
    }

    const resetPassword = async (email) => {
        try {
            // Here you would call your API to send password reset email
            console.log('Password reset for:', email)
            
            // Simulate successful password reset request
            return { 
                success: true,
                message: 'Password reset link sent! Check your email.'
            }
        } catch (error) {
            console.error('Password reset failed:', error)
            return { 
                success: false, 
                error: error.message || 'Password reset failed. Please try again.'
            }
        }
    }

    const value = {
        user,
        setUser,
        loading,
        authModalOpen,
        authMode,
        openAuthModal,
        closeAuthModal,
        switchAuthMode,
        signIn,
        signUp,
        signOut,
        resetPassword
    }

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    )
}

export const useAuth = () => {
    const context = useContext(AuthContext)
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider")
    }
    return context
}