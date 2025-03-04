describe("Chat Interface", () => {
  beforeEach(() => {
    // Visit the chat page directly
    cy.visit("/chat")
  })

  it("should display the chat interface", () => {
    // Verify the chat interface is visible
    cy.contains("Welcome to Sage AI").should("be.visible")
  })

  it("should send a message and receive a response", () => {
    // Type a message
    cy.get('input[placeholder="Type your message here..."]').type("Show me a SQL query for user data")

    // Send the message
    cy.contains("button", "Send").click()

    // Verify the message was sent
    cy.contains("Show me a SQL query for user data").should("be.visible")

    // Wait for and verify the response
    cy.contains("Here is your SQL query").should("be.visible", { timeout: 10000 })
  })

  it("should switch between chat history items", () => {
    // Click on a chat history item
    cy.contains("SQL Query Optimization").click()

    // Verify the chat title changes
    cy.contains("Chat 1").should("be.visible")

    // Click on another chat history item
    cy.contains("Database Schema Design").click()

    // Verify the chat title changes
    cy.contains("Chat 2").should("be.visible")
  })
})

