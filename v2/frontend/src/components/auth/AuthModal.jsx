"use client"

import { useState, useEffect } from "react"
import { X } from "lucide-react"
import { FeatureList } from "./FeatureList"

export const AuthModal = ({ isOpen, onClose, initialView = "signup" }) => {
  const [view, setView] = useState(initialView)
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = "hidden"
    } else {
      document.body.style.overflow = "unset"
    }
    return () => {
      document.body.style.overflow = "unset"
    }
  }, [isOpen])

  if (!mounted) return null
  if (!isOpen) return null

  const handleBackdropClick = (e) => {
    if (e.target === e.currentTarget) {
      onClose()
    }
  }

  return (
    <div className="fixed inset-0 z-50">
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-[#07031A]/90 backdrop-blur-sm transition-opacity"
        onClick={handleBackdropClick}
      />

      {/* Modal */}
      <div className="fixed left-1/2 top-1/2 w-full max-w-[800px] -translate-x-1/2 -translate-y-1/2 rounded-xl bg-[#07031A] p-6 shadow-2xl border border-white/5">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-2 bg-white/5 rounded-full px-3 py-1">
            <span className="text-[11px] uppercase tracking-wider text-white/60 font-medium">Built By</span>
            <img src="/logo.png" alt="SAGE.AI" className="h-3" />
          </div>
          <button onClick={onClose} className="rounded-full p-1.5 hover:bg-white/5 transition-colors">
            <X className="h-4 w-4 text-white/40" />
          </button>
        </div>

        {/* Content Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 md:gap-12">
          {/* Left Column */}
          <div>
            <FeatureList />
          </div>

          {/* Right Column */}
          <div>
            {view === "signup" ? (
              <div className="space-y-6">
                <div>
                  <h2 className="text-2xl font-semibold text-white mb-1">Create Account</h2>
                  <p className="text-sm text-white/60">Start creating amazing SQL queries with AI</p>
                </div>

                <form className="space-y-4" onSubmit={(e) => e.preventDefault()}>
                  <div>
                    <label className="block text-sm font-medium text-white/80 mb-1.5">Email</label>
                    <input
                      type="email"
                      className="w-full h-10 rounded-lg border border-white/10 bg-white/5 px-3 text-white placeholder-white/30 focus:border-[#5865F2] focus:ring-1 focus:ring-[#5865F2] focus:outline-none transition-colors"
                      required
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-white/80 mb-1.5">Password</label>
                    <input
                      type="password"
                      className="w-full h-10 rounded-lg border border-white/10 bg-white/5 px-3 text-white placeholder-white/30 focus:border-[#5865F2] focus:ring-1 focus:ring-[#5865F2] focus:outline-none transition-colors"
                      required
                    />
                  </div>

                  <button
                    type="submit"
                    className="w-full h-10 rounded-lg bg-white font-medium text-[#07031A] hover:bg-white/90 transition-colors mt-2"
                  >
                    Create Account
                  </button>
                </form>

                <p className="text-center text-sm text-white/40">
                  Already have an account?{" "}
                  <button
                    onClick={() => setView("signin")}
                    className="text-[#5865F2] hover:text-[#5865F2]/80 transition-colors font-medium"
                  >
                    Sign in
                  </button>
                </p>
              </div>
            ) : (
              <div className="space-y-6">
                <div>
                  <h2 className="text-2xl font-semibold text-white mb-1">Welcome Back</h2>
                  <p className="text-sm text-white/60">Sign in to continue your creative journey</p>
                </div>

                <form className="space-y-4" onSubmit={(e) => e.preventDefault()}>
                  <div>
                    <label className="block text-sm font-medium text-white/80 mb-1.5">Email</label>
                    <input
                      type="email"
                      className="w-full h-10 rounded-lg border border-white/10 bg-white/5 px-3 text-white placeholder-white/30 focus:border-[#5865F2] focus:ring-1 focus:ring-[#5865F2] focus:outline-none transition-colors"
                      required
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-white/80 mb-1.5">Password</label>
                    <input
                      type="password"
                      className="w-full h-10 rounded-lg border border-white/10 bg-white/5 px-3 text-white placeholder-white/30 focus:border-[#5865F2] focus:ring-1 focus:ring-[#5865F2] focus:outline-none transition-colors"
                      required
                    />
                  </div>

                  <button
                    type="submit"
                    className="w-full h-10 rounded-lg bg-white font-medium text-[#07031A] hover:bg-white/90 transition-colors mt-2"
                  >
                    Sign In
                  </button>
                </form>

                <p className="text-center text-sm text-white/40">
                  Don't have an account?{" "}
                  <button
                    onClick={() => setView("signup")}
                    className="text-[#5865F2] hover:text-[#5865F2]/80 transition-colors font-medium"
                  >
                    Create one
                  </button>
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

