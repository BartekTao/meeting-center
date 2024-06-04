describe('Check each url', () => {
  beforeEach(() => {
    // 訪問應用的根URL
    cy.visit('http://localhost:8888/#/');
  });

  it('should check the presence of navigation links', () => {
    // 檢查導航鏈接
    cy.get('.main-nav').within(() => {
      cy.get('.nav-link').eq(0).should('have.text', '空間查詢');
      cy.get('.nav-link').eq(1).should('have.text', '已預定空間');
      cy.get('.nav-link').eq(2).should('have.text', '後台管理');
    });
  });
  it('should navigate to 空間查詢 page', () => {
    cy.get('.main-nav').within(() => {
      cy.get('.nav-link').eq(0).click();
    });
    cy.url().should('include', '/#/'); // 確認URL包含 '/'
    cy.contains('空間查詢'); // 確認頁面上有 '空間查詢' 字樣
  });
  it('should navigate to 已預定空間 page', () => {
    cy.get('.main-nav').within(() => {
      cy.get('.nav-link').eq(1).click();
    });
    cy.url().should('include', '/#/reserved-page'); // 確認URL包含 '/reserved-page'
    cy.contains('已預定空間'); // 確認頁面上有 '已預定空間' 字樣
  });

  });
