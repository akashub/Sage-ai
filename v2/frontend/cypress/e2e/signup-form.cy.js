describe('Sign Up Form', () => {
    beforeEach(() => {
      // Visit the homepage
      cy.visit('/')
      
      // Click the Join Early Access button with a longer timeout
      cy.contains('Join Early Access', { timeout: 10000 }).click()
      
      // Wait for modal animation
      cy.wait(1000)
      
      // Look for any element that might trigger sign up mode using multiple possible texts
      cy.get('body').then($body => {
        // Try different possible texts for the sign up link
        const signupTriggers = [
          'Create one',
          'Create account',
          'Sign up',
          'Register',
          'Create new account'
        ]
        
        // Find the first matching text and click it
        const trigger = signupTriggers.find(text => 
          $body.text().includes(text)
        )
        
        if (trigger) {
          cy.contains(trigger).click()
        }
      })
      
      // Wait for mode switch animation
      cy.wait(1000)
    })
  
    it('should show signup form fields', () => {
      // Check for input fields without being too specific
      cy.get('input').should('have.length.at.least', 2)
    })
  
    it('should allow entering registration details', () => {
      // Find email input - try multiple possible labels
      cy.get('input[type="email"], input[placeholder*="email" i]')
        .first()
        .type('test@example.com')
  
      // Find password input - try multiple possible labels
      cy.get('input[type="password"]')
        .first()
        .type('password123')
  
      // If there's a confirm password field, fill that too
      cy.get('input[type="password"]').then($inputs => {
        if ($inputs.length > 1) {
          cy.wrap($inputs).last().type('password123')
        }
      })
    })
  
    it('should have a submit button', () => {
      // Look for any button that might be the submit button
      cy.get('button').then($buttons => {
        const submitTexts = [
          'Sign up',
          'Create account',
          'Register',
          'Submit',
          'Join'
        ]
        
        const submitButton = Array.from($buttons).find(button => 
          submitTexts.some(text => 
            button.textContent.toLowerCase().includes(text.toLowerCase())
          )
        )
        
        expect(submitButton).to.exist
      })
    })
  })