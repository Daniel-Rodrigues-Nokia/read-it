describe("Example Cypress Tests", () => {
  beforeEach(() => {
    // Runs before every test
    cy.visit("https://example.cypress.io");
  });

  it("1. Page should load and have correct title", () => {
    cy.title().should("include", "Cypress");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");

    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
    cy.get("h1").should("be.visible").and("contain.text", "Kitchen Sink");
  });

  it("2. Should type into an input field and submit a form", () => {
    cy.contains("Forms").click();
    cy.contains("Input").click();

    cy.get(".action-email")
      .type("test@example.com")
      .should("have.value", "test@example.com");

    cy.get(".action-form").submit();
  });

  it("3. Should intercept an API request and mock a response", () => {
    cy.intercept("GET", "**/comments/*", {
      body: { id: 1, name: "Mocked User" },
    }).as("getComment");

    cy.contains("Utilities").click();
    cy.contains("Network Requests").click();

    cy.get(".network-btn").click();
    cy.wait("@getComment").its("response.body").should("deep.equal", {
      id: 1,
      name: "Mocked User",
    });
  });

  it("4. Should check if a table has expected data", () => {
    cy.contains("Commands").click();
    cy.contains("Traversal").click();

    cy.get(".traversal-table tbody tr").should("have.length.at.least", 1);
    cy.get(".traversal-table tbody tr")
      .first()
      .within(() => {
        cy.get("td").eq(0).should("not.be.empty");
        cy.get("td").eq(1).should("not.be.empty");
      });
  });

  it("5. Should verify a button becomes visible after an action", () => {
    cy.contains("Commands").click();
    cy.contains("Actions").click();

    cy.get(".action-btn").should("not.be.disabled").click();
    cy.get("#action-canvas").should("be.visible");
  });

  // 6. Data-driven style test (multiple scenarios in one loop)
  [
    { desc: "Type a valid email", value: "hello@world.com" },
    { desc: "Type a numeric string", value: "12345" },
    { desc: "Type a random word", value: "CypressRocks" },
  ].forEach(({ desc, value }) => {
    it(`6. Input field test case - ${desc}`, () => {
      cy.contains("Forms").click();
      cy.contains("Input").click();

      cy.get(".action-email").clear().type(value).should("have.value", value);
    });
  });
});
