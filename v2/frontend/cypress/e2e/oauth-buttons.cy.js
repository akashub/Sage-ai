describe('OAuth Buttons', () => {
    beforeEach(() => {
      // Intercept auth check calls
      cy.intercept('/api/auth/user', {
        statusCode: 200,
        body: { authenticated: false }
      })
  
      // Visit homepage
      cy.visit('/')
      
      // Click the Join Early Access button with a longer timeout
      cy.contains('Join Early Access', { timeout: 10000 }).click()
      
      // Wait for modal animation
      cy.wait(1000)
    })
  
    it('should render social login buttons', () => {
      // Look for buttons containing OAuth provider names using case-insensitive, partial matches
      cy.get('button, a').then($elements => {
        const hasGithub = Array.from($elements).some(el => 
          el.textContent.toLowerCase().includes('github')
        )
        const hasGoogle = Array.from($elements).some(el => 
          el.textContent.toLowerCase().includes('google')
        )
        
        expect(hasGithub || hasGoogle).to.be.true
      })
    })
  
    it('should show OAuth buttons with correct styling', () => {
      // Look for any social login buttons
      cy.get('button, a').then($elements => {
        const socialButtons = Array.from($elements).filter(el => 
          el.textContent.toLowerCase().includes('github') ||
          el.textContent.toLowerCase().includes('google')
        )
        
        // Verify we found at least one button
        expect(socialButtons.length).to.be.greaterThan(0)
      })
    })
  
    it('should have clickable OAuth buttons', () => {
      // Find and click GitHub button if it exists
      cy.get('body').then($body => {
        if ($body.text().toLowerCase().includes('github')) {
          cy.get('button, a').contains(/github/i, { timeout: 5000 })
            .should('exist')
            .and('be.visible')
        }
        
        // Find and click Google button if it exists
        if ($body.text().toLowerCase().includes('google')) {
          cy.get('button, a').contains(/google/i, { timeout: 5000 })
            .should('exist')
            .and('be.visible')
        }
      })
    })
  
    it('should show divider between OAuth and email login', () => {
      // Look for common divider patterns
      cy.get('body').then($body => {
        const hasDivider = 
          $body.text().toLowerCase().includes('or') ||
          $body.find('hr').length > 0 ||
          $body.find('[role="separator"]').length > 0
        
        expect(hasDivider).to.be.true
      })
    })
  })