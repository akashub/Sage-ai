describe('Auth Modal', () => {
    beforeEach(() => {
      // Visit the homepage before each test
      cy.visit('/')
    })
  
    it('should open auth modal when clicking Join Early Access button', () => {
      // Find and click the Join Early Access button
      cy.contains('Join Early Access').click()
      
      // Verify the auth modal is visible
      cy.contains('Sign in').should('be.visible')
      
      // Test form interaction
      cy.get('input[type="email"]').should('be.visible').type('test@example.com')
      cy.get('input[type="password"]').should('be.visible').type('password123')
      
      // Click the Sign In button
      cy.contains('button', 'Sign In').click()
    })
    
    
    it('should validate email and password fields', () => {
      // Open the auth modal
      cy.contains('Join Early Access').click()
      
      // Try to submit without entering data
      cy.contains('button', 'Sign In').click()
      
      // Enter invalid email
      cy.get('input[type="email"]').type('invalid-email')
      cy.contains('button', 'Sign In').click()
      
      // Enter valid email
      cy.get('input[type="email"]').clear().type('test@example.com')
      
      // Enter short password
      cy.get('input[type="password"]').type('short')
      cy.contains('button', 'Sign In').click()
      
      // Enter valid password
      cy.get('input[type="password"]').clear().type('password123')
      cy.contains('button', 'Sign In').click()
    })
    

  })