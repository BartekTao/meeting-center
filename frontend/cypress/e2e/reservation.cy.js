describe('ReserveBlock Component', () => {
  beforeEach(() => {
    cy.visit('http://localhost:8888'); // 根據你的本地開發伺服器URL
  });
  
  it('should navigate to 後台管理 page', () => {
    cy.get('.main-nav').within(() => {
      cy.get('.nav-link').eq(2).click();
    });
    cy.url().should('include', '/#/manager'); // 確認URL包含 '/後台管理'
    cy.contains('後台管理'); // 確認頁面上有 '後台管理' 字樣
  });
  it('should load the MainBanner component', () => {
    cy.contains('會議預約'); // 確認頁面上有 '會議預約' 字樣
  });
  it('should fill out the search form and submit', () => {
  });
  it('should check the presence and properties of buttons', () => {
    // 檢查預約按鈕
      cy.get('button.main-button').contains('查詢').click();
      cy.wait(2000);
      cy.get('.main-white-button').eq(1).within(() => {
      cy.get('a').should('contain.text', '預約');
      
    });
    cy.get('.main-white-button').eq(1).find('a').contains('預約').click();
      cy.get('input[name="name"]').clear().type('會議標題示例').should('have.value', '會議標題示例');
      cy.get('input[name="email"]').clear().type('Ted Lin').should('have.value', 'Ted Lin');
      cy.get('textarea[name="content"]').clear().type('會議內容示例').should('have.value', '會議內容示例');
      cy.get('button[type="submit"]').click();
      cy.get('.main-nav').within(() => {
        cy.get('.nav-link').eq(1).click();
      });
      cy.contains('會議標題示例').should('exist');
      cy.get('.main-white-button').find('a').contains('刪除').click();

  });
});
