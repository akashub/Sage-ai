describe('Sign In Form', () => {
    beforeEach(() => {
      // Visit the homepage
      cy.visit('/')
      
      // Open the auth modal - using a more flexible selector
      cy.contains('Join Early Access', { timeout: 10000 }).click()
      
      // Wait for modal to be fully visible
      cy.wait(1000)
    })
    
    it('should render email and password fields', () => {
      // Check if form elements exist without requiring them to be visible
      cy.get('input[type="email"]').should('exist')
      cy.get('input[type="password"]').should('exist')
    })
    
    it('should allow entering email and password', () => {
      // Enter email
      cy.get('input[type="email"]').type('test@example.com')
      
      // Enter password
      cy.get('input[type="password"]').type('password123')
      
      // Verify the values were entered
      cy.get('input[type="email"]').should('have.value', 'test@example.com')
      cy.get('input[type="password"]').should('have.value', 'password123')
    })
    
    it('should have a sign in button', () => {
      // Look for any button that might be the sign in button
      // Using a case-insensitive, partial match
      cy.get('button').contains(/sign.?in/i, { matchCase: false }).should('exist')
    })
    
  })