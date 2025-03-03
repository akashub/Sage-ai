"use client"

import { createContext, useContext, useState } from "react"

const AuthContext = createContext()

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null)
    const [authModalOpen, setAuthModalOpen] = useState(false)

    const openAuthModal = () => {
        setAuthModalOpen(true)
    }
    const closeAuthModal = () => setAuthModalOpen(false)

    return (
        <AuthContext.Provider
            value={{
                user,
                setUser,
                authModalOpen,
                openAuthModal,
                closeAuthModal,
            }}
        >
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

