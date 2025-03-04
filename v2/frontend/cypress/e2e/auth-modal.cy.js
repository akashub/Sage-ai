describe("Auth Modal", () => {
  beforeEach(() => {
    // Visit the homepage
    cy.visit("/")
  })

  it("should open auth modal when clicking Join Early Access button", () => {
    // Find and click the Join Early Access button
    cy.contains("Join Early Access").click()

    // Verify the auth modal is visible
    cy.contains("Sign in").should("be.visible")

    // Test form interaction
    cy.get('input[type="email"]').type("test@example.com")
    cy.get('input[type="password"]').type("password123")

    // Click the Sign In button
    cy.contains("button", "Sign In").click()
  })

  it("should switch between sign in and sign up modes", () => {
    // Open the auth modal
    cy.contains("Join Early Access").click()

    // Switch to sign up mode
    cy.contains("Create one").click()

    // Verify we're in sign up mode
    cy.contains("Create account").should("be.visible")

    // Switch back to sign in
    cy.contains("Sign in").click()

    // Verify we're back in sign in mode
    cy.contains("Welcome Back").should("be.visible")
  })
})

